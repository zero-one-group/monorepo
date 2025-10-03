package rest

import (
	"context"
	"{{ package_name | kebab_case }}/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*domain.LoginResponse, error)
}

type AuthHandler struct {
	Service AuthService
}

func NewAuthHandler(e *echo.Group, svc AuthService) {
	handler := &AuthHandler{
		Service: svc,
	}

	e.POST("/login", handler.Login)
}

// @Summary        Login In
// @Description    Authenticate user
// @Tags           Users Authentication
// @Accept         json
// @Produce        json
// @Param          json    body        domain.LoginRequest                     true         "User signin credentials"
// @Success        200     {object}    domain.ResponseSingleData[domain.LoginResponse]      "Successfully logged in"
// @Failure        400     {object}    domain.ResponseSingleData[domain.Empty]              "Bad request"
// @Failure        401     {object}    domain.ResponseSingleData[domain.Empty]              "Unauthorized"
// @Failure        500     {object}    domain.ResponseSingleData[domain.Empty]              "Internal server error"
// @Router         /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		})
	}

	ctx := c.Request().Context()
	result, err := h.Service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusUnauthorized,
			Message: "Invalid email or password",
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.LoginResponse]{
		Data:    *result,
		Code:    http.StatusOK,
		Message: "Successfully logged in",
	})
}
