package rest

import (
	"context"
	"go-app/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, *domain.User, error)
	ValidateToken(token string) (string, error)
}

type AuthHandler struct {
	Service AuthService
}


func NewAuthHandler(e *echo.Group, svc AuthService) {
	handler := &AuthHandler{
		Service: svc,
	}
	authGroup := e.Group("/auth")

	authGroup.POST("/login", handler.Login)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	token, user, err := h.Service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusUnauthorized,
			Status:  "error",
			Message: "Invalid email or password",
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[map[string]any]{
		Data: map[string]any{
			"user":  user,
			"token": token,
		},
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Successfully logged in",
	})
}
