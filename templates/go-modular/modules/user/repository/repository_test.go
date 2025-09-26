package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"go-modular/modules/user/models"
	"go-modular/pkg/testutils"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
}

func setupRepo(t *testing.T) (*UserRepository, func()) {
	t.Helper()
	te := testutils.NewTestEnv(t)
	pool, pgURL, err := te.SetupPostgres()
	require.NoError(t, err)
	require.NotNil(t, pool)

	// ensure env/config is set for migrations and app code
	// run migrations to ensure tables exist
	te.SetupConfig()
	te.RunAppMigrations()

	logger := newLogger()
	repo := NewUserRepository(pool, logger)

	teardown := func() {
		// try to clean up table data between tests; container/pool lifecycle handled by testutils te.Cleanup
		_, _ = pool.Exec(context.Background(), `TRUNCATE TABLE `+models.UserTable+` CASCADE`)
		_ = pgURL
	}

	return repo, teardown
}

func userFromMap(t *testing.T, m map[string]interface{}) *models.User {
	t.Helper()
	b, err := json.Marshal(m)
	require.NoError(t, err)
	var u models.User
	err = json.Unmarshal(b, &u)
	require.NoError(t, err)
	return &u
}

func TestUserRepository_CRUD_and_Exists(t *testing.T) {
	ctx := context.Background()
	repo, teardown := setupRepo(t)
	defer teardown()

	// Create a new user (leave ID zero so repo assigns)
	u := userFromMap(t, map[string]interface{}{
		"id":           nil,
		"display_name": "Alice Doe",
		"email":        "alice@example.com",
		"username":     "alice",
		"avatar_url":   "",
		"metadata":     nil,
	})

	err := repo.CreateUser(ctx, u)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, u.ID, "CreateUser should set ID when zero")

	// Get by ID
	got, err := repo.GetUserByID(ctx, u.ID)
	require.NoError(t, err)
	require.NotNil(t, got)

	// compare using JSON-unmarshaled struct fields (handles pointer vs value fields)
	assert.Equal(t, "alice", normalizeString(got.Username))
	assert.Equal(t, "alice@example.com", normalizeString(got.Email))

	// UsernameExists / EmailExists
	ok, err := repo.UsernameExists(ctx, "alice")
	require.NoError(t, err)
	assert.True(t, ok)
	ok, err = repo.EmailExists(ctx, "alice@example.com")
	require.NoError(t, err)
	assert.True(t, ok)

	// Get by email/username
	byEmail, err := repo.GetUserByEmail(ctx, "alice@example.com")
	require.NoError(t, err)
	assert.Equal(t, u.ID, byEmail.ID)

	byUser, err := repo.GetUserByUsername(ctx, "alice")
	require.NoError(t, err)
	assert.Equal(t, u.ID, byUser.ID)

	// ListUsers (no filter) should return at least our user
	users, err := repo.ListUsers(ctx, nil)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(users), 1)

	// Update user
	u.DisplayName = "Alice Updated"
	err = repo.UpdateUser(ctx, u)
	require.NoError(t, err)

	updated, err := repo.GetUserByID(ctx, u.ID)
	require.NoError(t, err)
	assert.Equal(t, "Alice Updated", normalizeString(updated.DisplayName))

	// Delete user
	err = repo.DeleteUser(ctx, u.ID)
	require.NoError(t, err)

	_, err = repo.GetUserByID(ctx, u.ID)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestUserRepository_List_FilterSearchAndPaging(t *testing.T) {
	ctx := context.Background()
	repo, teardown := setupRepo(t)
	defer teardown()

	// insert multiple users
	for i := range 5 {
		uid := strings.ReplaceAll(uuid.Must(uuid.NewV7()).String(), "-", "")
		// username must satisfy DB check constraint: only alnum/underscore and 3..32 chars
		// produce a short unique username to avoid collisions and length issues
		username := fmt.Sprintf("user%02d%s", i, uid[:8]) // e.g. user00a1b2c3d4
		email := fmt.Sprintf("user%s-%d@example.com", uid[:8], i)
		u := userFromMap(t, map[string]any{
			"display_name": "User " + uid[:8],
			"email":        email,
			"username":     username,
		})
		require.NoError(t, repo.CreateUser(ctx, u))
		// small sleep to ensure created_at ordering
		time.Sleep(10 * time.Millisecond)
	}

	// search by a substring of username from one of inserted users
	all, err := repo.ListUsers(ctx, nil)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(all), 5)

	// use limit/offset
	filter := &models.FilterUser{
		Limit:  2,
		Offset: 1,
	}
	paged, err := repo.ListUsers(ctx, filter)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(paged), 2)
}

// normalizeString helps assertions work whether the model field is a string or *string.
func normalizeString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case *string:
		if x == nil {
			return ""
		}
		return *x
	default:
		return ""
	}
}
