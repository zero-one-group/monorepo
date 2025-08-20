package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RequestIDMiddleware generates and attaches a unique request ID to each request
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if request ID already exists in headers
			reqID := c.Request().Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = c.Request().Header.Get("X-Correlation-ID")
			}
			
			// Generate new request ID if none exists
			if reqID == "" {
				reqID = uuid.New().String()
			}
			
			// Set request ID in response header
			c.Response().Header().Set(echo.HeaderXRequestID, reqID)
			c.Response().Header().Set("X-Correlation-ID", reqID)
			
			// Continue with next middleware/handler
			return next(c)
		}
	}
}
