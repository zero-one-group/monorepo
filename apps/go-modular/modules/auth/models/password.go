package models

import (
	"time"

	"github.com/gofrs/uuid/v5"
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
