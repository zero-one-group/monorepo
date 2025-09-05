package app

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zero-one-group/go-modulith/internal/auth"
	"github.com/zero-one-group/go-modulith/internal/config"
	"github.com/zero-one-group/go-modulith/internal/product"
	"github.com/zero-one-group/go-modulith/internal/database"
	"github.com/zero-one-group/go-modulith/internal/errors"
	custommiddleware "github.com/zero-one-group/go-modulith/internal/middleware"
	"github.com/zero-one-group/go-modulith/internal/validator"
	"github.com/zero-one-group/go-modulith/internal/user"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type Server struct {
	echo   *echo.Echo
	config *config.Config
}

func NewServer(
	cfg *config.Config,
	db *database.DB,
	authHandler *auth.Handler,
	userHandler *user.Handler,
	productHandler *product.Handler,
	validator *validator.CustomValidator,
) *Server {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.HTTPErrorHandler = errors.ErrorHandler
	e.Validator = validator

	rateLimiter := custommiddleware.NewRateLimiter(cfg)

	e.Use(middleware.RequestID())
	e.Use(custommiddleware.Recovery())
	e.Use(custommiddleware.RequestLogging())
	e.Use(custommiddleware.CORS(cfg))
	e.Use(rateLimiter.Middleware())
	e.Use(otelecho.Middleware("go-modulith"))

	registerHealthRoutes(e, db)

	api := e.Group("/api/v1")
	authMiddleware := custommiddleware.JWTAuth(cfg)

	authHandler.RegisterRoutes(api, authMiddleware)
	userHandler.RegisterRoutes(api, authMiddleware)
	productHandler.RegisterRoutes(api, authMiddleware)

	return &Server{
		echo:   e,
		config: cfg,
	}
}

func (s *Server) Start() error {
	server := &http.Server{
		Addr:         s.config.Server.Address(),
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
	}

	s.echo.Logger.Infof("Starting server on %s", s.config.Server.Address())

	go func() {
		if err := s.echo.StartServer(server); err != nil && err != http.ErrServerClosed {
			s.echo.Logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.config.Server.ShutdownTimeout)
	defer cancel()

	return s.echo.Shutdown(shutdownCtx)
}

func registerHealthRoutes(e *echo.Echo, db *database.DB) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	e.GET("/ready", func(c echo.Context) error {
		if err := db.Health(c.Request().Context()); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{
				"status": "not ready",
				"error":  err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "ready",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	e.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"version": "1.0.0",
			"commit":  "unknown",
			"build":   time.Now().Format(time.RFC3339),
		})
	})
}