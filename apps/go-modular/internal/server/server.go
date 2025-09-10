package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go-modular/internal/adapter"
	"go-modular/internal/config"
	"go-modular/internal/notification"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	appMiddleware "go-modular/internal/middleware"
	modAuth "go-modular/modules/auth"
	modUser "go-modular/modules/user"
	templateFS "go-modular/templates"
)

// HTTPServer is the main HTTP server struct.
// Logger and Tracer are injected for observability.
type HTTPServer struct {
	httpAddr string
	logger   *slog.Logger
}

func NewHTTPServer(httpAddr string) *HTTPServer {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	return &HTTPServer{
		httpAddr: httpAddr,
		logger:   logger,
	}
}

func (s *HTTPServer) Start() error {
	cfg := config.Get()

	// Initialize Postgres database connection with retry mechanism
	pg, err := s.initializeDatabase(cfg)
	if err != nil {
		s.logger.Error("Failed to connect to Postgres database", "err", err)
		os.Exit(1)
	}
	defer pg.Close()

	var mailer *notification.Mailer
	m, err := notification.NewMailer(notification.MailerOptions{
		SMTPHost:     cfg.Mailer.SMTPHost,
		SMTPPort:     cfg.Mailer.SMTPPort,
		SMTPUsername: cfg.Mailer.SMTPUsername,
		SMTPPassword: cfg.Mailer.SMTPPassword,
		FromName:     cfg.Mailer.SenderName,
		FromAddress:  cfg.Mailer.SenderEmail,
		TemplateFS:   templateFS.TemplateDir,
		Logger:       s.logger,
	})
	if err != nil {
		s.logger.Info("mailer service not configured or failed to initialize, continuing without mailer", "err", err)
	} else {
		mailer = m
		s.logger.Info("mailer service initialized", "host", cfg.Mailer.SMTPHost, "port", cfg.Mailer.SMTPPort)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Register global middlewares
	e.Use(middleware.RequestID())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		// Use the RecoverConfig.LogErrorFunc signature: func(c echo.Context, err error, stack []byte) error
		// Integrate with slog logger and return the error so the centralized HTTPErrorHandler still runs.
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			req := c.Request()
			s.logger.Error("Recovered from panic", "method", req.Method, "path", req.URL.Path, "error", err)
			// return the error to allow downstream error handler to process it
			return err
		},
	}))
	e.Use(appMiddleware.LoggerMiddleware(s.logger))

	// Register primary HTTP server routes
	serverHandler := NewServerHandler(pg.Pool, s.logger)
	serverHandler.RegisterRoutes(e)

	// Create API v1 route group
	apiV1Route := e.Group("/api/v1")

	// Register middlewares for API routes
	apiV1Route.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.App.CORSOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		ExposeHeaders: []string{
			echo.HeaderAccept, echo.HeaderAcceptEncoding, echo.HeaderAuthorization, echo.HeaderCacheControl,
			echo.HeaderConnection, echo.HeaderContentType, echo.HeaderContentLength, echo.HeaderOrigin,
			echo.HeaderXCSRFToken, echo.HeaderXRequestID, "Pragma", "User-Agent", "X-App-Audience",
		},
		AllowCredentials: cfg.App.CORSCredentials, // The request can include user credentials like cookies
		MaxAge:           cfg.App.CORSMaxAge,      // Maximum value not ignored by any of major browsers
	}))

	// Load user module (no auth middleware yet)
	userModule := modUser.NewModule(&modUser.Options{PgPool: pg.Pool, Logger: s.logger})

	// Load auth module (requires user service)
	authModule := modAuth.NewModule(&modAuth.Options{
		PgPool:       pg.Pool,
		Logger:       s.logger,
		UserService:  userModule.GetUserService(),
		JWTSecretKey: []byte(cfg.App.JWTSecretKey),
		BaseURL:      cfg.GetAppBaseURL(),
		Mailer:       mailer,
	})

	// Inject auth middleware into user module so protected user routes use same JWT config
	userModule.Use(authModule.JWTMiddleware())

	// Register the module routes after injecting middleware
	userModule.RegisterRoutes(apiV1Route)
	authModule.RegisterRoutes(apiV1Route)

	s.logger.Info("Starting HTTP server", "addr", s.httpAddr)

	// Start the server
	return e.Start(s.httpAddr)
}

// Initialize PostgreSQL database connection with retry mechanism
func (s *HTTPServer) initializeDatabase(cfg *config.Config) (*adapter.PostgresDB, error) {
	const baseDelay = 2 * time.Second
	const maxDelay = 30 * time.Second
	const defaultMaxRetries = 5

	maxRetries := defaultMaxRetries
	if cfg.Database.PgMaxRetries == -1 {
		maxRetries = -1
	} else if cfg.Database.PgMaxRetries > 0 {
		maxRetries = cfg.Database.PgMaxRetries
	}

	s.logger.Info("Initializing database connection", "max_retries", maxRetries)

	var lastErr error
	attempt := 1

	for {
		pgCfg := adapter.PostgresConfig{URL: cfg.GetDatabaseURL()}
		pg, err := adapter.NewPostgres(pgCfg)
		if err == nil {
			// verify connection with Ping and timeout
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			pingErr := pg.Ping(ctx)
			cancel()

			if pingErr == nil {
				s.logger.Info("Database connection established", "attempt", attempt)
				return pg, nil
			}

			// ping failed, close pool and treat as error
			pg.Close()
			err = fmt.Errorf("ping failed: %w", pingErr)
		}

		lastErr = err
		s.logger.Warn("Database connection attempt failed", "attempt", attempt, "err", lastErr)

		// If not infinite and we've reached max, stop retrying
		if maxRetries != -1 && attempt >= maxRetries {
			break
		}

		// Exponential-ish backoff bounded by maxDelay
		delay := min(baseDelay*time.Duration(attempt), maxDelay)

		s.logger.Info("Retrying database connection", "next_try_in", delay, "attempt", attempt+1)
		time.Sleep(delay)
		attempt++
	}

	return nil, fmt.Errorf("failed to establish database connection after %d attempts: %w", attempt, lastErr)
}
