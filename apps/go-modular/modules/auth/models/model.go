package models

import (
	"net"
	"time"

	"github.com/gofrs/uuid/v5"
	user_models "go-modular/modules/user/models"
)

// -- MARK: UserPassword section

// Define table name for UserPassword model
const UserPasswordTable = "public.user_passwords"

// UserPassword represents user password model in the database
type UserPassword struct {
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	PasswordHash string     `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
}

// -- MARK: Session section

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

// -- MARK: RefreshToken section

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

// -- MARK: OneTimeToken section

// Define table name for OneTimeToken model
const OneTimeTokenTable = "public.one_time_tokens"

// OneTimeTokenSubject is an enum for the subject field in OneTimeToken
type OneTimeTokenSubject string

const (
	OneTimeTokenSubjectEmailOTP          OneTimeTokenSubject = "email_otp"
	OneTimeTokenSubjectEmailVerification OneTimeTokenSubject = "email_verification"
)

// OneTimeToken represents a one-time-use token for sensitive authentication flows.
// This table is used for various flows such as email verification, password reset,
// multi-factor authentication (MFA), and reauthentication. Each token is associated
// with a user (optional), a subject (purpose of the token), a hashed token value,
// and a "relates_to" field (such as email or phone). The token has a creation
// and expiration time, and optionally tracks the last time it was sent (for rate limiting).
//
// Metadata is a JSONB blob stored in the database and exposed as a map[string]interface{}
// to allow storing flow-specific data (e.g. ip, user agent, redirect_to, etc.).
type OneTimeToken struct {
	ID         uuid.UUID           `json:"id" db:"id"`                       // Unique identifier for the token
	UserID     *uuid.UUID          `json:"user_id" db:"user_id"`             // Optional user ID the token is associated with
	Subject    OneTimeTokenSubject `json:"subject" db:"subject"`             // Purpose/type of the token (e.g., email_otp, email_verification)
	TokenHash  string              `json:"token_hash" db:"token_hash"`       // Hashed token value (SHA256), don't store plain text or plain JWT
	RelatesTo  string              `json:"relates_to" db:"relates_to"`       // Context for the token (e.g., email, phone)
	Metadata   map[string]any      `json:"metadata,omitempty" db:"metadata"` // Arbitrary JSON metadata (stored as JSONB)
	CreatedAt  time.Time           `json:"created_at" db:"created_at"`       // When the token was created
	ExpiresAt  time.Time           `json:"expires_at" db:"expires_at"`       // When the token expires
	LastSentAt *time.Time          `json:"last_sent_at" db:"last_sent_at"`   // Last time the token was sent (for throttling)
}

// -- MARK: User credentials section

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
