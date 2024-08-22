package routes

import (
	"github.com/labstack/echo/v4"
	"go-app/internal/handler"
)

// SetupRoutes mengonfigurasi semua rute untuk aplikasi Echo
func SetupRoutes(e *echo.Echo) {
	// Root and health check routes
	e.GET("/", handler.RootHandler)
	e.GET("/health", handler.HealthCheckHandler)

	// Routes group for authenticated routes
	auth := e.Group("/auth")
	{
		auth.POST("/login", handler.LoginHandler)
		auth.POST("/forgot-password", handler.ForgotPasswordHandler)
		auth.PATCH("/reset-password", handler.ResetPasswordHandler)
	}
}
