package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"time"

	"go-modular/modules/auth/models"
	"go-modular/pkg/apputils"
)

// helper: try to detect if a user struct indicates the email is already verified.
// checks common field names: EmailVerified, IsEmailVerified, Verified (bool)
// and EmailVerifiedAt, VerifiedAt (time.Time or *time.Time non-zero)
func isUserEmailVerified(u any) bool {
	if u == nil {
		return false
	}
	v := reflect.ValueOf(u)
	if !v.IsValid() {
		return false
	}
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return false
	}

	// bool style fields
	for _, name := range []string{"EmailVerified", "IsEmailVerified", "Verified"} {
		f := v.FieldByName(name)
		if f.IsValid() && f.Kind() == reflect.Bool && f.Bool() {
			return true
		}
	}

	// time/timestamp style fields
	for _, name := range []string{"EmailVerifiedAt", "VerifiedAt"} {
		f := v.FieldByName(name)
		if !f.IsValid() {
			continue
		}
		// pointer to time
		if f.Kind() == reflect.Pointer {
			if !f.IsNil() {
				return true
			}
		}
		// struct (likely time.Time)
		if f.Kind() == reflect.Struct {
			if t, ok := f.Interface().(time.Time); ok && !t.IsZero() {
				return true
			}
		}
	}

	return false
}

// InitiateEmailVerification generates and stores a new email verification token for the user.
// If a valid token already exists, it only updates the last_sent_at field and does not generate a new token.
// redirectTo (optional) will be stored inside token.Metadata["redirect_to"] and, when provided,
// appended to the verification link sent to the user.
func (s *AuthService) InitiateEmailVerification(ctx context.Context, email string, redirectTo string) error {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// If user is already verified, short-circuit
	if isUserEmailVerified(user) {
		return fmt.Errorf("email already verified")
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
		// If a valid token exists, update last_sent_at and do not generate a new token
		existingToken.LastSentAt = &now
		if err := s.authRepo.UpdateOneTimeTokenLastSentAt(ctx, existingToken.ID, now); err != nil {
			return err
		}
		// Can't retrieve raw token from the DB (we only store its hash), so we cannot resend the exact token.
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

	// Prepare metadata and store the new token hash in the database
	var metadata map[string]any
	if redirectTo != "" {
		metadata = map[string]any{
			"redirect_to": redirectTo,
		}
	}

	token := &models.OneTimeToken{
		UserID:     &userID,
		Subject:    models.OneTimeTokenSubjectEmailVerification,
		TokenHash:  tokenHash,
		RelatesTo:  email,
		Metadata:   metadata,
		CreatedAt:  now,
		ExpiresAt:  expiresAt,
		LastSentAt: &now,
	}
	if err := s.authRepo.CreateOneTimeToken(ctx, token); err != nil {
		return err
	}

	// Send the rawToken to the user's email address, include redirectTo when present
	if err := s.sendVerificationEmail(ctx, email, rawToken, redirectTo); err != nil {
		// If sending fails, propagate the error (caller can decide what to do)
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// ValidateEmailVerification checks if the provided token is valid.
// It resolves the user from the stored one-time token, marks the user's email as verified,
// deletes the one-time token (one-time use) and returns true on success.
func (s *AuthService) ValidateEmailVerification(ctx context.Context, token string) (bool, error) {
	if token == "" {
		return false, errors.New("token is required")
	}

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// Retrieve the token from the database using its hash
	oneTimeToken, err := s.authRepo.GetOneTimeTokenByTokenHash(ctx, tokenHash)
	if err != nil || oneTimeToken == nil {
		return false, errors.New("invalid or expired token")
	}

	// Check expiration
	if time.Now().After(oneTimeToken.ExpiresAt) {
		return false, errors.New("token expired")
	}

	// Ensure token is bound to a user
	if oneTimeToken.UserID == nil {
		return false, errors.New("token not bound to a user")
	}
	userID := *oneTimeToken.UserID

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
		return errors.New("invalid or expired token")
	}
	return s.authRepo.DeleteOneTimeToken(ctx, oneTimeToken.ID)
}

// ResendEmailVerification either resends an existing valid token (by updating last_sent_at)
// or generates a new token if none is valid. Old tokens are revoked if a new one is created.
// redirectTo (optional) will be stored for newly created tokens.
func (s *AuthService) ResendEmailVerification(ctx context.Context, email string, redirectTo string) error {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// If user is already verified, short-circuit
	if isUserEmailVerified(user) {
		return fmt.Errorf("email already verified")
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
		// If a valid token exists, update last_sent_at and do not generate a new token
		existingToken.LastSentAt = &now
		if err := s.authRepo.UpdateOneTimeTokenLastSentAt(ctx, existingToken.ID, now); err != nil {
			return err
		}
		// Can't retrieve raw token from the DB (we only store its hash), so we cannot resend the exact token.
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

	var metadata map[string]any
	if redirectTo != "" {
		metadata = map[string]any{
			"redirect_to": redirectTo,
		}
	}

	token := &models.OneTimeToken{
		UserID:     &userID,
		Subject:    models.OneTimeTokenSubjectEmailVerification,
		TokenHash:  tokenHash,
		RelatesTo:  email,
		Metadata:   metadata,
		CreatedAt:  now,
		ExpiresAt:  expiresAt,
		LastSentAt: &now,
	}
	if err := s.authRepo.CreateOneTimeToken(ctx, token); err != nil {
		return err
	}

	// Send the new rawToken to the user's email address
	if err := s.sendVerificationEmail(ctx, email, rawToken, redirectTo); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// sendVerificationEmail constructs the verification URL and sends the email using the injected mailer.
// If no mailer is configured, it logs the URL to stdout (useful for local dev).
// redirectTo (optional) will be appended to the verification link as query parameter `redirect_to`.
func (s *AuthService) sendVerificationEmail(ctx context.Context, toEmail, rawToken, redirectTo string) error {
	// Determine base URL:
	// 1) prefer configured s.baseURL
	// 2) fallback to environment SERVER_HOST/SERVER_PORT
	// 3) final fallback to localhost:8000
	base := s.baseURL
	if base == "" {
		host := os.Getenv("SERVER_HOST")
		port := os.Getenv("SERVER_PORT")
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "8000"
		}
		base = fmt.Sprintf("http://%s:%s", host, port)
	}

	u, err := url.Parse(base)
	if err != nil || u.Scheme == "" || u.Host == "" {
		// fallback to SERVER_HOST/SERVER_PORT env vars explicitly
		host := os.Getenv("SERVER_HOST")
		port := os.Getenv("SERVER_PORT")
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "8000"
		}
		u = &url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%s", host, port)}
	}

	// Use only the token in the verification link (do NOT include the email)
	u.Path = "/api/v1/auth/verify-email"
	q := u.Query()
	q.Set("token", rawToken)
	if redirectTo != "" {
		q.Set("redirect_to", redirectTo)
	}
	u.RawQuery = q.Encode()
	verifyURL := u.String()

	// Try to fetch user to pass display name to template
	var displayName string
	if s.userService != nil {
		if user, err := s.userService.GetUserByEmail(ctx, toEmail); err == nil && user != nil {
			displayName = user.DisplayName
		}
	}

	// Template data passed to the email template; template can access .VerifyURL, .Email and .DisplayName
	data := map[string]any{
		"Email":       toEmail,
		"DisplayName": displayName,
		"VerifyURL":   verifyURL,
	}

	subject := "Verify your email address"
	templateName := "email_verification.html" // ensure this template exists in templates/emails/

	// If the service has a mailer configured, use it. Otherwise print the URL.
	// Note: AuthService struct is expected to have a mailer field injected when the service is created.
	if s.mailer != nil {
		if err := s.mailer.SendEmail(ctx, []string{toEmail}, subject, templateName, data); err != nil {
			return err
		}
		return nil
	}

	// Fallback for development: print verification link
	fmt.Println("No mailer configured, verification link for", toEmail, ":", verifyURL)
	return nil
}
