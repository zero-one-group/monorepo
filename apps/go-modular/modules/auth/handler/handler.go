package handler

import (
	"log/slog"

	"go-modular/modules/auth/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// HandlerInterface defines the contract for auth handlers.
type HandlerInterface interface {
	SetUserPassword(c echo.Context) error
	UpdateUserPassword(c echo.Context) error

	// Session handlers
	CreateSession(c echo.Context) error
	UpdateSession(c echo.Context) error
	GetSession(c echo.Context) error
	DeleteSession(c echo.Context) error

	// Refresh token handlers
	CreateRefreshToken(c echo.Context) error
	UpdateRefreshToken(c echo.Context) error
	GetRefreshToken(c echo.Context) error
	DeleteRefreshToken(c echo.Context) error
}

// Ensure Handler implements HandlerInterface
var _ HandlerInterface = (*Handler)(nil)

// Handler holds dependencies for auth handlers.
type Handler struct {
	logger      *slog.Logger
	authService services.AuthServiceInterface
	validator   *validator.Validate
}

type HandlerOpts struct {
	Logger      *slog.Logger
	AuthService services.AuthServiceInterface
}

// NewHandler creates a new Handler instance.
func NewHandler(opts *HandlerOpts) *Handler {
	return &Handler{
		logger:      opts.Logger,
		authService: opts.AuthService,
		validator:   validator.New(),
	}
}

// Helper to convert validator.ValidationErrors to a readable map
func validationErrorsToMap(err error) map[string]string {
	errs := map[string]string{}
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()
			var msg string
			switch tag {
			case "required":
				msg = "This field is required"
			case "uuid":
				msg = "Must be a valid UUID"
			case "min":
				msg = "Minimum length is " + fe.Param()
			case "eqfield":
				msg = "Must match " + fe.Param()
			default:
				msg = "Invalid value"
			}
			errs[field] = msg
		}
	} else {
		errs["error"] = err.Error()
	}
	return errs
}
