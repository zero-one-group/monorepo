package repository

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"go-modular/modules/auth/models"
// 	"go-modular/pkg/apputils"
// 	"go-modular/pkg/testutils"

// 	"github.com/gofrs/uuid/v5"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestUserPasswordRepo(t *testing.T) {
// 	t.Run("SetValidateUpdatePassword", func(t *testing.T) {
// 		ctx := context.Background()
// 		te := testutils.NewTestEnv(t)
// 		pool, _, err := te.SetupPostgres()
// 		require.NoError(t, err)
// 		require.NotNil(t, pool)
// 		te.SetupConfig()
// 		te.RunAppMigrations()

// 		// use a seeded user created by seeders (avoid inserting/deleting seed users)
// 		var uid uuid.UUID
// 		err = pool.QueryRow(ctx, `SELECT id FROM public.users WHERE username = $1 LIMIT 1`, "admin").Scan(&uid)
// 		if err != nil {
// 			err = pool.QueryRow(ctx, `SELECT id FROM public.users LIMIT 1`).Scan(&uid)
// 			require.NoError(t, err, "failed to find seeded user required by tests")
// 		}

// 		repo := NewAuthRepository(pool, newLogger())
// 		defer func() {
// 			_, _ = pool.Exec(ctx, `DELETE FROM public.user_passwords WHERE user_id = $1`, uid)
// 			pool.Close()
// 		}()

// 		hasher := apputils.NewPasswordHasher()
// 		hash1, err := hasher.Hash("secret123")
// 		require.NoError(t, err)

// 		up := &models.UserPassword{
// 			UserID:       uid,
// 			PasswordHash: hash1,
// 			CreatedAt:    time.Now().UTC(),
// 		}

// 		// set password
// 		require.NoError(t, repo.SetUserPassword(ctx, up))

// 		// validate correct password
// 		ok, err := repo.ValidateUserPassword(ctx, uid, "secret123")
// 		require.NoError(t, err)
// 		assert.True(t, ok)

// 		// update password
// 		hash2, err := hasher.Hash("newpass456")
// 		require.NoError(t, err)
// 		require.NoError(t, repo.UpdateUserPassword(ctx, uid, hash2))

// 		// validate new password
// 		ok, err = repo.ValidateUserPassword(ctx, uid, "newpass456")
// 		require.NoError(t, err)
// 		assert.True(t, ok)
// 	})
// }
