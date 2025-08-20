package domain

import (
	"errors"

	apperrors "{{ package_name }}/internal/errors"
)

// Legacy errors for backward compatibility
var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal server error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given param is not valid")
	// ErrUserNotFound indicates user was not found
	ErrUserNotFound = errors.New("user not found")
)

// Application-specific error constructors using the new error system

// NewUserNotFoundError creates a user not found error
func NewUserNotFoundError(userID string) *apperrors.AppError {
	return apperrors.NewNotFoundError("User not found", ErrUserNotFound).
		WithContext("user_id", userID)
}

// NewUserValidationError creates a user validation error
func NewUserValidationError(field, message string) *apperrors.AppError {
	return apperrors.NewValidationError("User validation failed", ErrBadParamInput).
		WithContext("field", field).
		WithContext("validation_message", message)
}

// NewUserConflictError creates a user conflict error (e.g., email already exists)
func NewUserConflictError(email string) *apperrors.AppError {
	return apperrors.NewConflictError("User with this email already exists", ErrConflict).
		WithContext("email", email)
}

// NewUserInternalError creates an internal error for user operations
func NewUserInternalError(operation string, cause error) *apperrors.AppError {
	return apperrors.NewInternalError("Internal error during user operation", cause).
		WithContext("operation", operation)
}

// NewUserDatabaseError creates a database error for user operations
func NewUserDatabaseError(operation string, cause error) *apperrors.AppError {
	return apperrors.NewDatabaseError("Database error during user operation", cause).
		WithContext("operation", operation)
}
