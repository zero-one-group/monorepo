package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-app/internal/routes"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	// Register logger middleware
	// @see: https://echo.labstack.com/docs/middleware/logger
	e.Use(middleware.Logger())

	// Register CORS middleware
	// @see: https://echo.labstack.com/docs/middleware/cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8000", "https://example.com"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Client-Info"},
	}))

	// Register the routes
	routes.SetupRoutes(e)

	// TODO: read the listen address from the environment variable
	listenAddr := "127.0.0.1:8080"
	e.Logger.Fatal(e.Start(listenAddr))
}
