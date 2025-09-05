package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zero-one-group/go-modulith/internal/errors"
	"github.com/zero-one-group/go-modulith/internal/middleware"
	"go.opentelemetry.io/otel"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(e *echo.Group, authMiddleware echo.MiddlewareFunc) {
	auth := e.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout, authMiddleware)
	}
}

func (h *Handler) Register(c echo.Context) error {
	ctx, span := otel.Tracer("auth").Start(c.Request().Context(), "handler.register")
	defer span.End()

	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		span.RecordError(err)
		return errors.ErrBadRequest
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	response, err := h.service.Register(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *Handler) Login(c echo.Context) error {
	ctx, span := otel.Tracer("auth").Start(c.Request().Context(), "handler.login")
	defer span.End()

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		span.RecordError(err)
		return errors.ErrBadRequest
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	response, err := h.service.Login(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) RefreshToken(c echo.Context) error {
	ctx, span := otel.Tracer("auth").Start(c.Request().Context(), "handler.refresh_token")
	defer span.End()

	var req RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		span.RecordError(err)
		return errors.ErrBadRequest
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	response, err := h.service.RefreshToken(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) Logout(c echo.Context) error {
	ctx, span := otel.Tracer("auth").Start(c.Request().Context(), "handler.logout")
	defer span.End()

	userID := middleware.GetUserID(c)
	if userID == "" {
		return errors.ErrUnauthorized
	}

	if err := h.service.Logout(ctx, userID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}