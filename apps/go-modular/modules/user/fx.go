package user

import (
	"log/slog"

	appMiddleware "go-modular/internal/middleware"
	"go-modular/internal/server"
	"go-modular/modules/user/handler"
	"go-modular/modules/user/repository"
	"go-modular/modules/user/services"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// Module wires the User feature.
var Module = fx.Module("user",
	fx.Provide(
		repository.NewUserRepository,
		fx.Annotate(
			newUserService,
			fx.As(new(services.UserServiceInterface)),
		),
		newUserHandler,
	),
	fx.Invoke(registerRoutes),
)

func newUserService(repo *repository.UserRepository) *services.UserService {
	return services.NewUserService(services.UserServiceOpts{UserRepo: repo})
}

func newUserHandler(log *slog.Logger, svc services.UserServiceInterface) *handler.Handler {
	return handler.NewHandler(&handler.HandlerOpts{
		Logger:      log,
		UserService: svc,
	})
}

// All /users routes require JWT (admin-managed; swagger BearerAuth on CreateUser too).
func registerRoutes(api server.APIV1, h *handler.Handler, jwt appMiddleware.JWTAccess) {
	g := api.Root.Group("/users", echo.MiddlewareFunc(jwt))
	g.POST("", h.CreateUser)
	g.GET("", h.ListUsers)
	g.GET("/:userId", h.GetUser)
	g.PUT("/:userId", h.UpdateUser)
	g.DELETE("/:userId", h.DeleteUser)
}
