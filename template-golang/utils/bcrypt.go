package utils

import (
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash from the provided password
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password",
			slog.String("error", err.Error()),
		)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// ComparePassword compares a plaintext password with a bcrypt hash
func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		slog.Warn("Password comparison failed",
			slog.String("error", err.Error()),
		)
		return false
	}
	return true
}
