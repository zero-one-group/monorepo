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

type CreateSessionRequest struct {
	UserID            string  `json:"user_id" validate:"required,uuid"`
	TokenHash         string  `json:"token_hash" validate:"required"`
	UserAgent         *string `json:"user_agent,omitempty"`
	DeviceName        *string `json:"device_name,omitempty"`
	DeviceFingerprint *string `json:"device_fingerprint,omitempty"`
	IPAddress         *string `json:"ip_address,omitempty" example:"192.168.1.1"`
	ExpiresAt         string  `json:"expires_at" validate:"required,datetime" example:"2025-12-31T23:59:59Z"`
}

type UpdateSessionRequest struct {
	SessionID         string  `json:"session_id" validate:"required,uuid"`
	UserAgent         *string `json:"user_agent,omitempty"`
	DeviceName        *string `json:"device_name,omitempty"`
	DeviceFingerprint *string `json:"device_fingerprint,omitempty"`
	IPAddress         *string `json:"ip_address,omitempty" example:"192.168.1.1"`
	RefreshedAt       *string `json:"refreshed_at,omitempty" example:"2025-12-31T23:59:59Z"`
	RevokedAt         *string `json:"revoked_at,omitempty" example:"2025-12-31T23:59:59Z"`
	RevokedBy         *string `json:"revoked_by,omitempty" validate:"omitempty,uuid"`
}

type CreateRefreshTokenRequest struct {
	UserID    string  `json:"user_id" validate:"required,uuid"`
	SessionID *string `json:"session_id,omitempty" validate:"omitempty,uuid"`
	TokenHash string  `json:"token_hash" validate:"required"`
	IPAddress *string `json:"ip_address,omitempty" example:"192.168.1.1"`
	UserAgent *string `json:"user_agent,omitempty"`
	ExpiresAt string  `json:"expires_at" validate:"required,datetime" example:"2025-12-31T23:59:59Z"`
}

type UpdateRefreshTokenRequest struct {
	TokenID   string  `json:"token_id" validate:"required,uuid"`
	IPAddress *string `json:"ip_address,omitempty" example:"192.168.1.1"`
	UserAgent *string `json:"user_agent,omitempty"`
	RevokedAt *string `json:"revoked_at,omitempty" example:"2025-12-31T23:59:59Z"`
	RevokedBy *string `json:"revoked_by,omitempty" validate:"omitempty,uuid"`
}

type SignInWithEmailRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"secure.password"`
}

type SignInWithUsernameRequest struct {
	Username string `json:"username" validate:"required" example:"user123"`
	Password string `json:"password" validate:"required" example:"secure.password"`
}

type SignInResponse struct {
	AuthenticatedUser
}

type AccessTokenPayload struct {
	UserID string `json:"user_id"` // User ID
	Email  string `json:"email"`   // User Email
	SID    string `json:"sid"`     // Session ID
}

// InitiateEmailVerificationRequest represents the request payload for initiating email verification.
type InitiateEmailVerificationRequest struct {
	Email      string `json:"email" validate:"required,email" example:"user@example.com"`
	RedirectTo string `json:"redirect_to,omitempty" validate:"omitempty,url" example:"https://example.com/verified"`
}

// ValidateEmailVerificationRequest represents the request payload for validating email verification.
type ValidateEmailVerificationRequest struct {
	Token string `json:"token" validate:"required" example:"01FZ..."`
}

// RevokeEmailVerificationRequest represents the request payload for revoking email verification.
type RevokeEmailVerificationRequest struct {
	Token string `json:"token" validate:"required" example:"01FZ..."`
}

// ResendEmailVerificationRequest represents the request payload for resending email verification.
type ResendEmailVerificationRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}
