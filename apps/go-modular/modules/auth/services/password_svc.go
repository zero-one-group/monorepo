package services

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
	"go-modular/modules/auth/models"
	"go-modular/pkg/apputils"
)

// SetUserPassword creates a new user password (with hashing).
func (s *AuthService) SetUserPassword(ctx context.Context, userPassword *models.UserPassword) error {
	if userPassword == nil || userPassword.PasswordHash == "" {
		return errors.New("password is required")
	}
	hasher := apputils.NewPasswordHasher()
	hashed, err := hasher.Hash(userPassword.PasswordHash)
	if err != nil {
		return err
	}
	userPassword.PasswordHash = hashed
	return s.authRepo.SetUserPassword(ctx, userPassword)
}

// UpdateUserPassword updates an existing user password (with current password validation and hashing).
func (s *AuthService) UpdateUserPassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	if newPassword == "" {
		return errors.New("new password is required")
	}
	if currentPassword == "" {
		return errors.New("current password is required")
	}

	// Validate current password
	ok, err := s.authRepo.ValidateUserPassword(ctx, userID, currentPassword)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("current password is incorrect")
	}

	hasher := apputils.NewPasswordHasher()
	hashed, err := hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	return s.authRepo.UpdateUserPassword(ctx, userID, hashed)
}

// ValidateUserPassword checks if the provided password matches the user's current password.
func (s *AuthService) ValidateUserPassword(ctx context.Context, userID uuid.UUID, password string) (bool, error) {
	if password == "" {
		return false, errors.New("password is required")
	}
	return s.authRepo.ValidateUserPassword(ctx, userID, password)
}
