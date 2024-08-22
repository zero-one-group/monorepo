package model

// User represents a user in the system
type User struct {
	ID               int64   `json:"id"`                 // Unique identifier for the user
	Email            string  `json:"email"`              // User's email address
	FirstName        string  `json:"first_name"`         // User's first name
	LastName         *string `json:"last_name"`          // User's last name
	AvatarURL        *string `json:"avatar_url"`         // URL to user's avatar
	EmailConfirmedAt *int64  `json:"email_confirmed_at"` // Timestamp when email was confirmed
	LastSeenAt       *int64  `json:"last_seen_at"`       // Timestamp of last user activity
	BannedUntil      *int64  `json:"banned_until"`       // Timestamp when the user is banned until
	CreatedAt        int64   `json:"created_at"`         // Timestamp when the user was created
	UpdatedAt        *int64  `json:"updated_at"`         // Timestamp when the user was last updated
}
