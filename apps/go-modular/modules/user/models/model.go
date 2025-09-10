// Package models contains struct definitions related to the database layer for the user module.
// This file defines models that map to database tables and are used for database operations.

package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

// Define table name for User model
const UserTable = "public.users"

// User represents user model in the database
type User struct {
	ID          uuid.UUID     `json:"id" db:"id"`
	DisplayName string        `json:"display_name" db:"display_name"`
	Email       string        `json:"email" db:"email"`
	Username    *string       `json:"username" db:"username"`
	AvatarURL   *string       `json:"avatar_url" db:"avatar_url"`
	Metadata    *UserMetadata `json:"metadata" db:"metadata"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time    `json:"updated_at" db:"updated_at"`
	LastLoginAt *time.Time    `json:"last_login_at" db:"last_login_at"`
	BannedAt    *time.Time    `json:"banned_at" db:"banned_at"`
	BanExpires  *time.Time    `json:"ban_expires" db:"ban_expires"`
	BanReason   *string       `json:"ban_reason" db:"ban_reason"`
}

type UserMetadata struct {
	Timezone string `json:"timezone,omitempty"`
	// Add more fields as needed
}

type FilterUser struct {
	Search *string `json:"search,omitempty" query:"search"`
	Limit  int     `json:"limit,omitempty" query:"limit"`
	Offset int     `json:"offset,omitempty" query:"offset"`
}

type UserWithCredential struct {
	User
	PasswordHash []byte `json:"password_hash" db:"-"`
}

// ParseUserID parses a string to uuid.UUID, returns error if invalid.
func ParseUserID(s string) (uuid.UUID, error) {
	return uuid.FromString(s)
}

func (u *User) GetID() uuid.UUID {
	return u.ID
}
func (u *User) GetEmail() string {
	return u.Email
}
func (u *User) AsUserModel() User {
	return *u
}
