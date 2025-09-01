package middleware

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "request_id"
)

// RequestIDMiddleware adds a unique request ID to each request for log correlation
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if request ID is already provided in headers
			requestID := c.Request().Header.Get(RequestIDHeader)
			if requestID == "" {
				// Generate new request ID if not provided
				requestID = uuid.New().String()
			}

			// Set request ID in response headers
			c.Response().Header().Set(RequestIDHeader, requestID)

			// Store request ID in context for use in handlers
			ctx := context.WithValue(c.Request().Context(), RequestIDKey, requestID)
			c.SetRequest(c.Request().WithContext(ctx))

			// Add request ID to the Echo context for easy access
			c.Set(RequestIDKey, requestID)

			return next(c)
		}
	}
}

// GetRequestID extracts the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
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
