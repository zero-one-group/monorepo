package middleware

import (
	"log/slog"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/lmittmann/tint"
)

func SlogLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start)

			req := c.Request()
			res := c.Response()
            status := res.Status


            args := []any{
                slog.Int("status", status),
                slog.Any("duration", duration),
                slog.String("client_ip", c.RealIP()),
                slog.String("method", req.Method),
                slog.String("path", req.URL.Path),
            }

			switch {
			case status >= 500 || status >= 400:
				slog.Error("Middleware logger", args...)
			default:
				slog.Info("Middleware logger", args...)
			}

            return err
        }
	}
}

func ColorizeLogging(groups []string, a slog.Attr) slog.Attr {
    const LevelTrace = slog.LevelDebug
    if a.Key == slog.LevelKey && len(groups) == 0 {
        level, ok := a.Value.Any().(slog.Level)
        if ok {
            switch level {
            case slog.LevelError:
                return tint.Attr(9, slog.String(a.Key, "ERR"))
            case slog.LevelWarn:
                return tint.Attr(12, slog.String(a.Key, "WRN"))
            case slog.LevelInfo:
                return tint.Attr(10, slog.String(a.Key, "INF"))
            case LevelTrace:
                return tint.Attr(10, slog.String(a.Key, "TRC"))
            }
        }
    }
    return a
}

