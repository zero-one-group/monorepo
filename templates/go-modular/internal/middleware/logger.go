package middleware

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

			status := res.Status
			latency := stop.Sub(start)
			clientIP := c.RealIP()
			realUserAgent := req.UserAgent()
			userAgent := apputils.SummarizeUserAgent(realUserAgent)
			requestID := req.Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = res.Header().Get(echo.HeaderXRequestID)
			}

			logAttrs := []slog.Attr{
				slog.Int("status", status),
				slog.String("method", req.Method),
				slog.String("request_id", requestID),
				slog.String("path", req.RequestURI),
				slog.String("client_ip", clientIP),
				slog.String("user_agent", userAgent),
				slog.String("duration", formatLatency(latency)),
				slog.Int64("bytes_in", req.ContentLength),
				slog.Int64("bytes_out", res.Size),
			}

			// Add query parameters if present
			if req.URL.RawQuery != "" {
				logAttrs = append(logAttrs, slog.String("query", req.URL.RawQuery))
			}

			// Log with appropriate level based on status code
			switch {
			case status >= 500:
				slog.Error("HTTP Request", attrsToArgs(logAttrs)...)
			case status >= 400:
				slog.Warn("HTTP Request", attrsToArgs(logAttrs)...)
			default:
				slog.Info("HTTP Request", attrsToArgs(logAttrs)...)
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
