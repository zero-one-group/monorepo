package services

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"{{ package_name | kebab_case }}/modules/auth/models"
)

// CreateRefreshToken creates a new refresh token.
func (s *AuthService) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	if token == nil {
		return errors.New("refresh token is required")
	}
	if token.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}
	if len(token.TokenHash) == 0 {
		return errors.New("token_hash is required")
	}
	if token.ExpiresAt.IsZero() || token.ExpiresAt.Before(time.Now()) {
		return errors.New("expires_at must be set and in the future")
	}
	return s.authRepo.CreateRefreshToken(ctx, token)
}

// GetRefreshToken retrieves a refresh token by its ID.
func (s *AuthService) GetRefreshToken(ctx context.Context, tokenID uuid.UUID) (*models.RefreshToken, error) {
	if tokenID == uuid.Nil {
		return nil, errors.New("refresh_token_id is required")
	}
	return s.authRepo.GetRefreshToken(ctx, tokenID)
}

// UpdateRefreshToken updates an existing refresh token.
func (s *AuthService) UpdateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	if token == nil || token.ID == uuid.Nil {
		return errors.New("refresh token and token.ID are required")
	}
	return s.authRepo.UpdateRefreshToken(ctx, token)
}

// DeleteRefreshToken deletes a refresh token by its ID.
func (s *AuthService) DeleteRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	if tokenID == uuid.Nil {
		return errors.New("refresh_token_id is required")
	}
	return s.authRepo.DeleteRefreshToken(ctx, tokenID)
}

// ValidateRefreshToken checks if a refresh token is valid (not revoked and not expired).
func (s *AuthService) ValidateRefreshToken(ctx context.Context, tokenID uuid.UUID) (bool, error) {
	if tokenID == uuid.Nil {
		return false, errors.New("refresh_token_id is required")
	}
	return s.authRepo.ValidateRefreshToken(ctx, tokenID)
}
