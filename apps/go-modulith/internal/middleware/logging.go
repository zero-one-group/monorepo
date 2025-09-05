package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zero-one-group/go-modulith/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func RequestLogging() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			ctx, span := otel.Tracer("middleware").Start(c.Request().Context(), "request_logging")
			defer span.End()

			req := c.Request()
			log := logger.FromContext(ctx)

			span.SetAttributes(
				attribute.String("http.method", req.Method),
				attribute.String("http.url", req.URL.String()),
				attribute.String("http.remote_addr", c.RealIP()),
				attribute.String("http.user_agent", req.UserAgent()),
			)

			log.Info("Request started",
				"method", req.Method,
				"path", req.URL.Path,
				"query", req.URL.RawQuery,
				"remote_addr", c.RealIP(),
				"user_agent", req.UserAgent(),
				"request_id", c.Response().Header().Get(echo.HeaderXRequestID),
			)

			err := next(c)

			duration := time.Since(start)
			status := c.Response().Status

			span.SetAttributes(
				attribute.Int("http.status_code", status),
				attribute.String("http.duration", duration.String()),
				attribute.Int64("http.response_size", c.Response().Size),
			)

			logLevel := slog.LevelInfo
			if status >= 400 {
				logLevel = slog.LevelError
				if err != nil {
					span.RecordError(err)
				}
			}

			log.Log(ctx, logLevel, "Request completed",
				"status", status,
				"duration", duration.String(),
				"response_size", c.Response().Size,
				"error", err,
			)

			return err
		}
	}
}