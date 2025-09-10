package server

import (
	"go-modular/internal/adapter"
	"go-modular/internal/config"
	"go-modular/internal/middleware"
	"go-modular/internal/notification"

	"github.com/labstack/echo/v4"

	modAuth "go-modular/modules/auth"
	modUser "go-modular/modules/user"
)

// registerModules registers application modules, injects middleware and attaches routes.
// Keeps Start() concise and centralizes module wiring for easier testing/refactor.
func (s *HTTPServer) registerModules(cfg *config.Config, pg *adapter.PostgresDB, mailer *notification.Mailer, e *echo.Echo) error {
	// Register primary HTTP server routes
	serverHandler := NewServerHandler(pg.Pool, s.logger)
	serverHandler.RegisterRoutes(e)

	// Create API v1 route group
	apiV1Route := e.Group("/api/v1")

	// Register middlewares for API routes
	apiV1Route.Use(middleware.CORSMiddleware(cfg))
	apiV1Route.Use(middleware.RateLimitMiddleware(
		cfg.App.RateLimitRequests, cfg.App.RateLimitBurstSize,
	))
	apiV1Route.Use(middleware.CompressionMiddleware())

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

	return nil
}
