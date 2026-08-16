package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"go-modular/internal/adapter"
	"go-modular/internal/config"
	appMiddleware "go-modular/internal/middleware"
	"go-modular/internal/server"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

var httpModule = fx.Module("http",
	fx.Provide(newEcho),
	fx.Provide(newAPIV1),
	fx.Provide(server.NewServerHandler),
	fx.Invoke(registerServerRoutes),
	fx.Invoke(runHTTPServer),
)

func newEcho(cfg *config.Config, log *slog.Logger) *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(cfg.GetEchoLogLevel())
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			req := c.Request()
			log.Error("Recovered from panic", "method", req.Method, "path", req.URL.Path, "error", err)
			return err
		},
	}))
	e.Use(appMiddleware.RequestIDMiddleware())
	e.Use(appMiddleware.SecurityHeadersMiddleware())
	e.Use(appMiddleware.TimeoutMiddleware(30 * time.Second))
	e.Use(appMiddleware.LoggerMiddleware(log))
	e.Use(appMiddleware.CORSMiddleware(cfg))
	e.Use(appMiddleware.RateLimitMiddleware(cfg.App.RateLimitRequests, cfg.App.RateLimitBurstSize))
	e.Use(appMiddleware.CompressionMiddleware())

	return e
}

func newAPIV1(e *echo.Echo) server.APIV1 {
	return server.APIV1{Root: e.Group("/api/v1")}
}

func registerServerRoutes(e *echo.Echo, h *server.ServerHandler) {
	h.RegisterRoutes(e)
}

func runHTTPServer(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
	e *echo.Echo,
	cfg *config.Config,
	log *slog.Logger,
	_ *adapter.PostgresDB, // DB must be ready before listen
) {
	addr := cfg.GetServerAddr()

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("Starting HTTP server", "addr", addr)
			go func() {
				if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
					log.Error("HTTP server error", "err", err)
					if shutErr := shutdowner.Shutdown(); shutErr != nil {
						log.Error("failed to trigger fx shutdown after HTTP error", "err", shutErr)
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down HTTP server")
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			return e.Shutdown(shutdownCtx)
		},
	})
}
