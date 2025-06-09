package middleware

import (
	"log/slog"
	"time"

	echo "github.com/labstack/echo/v4"
)

func SlogLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Since(start)

			req := c.Request()
			res := c.Response()

			slog.Info("HTTP Request",
				"status", res.Status,
				"duration", stop,
				"client_ip", c.RealIP(),
				"method", req.Method,
				"path", req.URL.Path,
			)

			return err
		}
	}
}
