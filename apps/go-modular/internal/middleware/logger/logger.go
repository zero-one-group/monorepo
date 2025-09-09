package logger

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mileusna/useragent"
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
			userAgent := summarizeUserAgent(realUserAgent)
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

// Returns a concise summary like "BrowserName vX.Y on OS X.Y"
func summarizeUserAgent(uaString string) string {
	ua := useragent.Parse(uaString)
	name := ua.Name
	version := ua.Version
	osName := ua.OS
	osVersion := ua.OSVersion

	// Ambil major.minor version browser
	majorMinor := ""
	if version != "" {
		parts := strings.Split(version, ".")
		if len(parts) >= 2 {
			majorMinor = parts[0] + "." + parts[1]
		} else {
			majorMinor = parts[0]
		}
	}

	// Ringkas nama OS populer dan ambil versi utama OS jika ada
	// macOS: "Mac OS X" → "macOS"
	if strings.Contains(osName, "Mac OS X") || strings.Contains(osName, "Intel Mac OS X") {
		osName = "macOS"
	}
	// Windows: "Windows NT 10.0" → "Windows 10"
	if strings.Contains(osName, "Windows") && osVersion == "10.0" {
		osName = "Windows"
		osVersion = "10"
	}
	if strings.Contains(osName, "Windows") && osVersion == "11.0" {
		osName = "Windows"
		osVersion = "11"
	}
	// iOS: "iPhone OS" → "iOS"
	if strings.Contains(osName, "iPhone OS") {
		osName = "iOS"
	}
	// Android: "Android"
	if strings.Contains(osName, "Android") {
		osName = "Android"
	}

	if name == "" && osName == "" {
		return "Unknown"
	}
	result := name
	if majorMinor != "" {
		result += " v" + majorMinor
	}
	if osName != "" {
		result += " on " + osName
		if osVersion != "" && osVersion != "0" {
			// Only append version if not empty or "0"
			result += " " + osVersion
		}
	}
	return strings.TrimSpace(result)
}
