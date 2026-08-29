package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"go-modular/modules/auth/models"
	user_models "go-modular/modules/user/models"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func hashOf(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func unverifiedUser() *user_models.User {
	return &user_models.User{ID: uuid.Must(uuid.NewV7()), Email: "new@example.com", DisplayName: "New"}
}

func TestIsUserEmailVerified(t *testing.T) {
	now := time.Now()
	assert.False(t, isUserEmailVerified(nil))
	assert.False(t, isUserEmailVerified((*user_models.User)(nil)))
	assert.False(t, isUserEmailVerified("not a struct"))
	assert.False(t, isUserEmailVerified(&user_models.User{}))
	assert.True(t, isUserEmailVerified(&user_models.User{EmailVerifiedAt: &now}))
	assert.True(t, isUserEmailVerified(struct{ Verified bool }{true}))
	assert.True(t, isUserEmailVerified(struct{ VerifiedAt time.Time }{now}))
	assert.False(t, isUserEmailVerified(struct{ VerifiedAt time.Time }{}))
}

func TestInitiateEmailVerification(t *testing.T) {
	ctx := context.Background()

	t.Run("creates a hashed 15-minute token bound to the user with redirect metadata", func(t *testing.T) {
		repo, users := newFakeAuthRepo(), &fakeUserService{users: []*user_models.User{unverifiedUser()}}
		svc := newTestAuthService(repo, users)

		require.NoError(t, svc.InitiateEmailVerification(ctx, "new@example.com", "https://app.example.com/ok"))

		require.Len(t, repo.oneTimeTokens, 1)
		for _, tok := range repo.oneTimeTokens {
			assert.Equal(t, models.OneTimeTokenSubjectEmailVerification, tok.Subject)
			assert.Equal(t, users.users[0].ID, *tok.UserID)
			assert.Equal(t, "new@example.com", tok.RelatesTo)
			assert.Len(t, tok.TokenHash, 64, "sha256 hex; the raw token is never stored")
			assert.Equal(t, "https://app.example.com/ok", tok.Metadata["redirect_to"])
			assert.WithinDuration(t, time.Now().Add(15*time.Minute), tok.ExpiresAt, 5*time.Second)
			assert.NotNil(t, tok.LastSentAt)
		}
	})
	t.Run("no metadata when there is no redirect", func(t *testing.T) {
		repo, users := newFakeAuthRepo(), &fakeUserService{users: []*user_models.User{unverifiedUser()}}
		require.NoError(t, newTestAuthService(repo, users).InitiateEmailVerification(ctx, "new@example.com", ""))
		for _, tok := range repo.oneTimeTokens {
			assert.Nil(t, tok.Metadata)
		}
	})
	t.Run("a still-valid token is only re-stamped, not regenerated", func(t *testing.T) {
		repo, users := newFakeAuthRepo(), &fakeUserService{users: []*user_models.User{unverifiedUser()}}
		svc := newTestAuthService(repo, users)
		require.NoError(t, svc.InitiateEmailVerification(ctx, "new@example.com", ""))
		require.NoError(t, svc.InitiateEmailVerification(ctx, "new@example.com", ""))
		assert.Len(t, repo.oneTimeTokens, 1)
		assert.Len(t, repo.resent, 1)
	})
	t.Run("unknown or already verified users", func(t *testing.T) {
		svc := newTestAuthService(newFakeAuthRepo(), &fakeUserService{err: errors.New("no rows")})
		assert.EqualError(t, svc.InitiateEmailVerification(ctx, "x@example.com", ""), "user not found")

		now := time.Now()
		verified := &user_models.User{ID: uuid.Must(uuid.NewV7()), Email: "done@example.com", EmailVerifiedAt: &now}
		svc = newTestAuthService(newFakeAuthRepo(), &fakeUserService{users: []*user_models.User{verified}})
		assert.EqualError(t, svc.InitiateEmailVerification(ctx, "done@example.com", ""), "email already verified")
	})
	t.Run("repository failure while storing the token", func(t *testing.T) {
		repo, users := newFakeAuthRepo(), &fakeUserService{users: []*user_models.User{unverifiedUser()}}
		repo.err = errors.New("db down")
		assert.EqualError(t, newTestAuthService(repo, users).InitiateEmailVerification(ctx, "new@example.com", ""), "db down")
	})
}

func TestValidateEmailVerification(t *testing.T) {
	ctx := context.Background()
	seedToken := func(repo *fakeAuthRepo, raw string, userID *uuid.UUID, expires time.Time) *models.OneTimeToken {
		tok := &models.OneTimeToken{ID: uuid.Must(uuid.NewV7()), UserID: userID, Subject: models.OneTimeTokenSubjectEmailVerification, TokenHash: hashOf(raw), RelatesTo: "new@example.com", ExpiresAt: expires}
		repo.oneTimeTokens[tok.ID] = tok
		return tok
	}

	t.Run("valid token verifies the user and is consumed", func(t *testing.T) {
		repo, users := newFakeAuthRepo(), &fakeUserService{}
		uid := uuid.Must(uuid.NewV7())
		tok := seedToken(repo, "raw-token", &uid, time.Now().Add(time.Minute))

		ok, err := newTestAuthService(repo, users).ValidateEmailVerification(ctx, "raw-token")
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, []uuid.UUID{uid}, users.verified)
		assert.Equal(t, []uuid.UUID{tok.ID}, repo.deletedOTT, "one-time use")

		ok, err = newTestAuthService(repo, users).ValidateEmailVerification(ctx, "raw-token")
		assert.False(t, ok)
		assert.EqualError(t, err, "invalid or expired token", "replaying the same token fails")
	})
	t.Run("rejections", func(t *testing.T) {
		repo, users := newFakeAuthRepo(), &fakeUserService{}
		svc := newTestAuthService(repo, users)
		uid := uuid.Must(uuid.NewV7())

		_, err := svc.ValidateEmailVerification(ctx, "")
		assert.EqualError(t, err, "token is required")
		_, err = svc.ValidateEmailVerification(ctx, "unknown")
		assert.EqualError(t, err, "invalid or expired token")

		seedToken(repo, "expired", &uid, time.Now().Add(-time.Minute))
		_, err = svc.ValidateEmailVerification(ctx, "expired")
		assert.EqualError(t, err, "token expired")

		seedToken(repo, "orphan", nil, time.Now().Add(time.Minute))
		_, err = svc.ValidateEmailVerification(ctx, "orphan")
		assert.EqualError(t, err, "token not bound to a user")

		users.err = errors.New("db down")
		seedToken(repo, "good", &uid, time.Now().Add(time.Minute))
		_, err = svc.ValidateEmailVerification(ctx, "good")
		assert.EqualError(t, err, "db down")
	})
}

func TestRevokeEmailVerification(t *testing.T) {
	ctx := context.Background()
	repo := newFakeAuthRepo()
	svc := newTestAuthService(repo, &fakeUserService{})
	tok := &models.OneTimeToken{ID: uuid.Must(uuid.NewV7()), TokenHash: hashOf("raw"), Subject: models.OneTimeTokenSubjectEmailVerification, ExpiresAt: time.Now().Add(time.Minute)}
	repo.oneTimeTokens[tok.ID] = tok

	require.NoError(t, svc.RevokeEmailVerification(ctx, "raw"))
	assert.Equal(t, []uuid.UUID{tok.ID}, repo.deletedOTT)
	assert.EqualError(t, svc.RevokeEmailVerification(ctx, "raw"), "invalid or expired token")
}

func TestResendEmailVerification(t *testing.T) {
	ctx := context.Background()

	t.Run("re-stamps a still-valid token", func(t *testing.T) {
		repo, users := newFakeAuthRepo(), &fakeUserService{users: []*user_models.User{unverifiedUser()}}
		svc := newTestAuthService(repo, users)
		require.NoError(t, svc.InitiateEmailVerification(ctx, "new@example.com", ""))
		require.NoError(t, svc.ResendEmailVerification(ctx, "new@example.com", ""))
		assert.Len(t, repo.oneTimeTokens, 1)
		assert.Len(t, repo.resent, 1)
	})
	t.Run("revokes expired tokens and issues a fresh one", func(t *testing.T) {
		repo, users := newFakeAuthRepo(), &fakeUserService{users: []*user_models.User{unverifiedUser()}}
		uid := users.users[0].ID
		stale := &models.OneTimeToken{ID: uuid.Must(uuid.NewV7()), UserID: &uid, Subject: models.OneTimeTokenSubjectEmailVerification, TokenHash: hashOf("old"), RelatesTo: "new@example.com", ExpiresAt: time.Now().Add(-time.Hour)}
		repo.oneTimeTokens[stale.ID] = stale

		require.NoError(t, newTestAuthService(repo, users).ResendEmailVerification(ctx, "new@example.com", "https://app.example.com/ok"))

		assert.Equal(t, []uuid.UUID{stale.ID}, repo.deletedOTT)
		require.Len(t, repo.oneTimeTokens, 1)
		for _, tok := range repo.oneTimeTokens {
			assert.NotEqual(t, stale.ID, tok.ID)
			assert.Equal(t, "https://app.example.com/ok", tok.Metadata["redirect_to"])
		}
	})
	t.Run("unknown / verified / repo failure", func(t *testing.T) {
		assert.EqualError(t, newTestAuthService(newFakeAuthRepo(), &fakeUserService{err: errors.New("x")}).ResendEmailVerification(ctx, "a@b.c", ""), "user not found")
		now := time.Now()
		verified := &user_models.User{ID: uuid.Must(uuid.NewV7()), Email: "done@example.com", EmailVerifiedAt: &now}
		assert.EqualError(t, newTestAuthService(newFakeAuthRepo(), &fakeUserService{users: []*user_models.User{verified}}).ResendEmailVerification(ctx, "done@example.com", ""), "email already verified")
		repo := newFakeAuthRepo()
		repo.err = errors.New("db down")
		assert.EqualError(t, newTestAuthService(repo, &fakeUserService{users: []*user_models.User{unverifiedUser()}}).ResendEmailVerification(ctx, "new@example.com", ""), "db down")
	})
}

// Without a mailer the service prints the link (dev fallback); the link must carry the token
// and redirect_to and point at the configured base URL, never at the email address.
func TestSendVerificationEmail_LinkShape(t *testing.T) {
	repo, users := newFakeAuthRepo(), &fakeUserService{users: []*user_models.User{unverifiedUser()}}
	svc := newTestAuthService(repo, users)
	require.NoError(t, svc.sendVerificationEmail(context.Background(), "new@example.com", "RAW", "https://app.example.com/ok"))

	svc.baseURL = "not a url"
	t.Setenv("SERVER_HOST", "")
	t.Setenv("SERVER_PORT", "")
	require.NoError(t, svc.sendVerificationEmail(context.Background(), "new@example.com", "RAW", ""))

	svc.baseURL = ""
	t.Setenv("SERVER_HOST", "api.internal")
	t.Setenv("SERVER_PORT", "9000")
	require.NoError(t, svc.sendVerificationEmail(context.Background(), "new@example.com", "RAW", ""))
}
