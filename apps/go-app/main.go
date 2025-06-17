package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go-app/config"
	"go-app/database"
	"go-app/domain"
	"go-app/internal/repository/postgres"
	"go-app/internal/rest"
	"go-app/internal/rest/middleware"
	"go-app/service"

	"github.com/labstack/echo/v4"
	"github.com/lmittmann/tint"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func init() {
	config.LoadEnv()
}

func main() {


    env := os.Getenv("APP_ENVIRONMENT")
    var handler slog.Handler

    w := os.Stdout
    if env == "local" {
        handler = tint.NewHandler(w, &tint.Options{
            ReplaceAttr: middleware.ColorizeLogging,
        })
    } else {
        // or continue setup log for another env
        handler = slog.NewTextHandler(w, nil)
    }

	logger := slog.New(handler)
	slog.SetDefault(logger)

	dbPool, err := database.SetupPgxPool()
	if err != nil {
		slog.Error("Failed to set up database", slog.String("error", err.Error()))
        os.Exit(1)
	}
	defer dbPool.Close()

	e := echo.New()
	e.HideBanner = true

	e.Logger.SetOutput(os.Stdout)
	e.Logger.SetLevel(0)

    e.Use(middleware.SlogLoggerMiddleware())
	e.Use(middleware.Cors())

	ctx := context.Background()
	tp, shutdown, err := config.InitTracer(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init tracer: %v\n", err)
		os.Exit(1)
	}
	// make sure we flush any spans on exit
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdown(ctx)
	}()

	// 2) Create Echo and register middleware
	e.Use(
		otelecho.Middleware(
			os.Getenv("SERVICE_NAME"),
			otelecho.WithTracerProvider(tp),
		),
	)

	// Register the routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, domain.Response{
			Code:    200,
			Status:  "Succes",
			Message: "All is well!",
		})
	})

	userRepo := postgres.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepo)

	apiV1 := e.Group("/api/v1")
	usersGroup := apiV1.Group("")

    apiV1 := e.Group("/api/v1")
    usersGroup := apiV1.Group("")
    rest.NewUserHandler(usersGroup, userService)

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
        slog.Info("Server starting", "address", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
            os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

    slog.Info("Shutting down server gracefully...")
	if err := e.Shutdown(ctx); err != nil {
		slog.Error("Shutdown error", "error", err)
	}
}


