package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"go-modular/modules/auth/models"
	"go-modular/pkg/apputils"
)

// InitiateEmailVerification generates and stores a new email verification token for the user.
// If a valid token already exists, it only updates the last_sent_at field and does not generate a new token.
// Returns an error if the user is not found or if token generation/storage fails.
func (s *AuthService) InitiateEmailVerification(ctx context.Context, email string) error {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	userID := user.ID

	// Check for an existing valid token for this user/email
	tokens, err := s.authRepo.FindAllOneTimeTokens(ctx)
	now := time.Now()
	var existingToken *models.OneTimeToken
	if err == nil {
		for _, t := range tokens {
			if t.UserID != nil && *t.UserID == userID &&
				t.Subject == models.OneTimeTokenSubjectEmailVerification &&
				t.RelatesTo == email &&
				now.Before(t.ExpiresAt) {
				existingToken = t
				break
			}
		}
	}

	if existingToken != nil {
		// If a valid token exists, update last_sent_at and resend (do not generate a new token)
		existingToken.LastSentAt = &now
		if err := s.authRepo.UpdateOneTimeTokenLastSentAt(ctx, existingToken.ID, now); err != nil {
			return err
		}
		// TODO: Resend the token to the user's email (raw token is not available, only hash)
		// sendVerificationEmail(email, <cannot retrieve raw token, only hash>)
		fmt.Println("Verification email resent to:", email)
		return nil
	}

	// Generate a new, cryptographically secure, URL-safe token (length 48, includes unix timestamp)
	rawToken, err := apputils.GenerateURLSafeToken(48)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])
	expiresAt := now.Add(15 * time.Minute)

	// Remove any old tokens for this user/email
	for _, t := range tokens {
		if t.UserID != nil && *t.UserID == userID && t.Subject == models.OneTimeTokenSubjectEmailVerification {
			_ = s.authRepo.DeleteOneTimeToken(ctx, t.ID)
		}
	}

	// Store the new token hash in the database
	token := &models.OneTimeToken{
		UserID:     &userID,
		Subject:    models.OneTimeTokenSubjectEmailVerification,
		TokenHash:  tokenHash,
		RelatesTo:  email,
		CreatedAt:  now,
		ExpiresAt:  expiresAt,
		LastSentAt: &now,
	}
	if err := s.authRepo.CreateOneTimeToken(ctx, token); err != nil {
		return err
	}

	// TODO: Send the rawToken to the user's email address
	// sendVerificationEmail(email, rawToken)
	fmt.Println("Verification email sent to:", email)
	fmt.Println("Raw token (for debugging):", rawToken)

	return nil
}

// ValidateEmailVerification checks if the provided token is valid for the user and email.
// If valid, marks the user's email as verified and deletes the token (one-time use).
func (s *AuthService) ValidateEmailVerification(ctx context.Context, email, token string) (bool, error) {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("user not found")
	}
	userID := user.ID

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// Retrieve the token from the database using its hash
	oneTimeToken, err := s.authRepo.GetOneTimeTokenByTokenHash(ctx, tokenHash)
	if err != nil || oneTimeToken == nil {
		return false, errors.New("invalid or expired token")
	}
	if oneTimeToken.UserID == nil || *oneTimeToken.UserID != userID {
		return false, errors.New("token does not belong to user")
	}
	if time.Now().After(oneTimeToken.ExpiresAt) {
		return false, errors.New("token expired")
	}

	// Delete the token after successful validation (one-time use)
	_ = s.authRepo.DeleteOneTimeToken(ctx, oneTimeToken.ID)

	// Mark the user's email as verified
	if err := s.userService.MarkEmailVerified(ctx, userID); err != nil {
		return false, err
	}

	return true, nil
}

// RevokeEmailVerification deletes the email verification token for the user, using only the token.
// Returns an error if the token is not found.
func (s *AuthService) RevokeEmailVerification(ctx context.Context, token string) error {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// Retrieve the token by its hash
	oneTimeToken, err := s.authRepo.GetOneTimeTokenByTokenHash(ctx, tokenHash)
	if err != nil || oneTimeToken == nil {
		return errors.New("token not found")
	}
	return s.authRepo.DeleteOneTimeToken(ctx, oneTimeToken.ID)
}

// ResendEmailVerification either resends an existing valid token (by updating last_sent_at)
// or generates a new token if none is valid. Old tokens are revoked if a new one is created.
// Returns an error if the user is not found or token generation/storage fails.
func (s *AuthService) ResendEmailVerification(ctx context.Context, email string) error {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	userID := user.ID

	// Check for an existing valid token for this user/email
	tokens, err := s.authRepo.FindAllOneTimeTokens(ctx)
	now := time.Now()
	var existingToken *models.OneTimeToken
	if err == nil {
		for _, t := range tokens {
			if t.UserID != nil && *t.UserID == userID &&
				t.Subject == models.OneTimeTokenSubjectEmailVerification &&
				t.RelatesTo == email &&
				now.Before(t.ExpiresAt) {
				existingToken = t
				break
			}
		}
	}

	if existingToken != nil {
		// If a valid token exists, update last_sent_at and resend (do not generate a new token)
		existingToken.LastSentAt = &now
		if err := s.authRepo.UpdateOneTimeTokenLastSentAt(ctx, existingToken.ID, now); err != nil {
			return err
		}
		// TODO: Resend the token to the user's email (raw token is not available, only hash)
		// sendVerificationEmail(email, <cannot retrieve raw token, only hash>)
		fmt.Println("Verification email resent to:", email)
		return nil
	}

	// If no valid token exists, revoke all old tokens and generate a new one
	for _, t := range tokens {
		if t.UserID != nil && *t.UserID == userID &&
			t.Subject == models.OneTimeTokenSubjectEmailVerification &&
			t.RelatesTo == email {
			_ = s.authRepo.DeleteOneTimeToken(ctx, t.ID)
		}
	}

	// Generate a new, cryptographically secure, URL-safe token (length 48, includes unix timestamp)
	rawToken, err := apputils.GenerateURLSafeToken(48)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])
	expiresAt := now.Add(15 * time.Minute)

	token := &models.OneTimeToken{
		UserID:     &userID,
		Subject:    models.OneTimeTokenSubjectEmailVerification,
		TokenHash:  tokenHash,
		RelatesTo:  email,
		CreatedAt:  now,
		ExpiresAt:  expiresAt,
		LastSentAt: &now,
	}
	if err := s.authRepo.CreateOneTimeToken(ctx, token); err != nil {
		return err
	}

	// TODO: Send the new rawToken to the user's email address
	// sendVerificationEmail(email, rawToken)
	fmt.Println("Verification email sent to:", email)
	fmt.Println("Raw token (for debugging):", rawToken)

	return nil
}
