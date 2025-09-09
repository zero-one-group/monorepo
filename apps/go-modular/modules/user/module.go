package user

import (
	"log/slog"
	"os"

	"go-modular/modules/user/handler"
	"go-modular/modules/user/repository"
	"go-modular/modules/user/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Options struct {
	PgPool *pgxpool.Pool // PostgreSQL connection pool (required)
	Logger *slog.Logger  // Slog logger instance (optional)
}

// UserModule holds dependencies for user-related handlers.
type UserModule struct {
	logger      *slog.Logger
	middlewares []echo.MiddlewareFunc
	handler     *handler.Handler
}

// NewModule creates a new UserModule.
func NewModule(opts *Options) *UserModule {
	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	// Initialize required services
	userService := services.NewUserService(services.UserServiceOpts{
		UserRepo: repository.NewUserRepository(opts.PgPool, logger),
	})

	h := handler.NewHandler(&handler.HandlerOpts{
		Logger:      logger,
		UserService: userService,
	})

	return &UserModule{logger: logger, handler: h}
}

// Use adds middleware(s) to the UserModule (grouped).
func (m *UserModule) Use(mw ...echo.MiddlewareFunc) {
	m.middlewares = append(m.middlewares, mw...)
}

// RegisterRoutes registers user CRUD endpoints to the given Echo group.
func (m *UserModule) RegisterRoutes(e *echo.Group) {
	g := e.Group("/users", m.middlewares...)
	g.POST("", m.handler.CreateUser)
	g.GET("", m.handler.ListUsers)
	g.GET("/:userId", m.handler.GetUser)
	g.PUT("/:userId", m.handler.UpdateUser)
	g.DELETE("/:userId", m.handler.DeleteUser)
}
