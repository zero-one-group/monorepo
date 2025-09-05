package user

import (
	"net/http"

	"github.com/google/uuid"
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
	users := e.Group("/users", authMiddleware)
	{
		users.GET("", h.GetUsers)
		users.GET("/profile", h.GetProfile)
		users.GET("/:id", h.GetUserByID)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}

func (h *Handler) GetUsers(c echo.Context) error {
	ctx, span := otel.Tracer("user").Start(c.Request().Context(), "handler.get_users")
	defer span.End()

	var filters UserFilters
	if err := c.Bind(&filters); err != nil {
		span.RecordError(err)
		return errors.ErrBadRequest
	}

	if err := c.Validate(&filters); err != nil {
		return err
	}

	response, err := h.service.GetUsers(ctx, filters)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetUserByID(c echo.Context) error {
	ctx, span := otel.Tracer("user").Start(c.Request().Context(), "handler.get_user_by_id")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return errors.ErrBadRequest.WithDetails(map[string]string{
			"id": "Invalid UUID format",
		})
	}

	response, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetProfile(c echo.Context) error {
	ctx, span := otel.Tracer("user").Start(c.Request().Context(), "handler.get_profile")
	defer span.End()

	userIDStr := middleware.GetUserID(c)
	if userIDStr == "" {
		return errors.ErrUnauthorized
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.ErrUnauthorized
	}

	response, err := h.service.GetProfile(ctx, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	ctx, span := otel.Tracer("user").Start(c.Request().Context(), "handler.update_user")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return errors.ErrBadRequest.WithDetails(map[string]string{
			"id": "Invalid UUID format",
		})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		span.RecordError(err)
		return errors.ErrBadRequest
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	response, err := h.service.UpdateUser(ctx, id, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	ctx, span := otel.Tracer("user").Start(c.Request().Context(), "handler.delete_user")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return errors.ErrBadRequest.WithDetails(map[string]string{
			"id": "Invalid UUID format",
		})
	}

	if err := h.service.DeleteUser(ctx, id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}