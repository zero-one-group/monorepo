package rest

import (
	"context"
	"log/slog"
	"net/http"
	"{{package_name}}/domain"
	apperrors "{{package_name}}/internal/errors"
	"{{package_name}}/internal/rest/middleware"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type UserService interface {
	CreateUser(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, error)
	GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserHandler struct {
	Service UserService
}

func NewUserHandler(e *echo.Group, svc UserService) {
	handler := &UserHandler{
		Service: svc,
	}
	userGroup := e.Group("/users") // users group

	userGroup.GET("", handler.GetUserList)
	userGroup.GET("/:id", handler.GetUser)
	userGroup.POST("", handler.CreateUser)
	userGroup.PUT("/:id", handler.UpdateUser)
	userGroup.DELETE("/:id", handler.DeleteUser)
}

func (h *UserHandler) GetUserList(c echo.Context) error {
	tracer := otel.Tracer("http.handler.user")
	ctx, span := tracer.Start(c.Request().Context(), "GetUserListHandler")
	defer span.End()

	filter := new(domain.UserFilter)
	if err := c.Bind(filter); err != nil {
		span.RecordError(err)
		
		slog.WarnContext(ctx, "Invalid filter parameters in user list request",
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.RealIP()),
		)
		
		// Return validation error using our error system
		validationErr := apperrors.NewValidationError("Invalid filter parameters", err).WithTrace(span)
		return validationErr
	}

	span.SetAttributes(
		attribute.Bool("filtered", filter.Search != ""),
		attribute.String("search_term", filter.Search),
	)

	users, err := h.Service.GetUserList(ctx, filter)
	if err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to get user list",
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.RealIP()),
		)
		
		// Let the error middleware handle the response
		return err
	}

	// Ensure we return an empty array instead of null
	if users == nil {
		users = []domain.User{}
	}

	slog.InfoContext(ctx, "User list retrieved successfully",
		slog.Int("count", len(users)),
		slog.Bool("filtered", filter.Search != ""),
	)

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.User]{
		Data:    users,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Successfully retrieved user list",
	})
}

func (h *UserHandler) GetUser(c echo.Context) error {
	tracer := otel.Tracer("http.handler.user")
	ctx, span := tracer.Start(c.Request().Context(), "GetUserHandler")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		span.RecordError(err)
		
		slog.WarnContext(ctx, "Invalid user ID format in request",
			slog.String("id_param", idParam),
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.RealIP()),
		)
		
		// Return validation error using our error system
		validationErr := apperrors.NewValidationError("Invalid user ID format", err).WithTrace(span)
		return validationErr
	}

	span.SetAttributes(attribute.String("user.id", id.String()))

	user, err := h.Service.GetUser(ctx, id)
	if err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to get user",
			slog.String("user_id", id.String()),
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.RealIP()),
		)
		
		// Let the error middleware handle the response
		return err
	}

	slog.InfoContext(ctx, "User retrieved successfully",
		slog.String("user_id", user.ID),
	)

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.User]{
		Data:    *user,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Successfully retrieved user",
	})
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	tracer := otel.Tracer("http.handler.user")
	ctx, span := tracer.Start(c.Request().Context(), "CreateUserHandler")
	defer span.End()

	var user domain.CreateUserRequest
	if err := c.Bind(&user); err != nil {
		span.RecordError(err)
		
		slog.WarnContext(ctx, "Invalid request payload for create user",
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.RealIP()),
		)
		
		// Return validation error using our error system
		validationErr := apperrors.NewValidationError("Invalid request payload", err).WithTrace(span)
		return validationErr
	}

	span.SetAttributes(
		attribute.String("user.email", user.Email),
		attribute.String("user.name", user.Name),
	)

	createdUser, err := h.Service.CreateUser(ctx, &user)
	if err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to create user",
			slog.String("user_email", user.Email),
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.RealIP()),
		)
		
		// Let the error middleware handle the response
		return err
	}

	slog.InfoContext(ctx, "User created successfully",
		slog.String("user_id", createdUser.ID),
		slog.String("user_email", createdUser.Email),
	)

	return c.JSON(http.StatusCreated, domain.ResponseSingleData[domain.User]{
		Data:    *createdUser,
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "User successfully created",
	})
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	tracer := otel.Tracer("http.handler.user")
	ctx, span := tracer.Start(c.Request().Context(), "UpdateUserHandler")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		span.RecordError(err)
		
		slog.WarnContext(ctx, "Invalid user ID format in update request",
			slog.String("id_param", idParam),
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.RealIP()),
		)
		
		// Return validation error using our error system
		validationErr := apperrors.NewValidationError("Invalid user ID format", err).WithTrace(span)
		return validationErr
	}

	var user domain.User
	if err := c.Bind(&user); err != nil {
		span.RecordError(err)
		
		slog.WarnContext(ctx, "Invalid request payload for update user",
			slog.String("user_id", id.String()),
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.RealIP()),
		)
		
		// Return validation error using our error system
		validationErr := apperrors.NewValidationError("Invalid request payload", err).WithTrace(span)
		return validationErr
	}

	span.SetAttributes(
		attribute.String("user.id", id.String()),
		attribute.String("user.email", user.Email),
		attribute.String("user.name", user.Name),
	)

	updatedUser, err := h.Service.UpdateUser(ctx, id, &user)
	if err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to update user",
			slog.String("user_id", id.String()),
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.RealIP()),
		)
		
		// Let the error middleware handle the response
		return err
	}

	slog.InfoContext(ctx, "User updated successfully",
		slog.String("user_id", updatedUser.ID),
		slog.String("user_email", updatedUser.Email),
	)

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.User]{
		Data:    *updatedUser,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "User successfully updated",
	})
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	tracer := otel.Tracer("http.handler.user")
	ctx, span := tracer.Start(c.Request().Context(), "DeleteUserHandler")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		span.RecordError(err)
		
		slog.WarnContext(ctx, "Invalid user ID format in delete request",
			slog.String("id_param", idParam),
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.RealIP()),
		)
		
		// Return validation error using our error system
		validationErr := apperrors.NewValidationError("Invalid user ID format", err).WithTrace(span)
		return validationErr
	}

	span.SetAttributes(attribute.String("user.id", id.String()))

	if err := h.Service.DeleteUser(ctx, id); err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to delete user",
			slog.String("user_id", id.String()),
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.RealIP()),
		)
		
		// Let the error middleware handle the response
		return err
	}

	slog.InfoContext(ctx, "User deleted successfully",
		slog.String("user_id", id.String()),
	)

	return c.JSON(http.StatusNoContent, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusNoContent,
		Status:  "success",
		Message: "User successfully deleted",
	})
}
