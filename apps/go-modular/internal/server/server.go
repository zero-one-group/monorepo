package server

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-modular/internal/adapter"
	"go-modular/internal/middleware/logger"
	"go-modular/internal/notification"

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
	databaseURL := os.Getenv("DATABASE_URL")

	// Initialize Postgres pool
	pg, err := adapter.NewPostgres(adapter.PostgresConfig{
		URL:        databaseURL,
		SearchPath: "public",
	})
	if err != nil {
		slog.Error("Failed to connect to Postgres database", "err", err)
		os.Exit(1)
	}
	defer pg.Close()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Register global middlewares
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(logger.LoggerMiddleware(s.logger))

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

	// Load and register user module
	userModule := user_module.NewModule(&user_module.Options{PgPool: pg.Pool, Logger: s.logger})
	userModule.RegisterRoutes(apiV1Route)

	// Read SMTP configuration from environment
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587
	if p := os.Getenv("SMTP_PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			smtpPort = v
		}
	}
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpSenderEmail := strings.Trim(os.Getenv("SMTP_SENDER_EMAIL"), "\"")
	smtpSenderName := strings.Trim(os.Getenv("SMTP_SENDER_NAME"), "\"")
	appBaseURL := os.Getenv("APP_BASE_URL")

	var mailer *notification.Mailer
	m, err := notification.NewMailer(notification.MailerOptions{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUser,
		SMTPPassword: smtpPass,
		FromName:     smtpSenderName,
		FromAddress:  smtpSenderEmail,
		TemplateFS:   templateFS.TemplateDir,
		Logger:       s.logger,
	})
	if err != nil {
		s.logger.Info("mailer not configured or failed to initialize, continuing without mailer", "err", err)
	} else {
		mailer = m
		s.logger.Info("mailer initialized", "host", smtpHost, "port", smtpPort, "from", smtpSenderEmail)
	}

	// Load and register auth module (pass mailer if available)
	authModule := auth_module.NewModule(&auth_module.Options{
		PgPool:       pg.Pool,
		Logger:       s.logger,
		UserService:  userModule.GetUserService(),
		JWTSecretKey: []byte(os.Getenv("APP_JWT_SECRET_KEY")),
		BaseURL:      appBaseURL,
		Mailer:       mailer,
	})
	authModule.RegisterRoutes(apiV1Route)

	s.logger.Info("Starting HTTP server", "addr", s.httpAddr)

	// Start the server
	return e.Start(s.httpAddr)
}
