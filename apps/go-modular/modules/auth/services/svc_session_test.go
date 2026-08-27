package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-modular/modules/auth/models"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validSession() *models.Session {
	return &models.Session{UserID: uuid.Must(uuid.NewV7()), TokenHash: "hash", ExpiresAt: time.Now().Add(time.Hour)}
}

func validRefreshToken() *models.RefreshToken {
	return &models.RefreshToken{ID: uuid.Must(uuid.NewV7()), UserID: uuid.Must(uuid.NewV7()), TokenHash: []byte("hash"), ExpiresAt: time.Now().Add(time.Hour)}
}

func TestCreateSession_Validation(t *testing.T) {
	svc := newTestAuthService(newFakeAuthRepo(), &fakeUserService{})
	ctx := context.Background()
	cases := []struct {
		name string
		in   *models.Session
		want string
	}{
		{"nil", nil, "session is required"},
		{"missing user", &models.Session{TokenHash: "h", ExpiresAt: time.Now().Add(time.Hour)}, "user_id is required"},
		{"missing hash", &models.Session{UserID: uuid.Must(uuid.NewV7()), ExpiresAt: time.Now().Add(time.Hour)}, "token_hash is required"},
		{"zero expiry", &models.Session{UserID: uuid.Must(uuid.NewV7()), TokenHash: "h"}, "expires_at must be set and in the future"},
		{"past expiry", &models.Session{UserID: uuid.Must(uuid.NewV7()), TokenHash: "h", ExpiresAt: time.Now().Add(-time.Minute)}, "expires_at must be set and in the future"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { assert.EqualError(t, svc.CreateSession(ctx, c.in), c.want) })
	}
}

func TestSessionLifecycle(t *testing.T) {
	repo := newFakeAuthRepo()
	svc := newTestAuthService(repo, &fakeUserService{})
	ctx := context.Background()

	s := validSession()
	require.NoError(t, svc.CreateSession(ctx, s))
	require.NotEqual(t, uuid.Nil, s.ID)

	got, err := svc.GetSession(ctx, s.ID)
	require.NoError(t, err)
	assert.Equal(t, s.TokenHash, got.TokenHash)

	repo.validSessions[s.ID] = true
	ok, err := svc.ValidateSession(ctx, s.ID)
	require.NoError(t, err)
	assert.True(t, ok)

	now := time.Now()
	s.RevokedAt = &now
	require.NoError(t, svc.UpdateSession(ctx, s))
	assert.NotNil(t, repo.sessions[s.ID].RevokedAt)

	require.NoError(t, svc.DeleteSession(ctx, s.ID))
	assert.Equal(t, []uuid.UUID{s.ID}, repo.deletedSession)
}

func TestSession_NilIDGuards(t *testing.T) {
	svc := newTestAuthService(newFakeAuthRepo(), &fakeUserService{})
	ctx := context.Background()
	_, err := svc.GetSession(ctx, uuid.Nil)
	assert.EqualError(t, err, "session_id is required")
	assert.EqualError(t, svc.UpdateSession(ctx, nil), "session and session.ID are required")
	assert.EqualError(t, svc.UpdateSession(ctx, &models.Session{}), "session and session.ID are required")
	assert.EqualError(t, svc.DeleteSession(ctx, uuid.Nil), "session_id is required")
	ok, err := svc.ValidateSession(ctx, uuid.Nil)
	assert.False(t, ok)
	assert.EqualError(t, err, "session_id is required")
}

func TestCreateRefreshToken_Validation(t *testing.T) {
	svc := newTestAuthService(newFakeAuthRepo(), &fakeUserService{})
	ctx := context.Background()
	uid := uuid.Must(uuid.NewV7())
	cases := []struct {
		name string
		in   *models.RefreshToken
		want string
	}{
		{"nil", nil, "refresh token is required"},
		{"missing user", &models.RefreshToken{TokenHash: []byte("h"), ExpiresAt: time.Now().Add(time.Hour)}, "user_id is required"},
		{"missing hash", &models.RefreshToken{UserID: uid, ExpiresAt: time.Now().Add(time.Hour)}, "token_hash is required"},
		{"past expiry", &models.RefreshToken{UserID: uid, TokenHash: []byte("h"), ExpiresAt: time.Now().Add(-time.Minute)}, "expires_at must be set and in the future"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) { assert.EqualError(t, svc.CreateRefreshToken(ctx, c.in), c.want) })
	}
}

func TestRefreshTokenLifecycle(t *testing.T) {
	repo := newFakeAuthRepo()
	svc := newTestAuthService(repo, &fakeUserService{})
	ctx := context.Background()

	tok := validRefreshToken()
	require.NoError(t, svc.CreateRefreshToken(ctx, tok))

	got, err := svc.GetRefreshToken(ctx, tok.ID)
	require.NoError(t, err)
	assert.Equal(t, tok.UserID, got.UserID)

	repo.validTokens[tok.ID] = true
	ok, err := svc.ValidateRefreshToken(ctx, tok.ID)
	require.NoError(t, err)
	assert.True(t, ok)

	now := time.Now()
	tok.RevokedAt = &now
	require.NoError(t, svc.UpdateRefreshToken(ctx, tok))
	assert.NotNil(t, repo.refreshTokens[tok.ID].RevokedAt)

	require.NoError(t, svc.DeleteRefreshToken(ctx, tok.ID))
	assert.Equal(t, []uuid.UUID{tok.ID}, repo.deletedTokens)
}

func TestRefreshToken_NilIDGuardsAndRepoErrors(t *testing.T) {
	repo := newFakeAuthRepo()
	svc := newTestAuthService(repo, &fakeUserService{})
	ctx := context.Background()
	_, err := svc.GetRefreshToken(ctx, uuid.Nil)
	assert.EqualError(t, err, "refresh_token_id is required")
	assert.EqualError(t, svc.UpdateRefreshToken(ctx, nil), "refresh token and token.ID are required")
	assert.EqualError(t, svc.UpdateRefreshToken(ctx, &models.RefreshToken{}), "refresh token and token.ID are required")
	assert.EqualError(t, svc.DeleteRefreshToken(ctx, uuid.Nil), "refresh_token_id is required")
	ok, err := svc.ValidateRefreshToken(ctx, uuid.Nil)
	assert.False(t, ok)
	assert.EqualError(t, err, "refresh_token_id is required")

	repo.err = errors.New("db down")
	assert.EqualError(t, svc.CreateRefreshToken(ctx, validRefreshToken()), "db down")
	assert.EqualError(t, svc.CreateSession(ctx, validSession()), "db down")
}
