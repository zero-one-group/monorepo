package handler

import (
	"log/slog"

	"{{ package_name | kebab_case }}/modules/auth/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// HandlerInterface defines the contract for auth handlers.
type HandlerInterface interface {
	// Password handlers
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

	// Email verification handlers
	InitiateEmailVerification(c echo.Context) error
	ValidateEmailVerification(c echo.Context) error
	ValidateEmailVerificationByLink(c echo.Context) error
	RevokeEmailVerification(c echo.Context) error
	ResendEmailVerification(c echo.Context) error
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
