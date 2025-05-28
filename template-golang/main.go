package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"{{ package_name }}/config"
	"{{ package_name }}/database"
	"{{ package_name }}/domain"
	router "{{ package_name }}/route"
	"{{ package_name }}/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func init() {
	config.LoadEnv()
}

func main() {

	dbPool, err := database.SetupPgxPool()
	if err != nil {
		log.Fatal("Failed to set up database: " + err.Error())
	}
	defer dbPool.Close()

	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.INFO)

	// Register logger middleware
	// @see: https://echo.labstack.com/docs/middleware/logger
	e.Use(middleware.Logger())

	// Register CORS middleware
	// @see: https://echo.labstack.com/docs/middleware/cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Signature"},
	}))

	// Register the routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, domain.Response{
			Success: true,
			Message: "All is well!",
			Time:    time.Now(),
		})
	})
	apiV1 := e.Group("/api/v1")
	svc := service.NewUserService()
	router.RegisterUserRoutes(apiV1, svc)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Get host from environment variable, default to 127.0.0.1 if not set
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	// Get port from environment variable, default to 8000 if not set
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	// Server address and port to listen on
	serverAddr := fmt.Sprintf("%s:%s", host, port)

	go func() {
		e.Logger.Infof("Server starting on %s", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
