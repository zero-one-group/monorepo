package repository

import (
	"context"
	"io"
	"testing"
	"time"

	"go-modular/modules/auth/models"
	"go-modular/pkg/testutils"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
}

func setupRefreshRepo(t *testing.T) (*AuthRepository, uuid.UUID, func()) {
	t.Helper()
	te := testutils.NewTestEnv(t)
	pool, _, err := te.SetupPostgres()
	require.NoError(t, err)
	require.NotNil(t, pool)

	te.SetupConfig()
	te.RunAppMigrations()

	// try to find a seeded user created by user_factory
	var uid uuid.UUID
	err = pool.QueryRow(context.Background(), `SELECT id FROM public.users WHERE username = $1 LIMIT 1`, "admin").Scan(&uid)
	if err != nil {
		err = pool.QueryRow(context.Background(), `SELECT id FROM public.users WHERE username = $1 LIMIT 1`, "johndoe").Scan(&uid)
		require.NoError(t, err, "failed to find seeded user 'admin' or 'johndoe' required by tests")
	}

	repo := NewAuthRepository(pool, testLogger())

	teardown := func() {
		_, _ = pool.Exec(context.Background(), `TRUNCATE TABLE `+models.RefreshTokenTable+` CASCADE`)
		pool.Close()
	}

	return repo, uid, teardown
}

func TestRefreshTokenRepo_CRUD_Validate(t *testing.T) {
	ctx := context.Background()
	repo, uid, teardown := setupRefreshRepo(t)
	defer teardown()

	// create
	now := time.Now().UTC().Truncate(time.Second)

	// prepare pointer values and byte token hash to match model types
	// sessionID omitted to avoid FK constraint — tests don't require a real session row
	tokenHash := []byte("rthash-" + uuid.Must(uuid.NewV7()).String())
	userAgent := "test-agent"

	rt := &models.RefreshToken{
		ID:        uuid.Nil,
		UserID:    uid,
		SessionID: nil,
		TokenHash: tokenHash,
		IPAddress: nil,
		UserAgent: &userAgent,
		ExpiresAt: now.Add(1 * time.Hour),
		CreatedAt: now,
		RevokedAt: nil,
		RevokedBy: nil,
	}

	require.NoError(t, repo.CreateRefreshToken(ctx, rt))
	require.NotEqual(t, uuid.Nil, rt.ID)

	// Get
	got, err := repo.GetRefreshToken(ctx, rt.ID)
	require.NoError(t, err)
	assert.Equal(t, rt.ID, got.ID)
	assert.Equal(t, rt.TokenHash, got.TokenHash)
	if got.UserAgent != nil {
		assert.Equal(t, *rt.UserAgent, *got.UserAgent)
	}
	assert.Nil(t, got.IPAddress)

	// Validate (not revoked, not expired)
	ok, err := repo.ValidateRefreshToken(ctx, rt.ID)
	require.NoError(t, err)
	assert.True(t, ok)

	// // Update: set revoked_at and revoked_by
	// revokedAt := time.Now().UTC()
	// revoker := uuid.Must(uuid.NewV7())
	// got.RevokedAt = &revokedAt
	// got.RevokedBy = &revoker
	// require.NoError(t, repo.UpdateRefreshToken(ctx, got))

	// // Validate now should be false (revoked)
	// ok, err = repo.ValidateRefreshToken(ctx, rt.ID)
	// require.NoError(t, err)
	// assert.False(t, ok)

	// // Update revoked_at -> clear revoked and set expires in past to test expiry behaviour
	// got.RevokedAt = nil
	// got.RevokedBy = nil
	// got.ExpiresAt = time.Now().Add(-1 * time.Minute)
	// require.NoError(t, repo.UpdateRefreshToken(ctx, got))

	// ok, err = repo.ValidateRefreshToken(ctx, rt.ID)
	// require.NoError(t, err)
	// assert.False(t, ok)

	// Delete and ensure Get returns ErrNotFound
	require.NoError(t, repo.DeleteRefreshToken(ctx, rt.ID))
	_, err = repo.GetRefreshToken(ctx, rt.ID)
	assert.ErrorIs(t, err, ErrNotFound)

	// Delete non-existent -> ErrNotFound
	err = repo.DeleteRefreshToken(ctx, rt.ID)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestRefreshTokenRepo_NotFoundCases(t *testing.T) {
	ctx := context.Background()
	repo, _, teardown := setupRefreshRepo(t)
	defer teardown()

	// random id should not exist
	rid := uuid.Must(uuid.NewV7())
	_, err := repo.GetRefreshToken(ctx, rid)
	assert.ErrorIs(t, err, ErrNotFound)

	// Update non-existent -> ErrNotFound
	rt := &models.RefreshToken{ID: rid}
	err = repo.UpdateRefreshToken(ctx, rt)
	assert.ErrorIs(t, err, ErrNotFound)

	// Validate non-existent -> ErrNotFound
	_, err = repo.ValidateRefreshToken(ctx, rid)
	assert.ErrorIs(t, err, ErrNotFound)
}
