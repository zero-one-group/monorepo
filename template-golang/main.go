package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// Response represents the API response structure
type Response struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func main() {
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
		return c.JSON(http.StatusOK, Response{
			Success: true,
			Message: "All is well!",
			Time:    time.Now(),
		})
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Get host from environment variable, default to 127.0.0.1 if not set
	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	// Get port from environment variable, default to {{ port_number }} if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "{{ port_number }}"
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
