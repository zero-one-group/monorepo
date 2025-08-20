package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ErrorCode represents different types of application errors
type ErrorCode string

const (
	ErrCodeInternal        ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeConflict        ErrorCode = "CONFLICT"
	ErrCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrCodeDatabase        ErrorCode = "DATABASE_ERROR"
	ErrCodeExternal        ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeRateLimit       ErrorCode = "RATE_LIMIT_EXCEEDED"
)

// AppError represents a structured application error
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	Cause      error                  `json:"-"`
	Context    map[string]interface{} `json:"context,omitempty"`
	Stack      string                 `json:"-"`
	TraceID    string                 `json:"trace_id,omitempty"`
	StatusCode int                    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap implements the unwrapper interface for Go 1.13+ error handling
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithContext adds context information to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithTrace adds tracing information to the error
func (e *AppError) WithTrace(span trace.Span) *AppError {
	if span.SpanContext().IsValid() {
		e.TraceID = span.SpanContext().TraceID().String()
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Message)
	}
	return e
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, cause error) *AppError {
	// Capture stack trace
	_, file, line, _ := runtime.Caller(1)
	
	err := &AppError{
		Code:       code,
		Message:    message,
		Cause:      cause,
		Stack:      fmt.Sprintf("%s:%d", file, line),
		StatusCode: getHTTPStatusCode(code),
	}
	
	return err
}

// Wrap creates a new AppError by wrapping an existing error
func Wrap(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}
	
	// If it's already an AppError, preserve the original cause
	if appErr, ok := err.(*AppError); ok {
		return NewAppError(code, message, appErr.Cause)
	}
	
	return NewAppError(code, message, err)
}

// WrapWithContext wraps an error with additional context
func WrapWithContext(err error, code ErrorCode, message string, ctx map[string]interface{}) *AppError {
	appErr := Wrap(err, code, message)
	if appErr != nil {
		appErr.Context = ctx
	}
	return appErr
}

// getHTTPStatusCode maps error codes to HTTP status codes
func getHTTPStatusCode(code ErrorCode) int {
	switch code {
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeBadRequest, ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrCodeInternal, ErrCodeDatabase, ErrCodeExternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Common error constructors
func NewNotFoundError(message string, cause error) *AppError {
	return NewAppError(ErrCodeNotFound, message, cause)
}

func NewValidationError(message string, cause error) *AppError {
	return NewAppError(ErrCodeValidation, message, cause)
}

func NewInternalError(message string, cause error) *AppError {
	return NewAppError(ErrCodeInternal, message, cause)
}

func NewDatabaseError(message string, cause error) *AppError {
	return NewAppError(ErrCodeDatabase, message, cause)
}

func NewBadRequestError(message string, cause error) *AppError {
	return NewAppError(ErrCodeBadRequest, message, cause)
}

func NewConflictError(message string, cause error) *AppError {
	return NewAppError(ErrCodeConflict, message, cause)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extracts AppError from error chain
func GetAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
