package server

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-modular/internal/adapter"
	"go-modular/internal/middleware/logger"

	auth_module "go-modular/modules/auth"
	user_module "go-modular/modules/user"
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

	// Load and register auth module
	authModule := auth_module.NewModule(&auth_module.Options{PgPool: pg.Pool, Logger: s.logger})
	authModule.RegisterRoutes(apiV1Route)

	s.logger.Info("Starting HTTP server", "addr", s.httpAddr)

	// Start the server
	return e.Start(s.httpAddr)
}
