package logger

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"go-modular/pkg/apputils"
)

// LoggerMiddleware returns an Echo middleware that logs HTTP requests using slog.
func LoggerMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			err := next(c)
			stop := time.Now()

			latency := stop.Sub(start)
			status := res.Status
			method := req.Method
			uri := req.RequestURI
			ip := c.RealIP()
			realUserAgent := req.UserAgent()
			userAgent := apputils.SummarizeUserAgent(realUserAgent)
			requestID := req.Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = res.Header().Get(echo.HeaderXRequestID)
			}

			logAttrs := []slog.Attr{
				slog.String("request_id", requestID),
				slog.String("method", method),
				slog.String("uri", uri),
				slog.Int("status", status),
				slog.String("ip", ip),
				slog.String("user_agent", userAgent),
				slog.String("latency", formatLatency(latency)),
			}
			if err != nil {
				logAttrs = append(logAttrs, slog.String("error", err.Error()))
				logger.Error("HTTP request", attrsToArgs(logAttrs)...)
			} else {
				logger.Info("HTTP request", attrsToArgs(logAttrs)...)
			}
			return err
		}
	}
}

func formatLatency(d time.Duration) string {
	ns := d.Nanoseconds()
	switch {
	case ns < 1_000:
		return fmt.Sprintf("%dns", ns)
	case ns < 1_000_000:
		return fmt.Sprintf("%.2fµs", float64(ns)/1_000)
	case ns < 1_000_000_000:
		return fmt.Sprintf("%.2fms", float64(ns)/1_000_000)
	default:
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

func attrsToArgs(attrs []slog.Attr) []any {
	args := make([]any, 0, len(attrs)*2)
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value.Any())
	}
	return args
}
