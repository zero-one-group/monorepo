package models

import (
	"net"
	"time"

	"github.com/gofrs/uuid/v5"
	user_models "go-modular/modules/user/models"
)

// Define table name for UserPassword model
const UserPasswordTable = "public.user_passwords"

// UserPassword represents user password model in the database
type UserPassword struct {
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	PasswordHash string     `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
}

// Define table name for Session model
const SessionTable = "public.sessions"

// Session represents session model in the database
type Session struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	UserID            uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash         string     `json:"token_hash" db:"token_hash"`
	UserAgent         *string    `json:"user_agent" db:"user_agent"`
	DeviceName        *string    `json:"device_name" db:"device_name"`
	DeviceFingerprint *string    `json:"device_fingerprint" db:"device_fingerprint"`
	IPAddress         *net.IP    `json:"ip_address" db:"ip_address"`
	ExpiresAt         time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	RefreshedAt       *time.Time `json:"refreshed_at" db:"refreshed_at"`
	RevokedAt         *time.Time `json:"revoked_at" db:"revoked_at"`
	RevokedBy         *uuid.UUID `json:"revoked_by" db:"revoked_by"`
}

// Define table name for RefreshToken model
const RefreshTokenTable = "public.refresh_tokens"

// RefreshToken represents refresh token model in the database
type RefreshToken struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	SessionID *uuid.UUID `json:"session_id" db:"session_id"`
	TokenHash []byte     `json:"token_hash" db:"token_hash"`
	IPAddress *net.IP    `json:"ip_address" db:"ip_address"`
	UserAgent *string    `json:"user_agent" db:"user_agent"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	RevokedAt *time.Time `json:"revoked_at" db:"revoked_at"`
	RevokedBy *uuid.UUID `json:"revoked_by" db:"revoked_by"`
}

type UserWithCredentials struct {
	User         user_models.User `json:"user"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
}

type AuthenticatedUser struct {
	UserWithCredentials
	SessionID   *uuid.UUID `json:"session_id"`
	TokenExpiry time.Time  `json:"token_expiry"`
}
