package rest

import (
	"context"
	"fmt"
	"net/http"
	"go-app/domain"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error)
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
    // userGroup.GET("/:id", handler.GetUserList) // example for get only one user
}

func (h *UserHandler) GetUserList(c echo.Context) error {
    filter := new(domain.UserFilter)
	if err := c.Bind(filter); err != nil {
        fmt.Println(err)
	}

	ctx := c.Request().Context()
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
        Data:       users,
		Code:       http.StatusOK,
		Status:     "Success",
		Message:    "Successfully retrieve user list",
	})

}
