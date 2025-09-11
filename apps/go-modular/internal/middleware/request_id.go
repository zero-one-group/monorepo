package middleware

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
	"go.jetify.com/typeid"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

// use an unexported type for context keys to avoid collisions (fixes SA1029)
type ctxKey string

var requestIDCtxKey = ctxKey(RequestIDKey)

// RequestIDMiddleware adds a unique request ID to each request for log correlation
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if request ID is already provided in headers
			requestID := c.Request().Header.Get(RequestIDHeader)
			if _, err := typeid.FromString(requestID); err != nil {
				requestID = generateTypeID()
			}

			// Set request ID in response headers
			c.Response().Header().Set(RequestIDHeader, requestID)

			// Store request ID in context for use in handlers (use typed key)
			ctx := context.WithValue(c.Request().Context(), requestIDCtxKey, requestID)
			c.SetRequest(c.Request().WithContext(ctx))

			// Add request ID to the Echo context for easy access (echo uses string keys)
			c.Set(RequestIDKey, requestID)

			return next(c)
		}
	}
}

// GetRequestID extracts the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDCtxKey).(string); ok {
		return requestID
	}
	return ""
}

// LogWithRequestID creates a logger with request ID context
func LogWithRequestID(ctx context.Context) *slog.Logger {
	requestID := GetRequestID(ctx)
	if requestID != "" {
		return slog.With(slog.String(RequestIDKey, requestID))
	}
	return slog.Default()
}

// GetRequestIDFromEcho extracts request ID from Echo context
func GetRequestIDFromEcho(c echo.Context) string {
	if requestID, ok := c.Get(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

func generateTypeID() string {
	// Generate TypeID with "req" prefix for request tracing
	tid, err := typeid.WithPrefix("req")
	if err != nil {
		// This should never happen with a valid prefix, but fallback to no prefix
		tid, _ = typeid.WithPrefix("")
	}
	return tid.String()
}
