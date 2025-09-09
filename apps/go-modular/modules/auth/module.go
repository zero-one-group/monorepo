package auth

import (
	"log/slog"
	"os"

	"go-modular/modules/auth/handler"
	"go-modular/modules/auth/repository"
	"go-modular/modules/auth/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Options struct {
	PgPool *pgxpool.Pool // PostgreSQL connection pool (required)
	Logger *slog.Logger  // Slog logger instance (optional)
}

// AuthModule holds dependencies for auth-related handlers.
type AuthModule struct {
	logger      *slog.Logger
	middlewares []echo.MiddlewareFunc
	handler     *handler.Handler
}

// NewModule creates a new AuthModule.
func NewModule(opts *Options) *AuthModule {
	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	authRepo := repository.NewAuthRepository(opts.PgPool, logger)
	authService := services.NewAuthService(services.AuthServiceOpts{
		AuthRepo: authRepo,
	})

	h := handler.NewHandler(&handler.HandlerOpts{
		Logger:      logger,
		AuthService: authService,
	})

	return &AuthModule{logger: logger, handler: h}
}

// Use adds middleware(s) to the AuthModule (grouped).
func (m *AuthModule) Use(mw ...echo.MiddlewareFunc) {
	m.middlewares = append(m.middlewares, mw...)
}

// RegisterRoutes registers auth endpoints to the given Echo group.
func (m *AuthModule) RegisterRoutes(e *echo.Group) {
	g := e.Group("/auth", m.middlewares...)
	g.POST("/password", m.handler.SetUserPassword)
	g.PUT("/password/:userId", m.handler.UpdateUserPassword)
}
