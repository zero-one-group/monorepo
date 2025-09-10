package services

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"go-modular/modules/auth/models"
)

// CreateSession creates a new session.
func (s *AuthService) CreateSession(ctx context.Context, session *models.Session) error {
	if session == nil {
		return errors.New("session is required")
	}
	if session.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}
	if session.TokenHash == "" {
		return errors.New("token_hash is required")
	}
	if session.ExpiresAt.IsZero() || session.ExpiresAt.Before(time.Now()) {
		return errors.New("expires_at must be set and in the future")
	}
	return s.authRepo.CreateSession(ctx, session)
}

// GetSession retrieves a session by its ID.
func (s *AuthService) GetSession(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	if sessionID == uuid.Nil {
		return nil, errors.New("session_id is required")
	}
	return s.authRepo.GetSession(ctx, sessionID)
}

// UpdateSession updates an existing session.
func (s *AuthService) UpdateSession(ctx context.Context, session *models.Session) error {
	if session == nil || session.ID == uuid.Nil {
		return errors.New("session and session.ID are required")
	}
	return s.authRepo.UpdateSession(ctx, session)
}

// DeleteSession deletes a session by its ID.
func (s *AuthService) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	if sessionID == uuid.Nil {
		return errors.New("session_id is required")
	}
	return s.authRepo.DeleteSession(ctx, sessionID)
}

// ValidateSession checks if a session is valid (not revoked and not expired).
func (s *AuthService) ValidateSession(ctx context.Context, sessionID uuid.UUID) (bool, error) {
	if sessionID == uuid.Nil {
		return false, errors.New("session_id is required")
	}
	return s.authRepo.ValidateSession(ctx, sessionID)
}
