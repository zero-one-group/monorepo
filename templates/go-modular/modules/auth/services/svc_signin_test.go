package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"{{ package_name | kebab_case }}/modules/auth/models"
	user_models "{{ package_name | kebab_case }}/modules/user/models"
	"{{ package_name | kebab_case }}/pkg/apputils"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const goodPassword = "correct horse battery staple"

// seedUser creates a verified user with a known password in both fakes.
func seedUser(t *testing.T, repo *fakeAuthRepo, users *fakeUserService, verified bool) *user_models.User {
	t.Helper()
	uname := "jane"
	u := &user_models.User{ID: uuid.Must(uuid.NewV7()), Email: "jane@example.com", Username: &uname, DisplayName: "Jane"}
	if verified {
		now := time.Now()
		u.EmailVerifiedAt = &now
	}
	users.users = append(users.users, u)
	svc := newTestAuthService(repo, users)
	require.NoError(t, svc.SetUserPassword(context.Background(), &models.UserPassword{UserID: u.ID, PasswordHash: goodPassword}))
	return u
}

func TestSignIn_HappyPathIssuesTokensSessionAndRefreshToken(t *testing.T) {
	repo, users := newFakeAuthRepo(), &fakeUserService{}
	u := seedUser(t, repo, users, true)
	svc := newTestAuthService(repo, users)

	authed, err := svc.SignInWithEmail(context.Background(), u.Email, goodPassword)
	require.NoError(t, err)

	assert.Equal(t, u.ID, authed.User.ID)
	assert.NotEmpty(t, authed.AccessToken)
	assert.NotEmpty(t, authed.RefreshToken)
	require.NotNil(t, authed.SessionID)
	assert.Contains(t, repo.sessions, *authed.SessionID, "a session row is created")
	assert.Len(t, repo.refreshTokens, 1, "a refresh token row is created")
	for _, rt := range repo.refreshTokens {
		assert.Equal(t, *authed.SessionID, *rt.SessionID, "refresh token is bound to the session")
		assert.Equal(t, u.ID, rt.UserID)
	}

	// The access token must carry the user, email and session id and be verifiable with our key.
	gen := apputils.NewJWTGenerator(apputils.JWTConfig{SecretKey: svc.secretKey, SigningAlg: svc.signingAlg, Issuer: svc.baseURL})
	claims, err := gen.ParseAndValidate(context.Background(), authed.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, u.ID.String(), claims["sub"])
	assert.Equal(t, u.Email, claims["email"])
	assert.Equal(t, authed.SessionID.String(), claims["sid"])
}

func TestSignIn_WithUsername(t *testing.T) {
	repo, users := newFakeAuthRepo(), &fakeUserService{}
	u := seedUser(t, repo, users, true)
	svc := newTestAuthService(repo, users)

	authed, err := svc.SignInWithUsername(context.Background(), *u.Username, goodPassword)
	require.NoError(t, err)
	assert.Equal(t, u.ID, authed.User.ID)
}

func TestSignIn_Failures(t *testing.T) {
	ctx := context.Background()

	t.Run("empty identifier or password", func(t *testing.T) {
		svc := newTestAuthService(newFakeAuthRepo(), &fakeUserService{})
		_, err := svc.SignInWithEmail(ctx, "", goodPassword)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		_, err = svc.SignInWithUsername(ctx, "jane", "")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})
	t.Run("unknown user is indistinguishable from a wrong password", func(t *testing.T) {
		svc := newTestAuthService(newFakeAuthRepo(), &fakeUserService{})
		_, err := svc.SignInWithEmail(ctx, "nobody@example.com", goodPassword)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})
	t.Run("wrong password", func(t *testing.T) {
		repo, users := newFakeAuthRepo(), &fakeUserService{}
		u := seedUser(t, repo, users, true)
		_, err := newTestAuthService(repo, users).SignInWithEmail(ctx, u.Email, "nope")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Empty(t, repo.sessions, "no session on failed sign-in")
	})
	t.Run("unverified email is rejected only after the password matched", func(t *testing.T) {
		repo, users := newFakeAuthRepo(), &fakeUserService{}
		u := seedUser(t, repo, users, false)
		svc := newTestAuthService(repo, users)
		_, err := svc.SignInWithEmail(ctx, u.Email, "nope")
		assert.ErrorIs(t, err, ErrInvalidCredentials, "wrong password wins over unverified")
		_, err = svc.SignInWithEmail(ctx, u.Email, goodPassword)
		assert.ErrorIs(t, err, ErrEmailNotVerified)
	})
	t.Run("user lookup error is propagated", func(t *testing.T) {
		users := &fakeUserService{err: errors.New("db down")}
		_, err := newTestAuthService(newFakeAuthRepo(), users).SignInWithEmail(ctx, "x@y.z", goodPassword)
		assert.EqualError(t, err, "db down")
	})
	t.Run("session creation error is propagated", func(t *testing.T) {
		repo, users := newFakeAuthRepo(), &fakeUserService{}
		u := seedUser(t, repo, users, true)
		svc := newTestAuthService(repo, users)
		repo.err = errors.New("db down") // after seeding: password check itself now fails
		_, err := svc.SignInWithEmail(ctx, u.Email, goodPassword)
		assert.EqualError(t, err, "db down")
	})
}

// The handler puts request headers on the context under apputils.HeadersContextKey so the
// service can pick the token audience from X-App-Audience. Both apps (customer PWA and
// backoffice) sign in through the same endpoint, so the audience must end up in the token.
func TestSignIn_UsesXAppAudienceFromRequestHeaders(t *testing.T) {
	repo, users := newFakeAuthRepo(), &fakeUserService{}
	u := seedUser(t, repo, users, true)
	svc := newTestAuthService(repo, users)

	ctx := context.WithValue(context.Background(), apputils.HeadersContextKey, map[string]string{"X-App-Audience": "admin-app"})
	authed, err := svc.SignInWithEmail(ctx, u.Email, goodPassword)
	require.NoError(t, err)

	gen := apputils.NewJWTGenerator(apputils.JWTConfig{SecretKey: svc.secretKey, SigningAlg: svc.signingAlg, Issuer: svc.baseURL})
	claims, err := gen.ParseAndValidate(context.Background(), authed.RefreshToken)
	require.NoError(t, err)
	assert.Contains(t, claims["aud"], "admin-app")
}

func TestSignIn_DefaultsAudienceToClientApp(t *testing.T) {
	repo, users := newFakeAuthRepo(), &fakeUserService{}
	u := seedUser(t, repo, users, true)
	svc := newTestAuthService(repo, users)

	authed, err := svc.SignInWithEmail(context.Background(), u.Email, goodPassword)
	require.NoError(t, err)
	gen := apputils.NewJWTGenerator(apputils.JWTConfig{SecretKey: svc.secretKey, SigningAlg: svc.signingAlg, Issuer: svc.baseURL})
	claims, err := gen.ParseAndValidate(context.Background(), authed.RefreshToken)
	require.NoError(t, err)
	assert.Contains(t, claims["aud"], "client-app")
}
