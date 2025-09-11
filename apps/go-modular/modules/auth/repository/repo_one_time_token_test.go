package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"testing"
	"time"

	"go-modular/modules/auth/models"
	"go-modular/pkg/apputils"
	"go-modular/pkg/testutils"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
}

// setupAuthRepo starts the test DB, runs migrations and returns repo and a seeded user id.
// It relies on the application's seeders (RunAppMigrations) to create test users.
func setupAuthRepo(t *testing.T) (*AuthRepository, uuid.UUID, func()) {
	t.Helper()
	te := testutils.NewTestEnv(t)
	pool, _, err := te.SetupPostgres()
	require.NoError(t, err)
	require.NotNil(t, pool)

	// ensure env/config is set for migrations and app code
	te.SetupConfig()
	// run migrations and seeders to ensure tables and seed data exist
	te.RunAppMigrations()

	logger := newLogger()
	repo := NewAuthRepository(pool, logger)

	// try to find a seeded user created by user_factory
	var uid uuid.UUID
	err = pool.QueryRow(context.Background(), `SELECT id FROM public.users WHERE username = $1 LIMIT 1`, "admin").Scan(&uid)
	if err != nil {
		// try alternate common seeded username
		err = pool.QueryRow(context.Background(), `SELECT id FROM public.users WHERE username = $1 LIMIT 1`, "johndoe").Scan(&uid)
		if err != nil {
			// fallback: take any existing user
			err = pool.QueryRow(context.Background(), `SELECT id FROM public.users LIMIT 1`).Scan(&uid)
			require.NoError(t, err, "failed to find any seeded user required by tests")
		}
	}

	teardown := func() {
		// cleanup tokens between tests (testutils manages container lifecycle)
		_, _ = pool.Exec(context.Background(), `TRUNCATE TABLE `+models.OneTimeTokenTable+` CASCADE`)
		pool.Close()
	}

	return repo, uid, teardown
}

func TestOneTimeToken_CRUD_and_NotFound(t *testing.T) {
	ctx := context.Background()
	repo, uid, teardown := setupAuthRepo(t)
	defer teardown()

	// create token
	now := time.Now().UTC().Truncate(time.Second)

	// generate raw token and compute tokenHash like svc_verification.go
	rawToken, err := apputils.GenerateURLSafeToken(48)
	require.NoError(t, err)
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	token := &models.OneTimeToken{
		ID:        uuid.Nil,
		UserID:    &uid,
		Subject:   "verify_email",
		TokenHash: tokenHash,
		RelatesTo: "email",
		Metadata:  map[string]any{"reason": "signup"},
		CreatedAt: now,
		ExpiresAt: now.Add(1 * time.Hour),
		// LastSentAt zero value
	}

	// Create
	require.NoError(t, repo.CreateOneTimeToken(ctx, token))
	require.NotEqual(t, uuid.Nil, token.ID)

	// Get by token hash (use same computed hash)
	got2, err := repo.GetOneTimeTokenByTokenHash(ctx, tokenHash)
	require.NoError(t, err)
	assert.Equal(t, token.ID, got2.ID)

	// (optional) validate retrieval by hash of rawToken matches
	g, err := repo.GetOneTimeTokenByTokenHash(ctx, tokenHash)
	require.NoError(t, err)
	assert.Equal(t, token.ID, g.ID)

	// Find all includes our token
	all, err := repo.FindAllOneTimeTokens(ctx)
	require.NoError(t, err)
	found := false
	for _, tt := range all {
		if tt.ID == token.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "created token should be present in FindAllOneTimeTokens")

	// Update last_sent_at
	newSent := time.Now().UTC().Truncate(time.Second)
	require.NoError(t, repo.UpdateOneTimeTokenLastSentAt(ctx, token.ID, newSent))

	updated, err := repo.GetOneTimeTokenByID(ctx, token.ID)
	require.NoError(t, err)
	require.NotNil(t, updated.LastSentAt, "LastSentAt should be set after update")
	// compare dereferenced pointer within a second
	assert.WithinDuration(t, newSent, *updated.LastSentAt, time.Second)

	// Delete
	require.NoError(t, repo.DeleteOneTimeToken(ctx, token.ID))

	// subsequent Get should return ErrNotFound
	_, err = repo.GetOneTimeTokenByID(ctx, token.ID)
	assert.ErrorIs(t, err, ErrNotFound)

	// Deleting again returns ErrNotFound
	err = repo.DeleteOneTimeToken(ctx, token.ID)
	assert.ErrorIs(t, err, ErrNotFound)

	// Update last_sent_at on missing ID returns ErrNotFound
	err = repo.UpdateOneTimeTokenLastSentAt(ctx, token.ID, time.Now())
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestOneTimeToken_Get_NotFound(t *testing.T) {
	ctx := context.Background()
	repo, _, teardown := setupAuthRepo(t)
	defer teardown()

	// random id should not exist
	rid := uuid.Must(uuid.NewV7())
	_, err := repo.GetOneTimeTokenByID(ctx, rid)
	assert.ErrorIs(t, err, ErrNotFound)

	_, err = repo.GetOneTimeTokenByTokenHash(ctx, "non-existent-hash-xyz")
	assert.ErrorIs(t, err, ErrNotFound)
}
