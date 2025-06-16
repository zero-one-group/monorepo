package rest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"{{package_name}}/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(
		c.Request().Context(),
		"RouteUser.GetUserList",
	)
	defer span.Finish()

	filter := new(domain.UserFilter)
	if err := c.Bind(filter); err != nil {
		fmt.Println(err)
	}

	span.SetTag("filter", filter)
	users, err := h.Service.GetUserList(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to list users: " + err.Error(),
		})
	}
	if users == nil {
		users = []domain.User{}
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.User]{
		Data:    users,
		Code:    http.StatusOK,
		Status:  "Success",
		Message: "Successfully retrieve user list",
	})
}

func (h *UserHandler) GetUser(c echo.Context) error {
	span, ctx := opentracing.StartSpanFromContext(
		c.Request().Context(),
		"RouteUser.GetUser",
	)
	defer span.Finish()

	rawID := c.Param("id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "invalid user id"},
		)
	}

	span.SetTag("user.id", id.String())
	user, err := h.Service.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: "User not found",
			})
		}

		fmt.Println("GetUser error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to get user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.User]{
		Data:    *user,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Successfully retrieved user",
	})
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	span, ctx := opentracing.StartSpanFromContext(
		c.Request().Context(),
		"RouteUser.CreateUser",
	)
	defer span.Finish()

	var user domain.CreateUserRequest
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	span.SetTag("request", user)
	createdUser, err := h.Service.CreateUser(ctx, &user)
	if err != nil {
		fmt.Println("CreateUser error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to create user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, domain.ResponseSingleData[domain.User]{
		Data:    *createdUser,
		Code:    http.StatusCreated,
		Status:  "success",
		Message: "User successfully created",
	})
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	span, ctx := opentracing.StartSpanFromContext(
		c.Request().Context(),
		"RouteUser.UpdateUser",
	)
	defer span.Finish()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid user ID format",
		})
	}

	var user domain.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	span.SetTag("payload", user)
	updatedUser, err := h.Service.UpdateUser(ctx, id, &user)
	if err != nil {
		fmt.Println("UpdateUser error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to update user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.User]{
		Data:    *updatedUser,
		Code:    http.StatusOK,
		Status:  "success",
		Message: "User successfully updated",
	})
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	span, ctx := opentracing.StartSpanFromContext(
		c.Request().Context(),
		"RouteUser.DeleteUser",
	)
	defer span.Finish()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid user ID format",
		})
	}

	span.SetTag("id", id)

	if err := h.Service.DeleteUser(ctx, id); err != nil {
		fmt.Println("DeleteUser error:", err)
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: "Failed to delete user: " + err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, domain.ResponseSingleData[domain.Empty]{
		Code:    http.StatusNoContent,
		Status:  "success",
		Message: "User successfully deleted",
	})
}
