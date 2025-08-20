package middleware

import (
	"database/sql"
	"log/slog"
	"net/http"

	"{{ package_name }}/domain"
	apperrors "{{ package_name }}/internal/errors"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Code      int                    `json:"code"`
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Details   string                 `json:"details,omitempty"`
	TraceID   string                 `json:"trace_id,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
}

// ErrorHandler returns an Echo HTTP error handler for centralized error handling
func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		// Get trace context
		span := trace.SpanFromContext(c.Request().Context())
		traceID := ""
		if span.SpanContext().IsValid() {
			traceID = span.SpanContext().TraceID().String()
		}

		// Get request ID from context if available
		requestID := getRequestID(c)

		// Handle different error types
		var errorResponse ErrorResponse

		switch {
		case apperrors.IsAppError(err):
			handleAppError(err, &errorResponse, traceID, requestID)
		case isEchoHTTPError(err):
			handleEchoHTTPError(err, &errorResponse, traceID, requestID)
		case isDatabaseError(err):
			handleDatabaseError(err, &errorResponse, traceID, requestID)
		default:
			handleGenericError(err, &errorResponse, traceID, requestID)
		}

		// Log the error with context
		logError(c, err, errorResponse)

		// Send the error response
		if err := c.JSON(errorResponse.Code, errorResponse); err != nil {
			slog.Error("Failed to send error response",
				slog.String("error", err.Error()),
				slog.String("trace_id", traceID),
				slog.String("request_id", requestID),
			)
		}
	}
}

// handleAppError processes application-specific errors
func handleAppError(err error, resp *ErrorResponse, traceID, requestID string) {
	appErr, _ := apperrors.GetAppError(err)
	
	resp.Code = appErr.StatusCode
	resp.Status = "error"
	resp.Message = appErr.Message
	resp.Details = appErr.Details
	resp.TraceID = traceID
	resp.RequestID = requestID
	resp.Context = appErr.Context

	// Add trace information to the error if span is available
	if traceID != "" && appErr.TraceID == "" {
		appErr.TraceID = traceID
	}
}

// handleEchoHTTPError processes Echo framework HTTP errors
func handleEchoHTTPError(err error, resp *ErrorResponse, traceID, requestID string) {
	echoErr := err.(*echo.HTTPError)
	
	resp.Code = echoErr.Code
	resp.Status = "error"
	resp.TraceID = traceID
	resp.RequestID = requestID
	
	if msg, ok := echoErr.Message.(string); ok {
		resp.Message = msg
	} else {
		resp.Message = http.StatusText(echoErr.Code)
	}
}

// handleDatabaseError processes database-related errors
func handleDatabaseError(err error, resp *ErrorResponse, traceID, requestID string) {
	resp.TraceID = traceID
	resp.RequestID = requestID
	resp.Status = "error"
	
	switch err {
	case sql.ErrNoRows:
		resp.Code = http.StatusNotFound
		resp.Message = "Resource not found"
	default:
		resp.Code = http.StatusInternalServerError
		resp.Message = "Database operation failed"
		resp.Details = "An error occurred while processing your request"
	}
}

// handleGenericError processes all other errors
func handleGenericError(err error, resp *ErrorResponse, traceID, requestID string) {
	resp.Code = http.StatusInternalServerError
	resp.Status = "error"
	resp.Message = "Internal server error"
	resp.Details = "An unexpected error occurred while processing your request"
	resp.TraceID = traceID
	resp.RequestID = requestID
}

// logError logs the error with appropriate context
func logError(c echo.Context, err error, resp ErrorResponse) {
	// Determine log level based on error type and status code
	logLevel := slog.LevelError
	if resp.Code < 500 {
		logLevel = slog.LevelWarn
	}

	// Create structured log entry
	logAttrs := []slog.Attr{
		slog.String("method", c.Request().Method),
		slog.String("path", c.Request().URL.Path),
		slog.String("user_agent", c.Request().UserAgent()),
		slog.String("remote_addr", c.RealIP()),
		slog.Int("status_code", resp.Code),
		slog.String("error", err.Error()),
	}

	// Add trace and request IDs if available
	if resp.TraceID != "" {
		logAttrs = append(logAttrs, slog.String("trace_id", resp.TraceID))
	}
	if resp.RequestID != "" {
		logAttrs = append(logAttrs, slog.String("request_id", resp.RequestID))
	}

	// Add error context if it's an AppError
	if appErr, ok := apperrors.GetAppError(err); ok && appErr.Context != nil {
		for key, value := range appErr.Context {
			logAttrs = append(logAttrs, slog.Any("ctx_"+key, value))
		}
	}

	slog.LogAttrs(c.Request().Context(), logLevel, "Request error", logAttrs...)
}

// Helper functions

func isEchoHTTPError(err error) bool {
	_, ok := err.(*echo.HTTPError)
	return ok
}

func isDatabaseError(err error) bool {
	return err == sql.ErrNoRows || err == sql.ErrConnDone || err == sql.ErrTxDone
}

func getRequestID(c echo.Context) string {
	// Try to get request ID from various possible sources
	if reqID := c.Response().Header().Get(echo.HeaderXRequestID); reqID != "" {
		return reqID
	}
	if reqID := c.Request().Header.Get("X-Request-ID"); reqID != "" {
		return reqID
	}
	if reqID := c.Request().Header.Get("X-Correlation-ID"); reqID != "" {
		return reqID
	}
	return ""
}

// CreateErrorResponse creates a standardized error response for domain errors
func CreateErrorResponse(code int, status, message string) domain.ResponseSingleData[domain.Empty] {
	return domain.ResponseSingleData[domain.Empty]{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    domain.Empty{},
	}
}
