package handler

import (
	"github.com/labstack/echo/v4"
	"go-app/internal/model"
	"net/http"
)

func LoginHandler(c echo.Context) error {
	response := model.LoginResponse{
		BaseAPIResponse: model.BaseAPIResponse{Status: http.StatusOK, Success: true},
		Data: model.LoginData{
			AccessToken:  "sample_access_token",
			RefreshToken: "sample_refresh_token",
			Role:         "admin",
			User:         nil, // Can be populated with actual user data if needed

		},
	}
	return c.JSON(http.StatusOK, response)
}

func ForgotPasswordHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Endpoint for forgot password"})
}

func ResetPasswordHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Endpoint for reset password"})
}
