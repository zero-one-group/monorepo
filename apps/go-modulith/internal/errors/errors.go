package errors

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error   string            `json:"error"`
	Code    string            `json:"code"`
	Details map[string]string `json:"details,omitempty"`
}

type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Details    map[string]string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    make(map[string]string),
	}
}

func (e *AppError) WithDetails(details map[string]string) *AppError {
	e.Details = details
	return e
}

var (
	ErrUnauthorized = NewAppError("UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized)
	ErrForbidden    = NewAppError("FORBIDDEN", "Forbidden", http.StatusForbidden)
	ErrNotFound     = NewAppError("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrConflict     = NewAppError("CONFLICT", "Resource conflict", http.StatusConflict)
	ErrBadRequest   = NewAppError("BAD_REQUEST", "Bad request", http.StatusBadRequest)
	ErrInternal     = NewAppError("INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
)

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var appErr *AppError
	var httpErr *echo.HTTPError
	var validationErr validator.ValidationErrors

	switch {
	case As(err, &appErr):
		c.JSON(appErr.StatusCode, ErrorResponse{
			Error:   appErr.Message,
			Code:    appErr.Code,
			Details: appErr.Details,
		})
	case As(err, &httpErr):
		code := httpErr.Code
		message := fmt.Sprintf("%v", httpErr.Message)
		
		c.JSON(code, ErrorResponse{
			Error: message,
			Code:  http.StatusText(code),
		})
	case As(err, &validationErr):
		details := make(map[string]string)
		for _, fieldErr := range validationErr {
			details[fieldErr.Field()] = getValidationErrorMessage(fieldErr)
		}
		
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Code:    "VALIDATION_ERROR",
			Details: details,
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}
}

func As(err error, target interface{}) bool {
	switch t := target.(type) {
	case **AppError:
		if appErr, ok := err.(*AppError); ok {
			*t = appErr
			return true
		}
	case **echo.HTTPError:
		if httpErr, ok := err.(*echo.HTTPError); ok {
			*t = httpErr
			return true
		}
	case *validator.ValidationErrors:
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			*t = validationErr
			return true
		}
	}
	return false
}

func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", fe.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Value must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Value must be less than or equal to %s", fe.Param())
	default:
		return "Invalid value"
	}
}