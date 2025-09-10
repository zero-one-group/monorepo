package server

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-modular/internal/adapter"
	"go-modular/internal/config"
	"go-modular/internal/notification"

	appMiddleware "go-modular/internal/middleware"
	auth_module "go-modular/modules/auth"
	user_module "go-modular/modules/user"
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

	// Initialize Postgres pool
	pg, err := adapter.NewPostgres(adapter.PostgresConfig{
		URL:        cfg.GetDatabaseURL(),
		SearchPath: "public",
	})
	if err != nil {
		slog.Error("Failed to connect to Postgres database", "err", err)
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
	e.Use(middleware.Recover())
	e.Use(appMiddleware.LoggerMiddleware(s.logger))

	// Register primary HTTP server routes
	serverHandler := NewServerHandler(pg.Pool, s.logger)
	serverHandler.RegisterRoutes(e)

	// Create API v1 route group
	apiV1Route := e.Group("/api/v1")

	// Register middlewares for API routes
	apiV1Route.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Load user module (no auth middleware yet)
	userModule := user_module.NewModule(&user_module.Options{PgPool: pg.Pool, Logger: s.logger})

	// Load auth module (requires user service)
	authModule := auth_module.NewModule(&auth_module.Options{
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
