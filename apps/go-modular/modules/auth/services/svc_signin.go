package services

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"go-modular/modules/auth/models"
	user_models "go-modular/modules/user/models"
	"go-modular/pkg/apputils"
)

// ErrInvalidCredentials is returned when authentication fails.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrEmailNotVerified is returned when the user's email is not verified.
var ErrEmailNotVerified = errors.New("email not verified")

// UserIdentity interface for user abstraction in sign-in
type UserIdentity interface {
	GetID() uuid.UUID
	GetEmail() string
	AsUserModel() user_models.User
}

// signinWithCredentials is a reusable function for both email and username sign-in.
// It validates user credentials, checks if the email is verified, and issues JWT tokens.
func (s *AuthService) signinWithCredentials(
	ctx context.Context,
	identifier string,
	password string,
	getUser func(context.Context, string) (UserIdentity, error),
) (*models.AuthenticatedUser, error) {
	if identifier == "" || password == "" {
		return nil, ErrInvalidCredentials
	}
	user, err := getUser(ctx, identifier)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Validate the user's password FIRST
	ok, err := s.ValidateUserPassword(ctx, user.GetID(), password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrInvalidCredentials
	}

	// Then check if the user's email is verified
	if u, ok := any(user).(interface{ GetEmailVerifiedAt() *time.Time }); ok {
		if u.GetEmailVerifiedAt() == nil {
			return nil, ErrEmailNotVerified
		}
	}

	// Prepare JWT generator with configuration from environment or service fields.
	jwtGen := apputils.NewJWTGenerator(apputils.JWTConfig{
		SecretKey:          s.secretKey,
		AccessTokenExpiry:  s.accessTokenExpiry,
		RefreshTokenExpiry: s.refreshTokenExpiry,
		SigningAlg:         s.signingAlg,
		Issuer:             s.baseURL,
	})

	// Determine audience for the token, default to "client-app"
	audience := "client-app"
	if md, ok := ctx.Value("headers").(map[string]string); ok {
		if aud, exists := md["X-App-Audience"]; exists && aud != "" {
			audience = aud
		}
	}

	// Generate a new UUID for the refresh token
	refreshTokenUUID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	refreshTokenID := refreshTokenUUID.String()

	// Generate the refresh token JWT and its hash
	refreshToken, err := jwtGen.GenerateRefreshTokenJWT(ctx, user.GetID().String(), audience, refreshTokenID)
	if err != nil {
		return nil, err
	}
	refreshTokenHash := jwtGen.GetHash(refreshToken)

	// Create a new session for the user
	session := &models.Session{
		UserID:    user.GetID(),
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(jwtGen.AccessTokenExpiry()),
	}
	if err := s.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	// Prepare the access token payload, including the session ID
	accessPayload := models.AccessTokenPayload{
		UserID: user.GetID().String(),
		Email:  user.GetEmail(),
		SID:    session.ID.String(),
	}
	accessToken, err := jwtGen.Sign(ctx, accessPayload, user.GetID().String())
	if err != nil {
		return nil, err
	}

	// Store the refresh token in the database
	refreshTokenModel := &models.RefreshToken{
		ID:        refreshTokenUUID,
		UserID:    user.GetID(),
		SessionID: &session.ID,
		TokenHash: []byte(refreshTokenHash),
		ExpiresAt: time.Now().Add(jwtGen.RefreshTokenExpiry()),
	}
	if err := s.CreateRefreshToken(ctx, refreshTokenModel); err != nil {
		return nil, err
	}

	// Return the authenticated user with tokens and session info
	authUser := &models.AuthenticatedUser{
		UserWithCredentials: models.UserWithCredentials{
			User:         user.AsUserModel(), // Cast to user_models.User
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		SessionID:   &session.ID,
		TokenExpiry: session.ExpiresAt,
	}

	return authUser, nil
}

// SignInWithEmail authenticates a user by email and password.
func (s *AuthService) SignInWithEmail(ctx context.Context, email, password string) (*models.AuthenticatedUser, error) {
	return s.signinWithCredentials(ctx, email, password, func(ctx context.Context, email string) (UserIdentity, error) {
		user, err := s.userService.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, nil
		}
		return user, nil
	})
}

// SignInWithUsername authenticates a user by username and password.
func (s *AuthService) SignInWithUsername(ctx context.Context, username, password string) (*models.AuthenticatedUser, error) {
	return s.signinWithCredentials(ctx, username, password, func(ctx context.Context, username string) (UserIdentity, error) {
		user, err := s.userService.GetUserByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, nil
		}
		return user, nil
	})
}
