// Package models contains HTTP request/response schema definitions for the auth module.
// This file defines struct types used for HTTP payloads, validation, and OpenAPI documentation.

package models

type SetPasswordRequest struct {
	UserID               string `json:"user_id" validate:"required,uuid"`
	Password             string `json:"password" validate:"required,min=8" example:"secure.password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password" example:"secret.password"`
}

type UpdatePasswordRequest struct {
	CurrentPassword      string `json:"current_password" validate:"required,min=8" example:"current.password"`
	NewPassword          string `json:"new_password" validate:"required,min=8" example:"secure.password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=NewPassword" example:"secret.password"`
}
