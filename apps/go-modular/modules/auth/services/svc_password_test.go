package services

import (
	"context"
	"errors"
	"testing"

	"go-modular/modules/auth/models"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthService_RequiredOptionsAndDefaults(t *testing.T) {
	repo, users := newFakeAuthRepo(), &fakeUserService{}
	assert.PanicsWithValue(t, "AuthRepo is required", func() {
		NewAuthService(AuthServiceOpts{UserService: users, JWTSecretKey: []byte("k"), BaseURL: "https://x"})
	})
	assert.PanicsWithValue(t, "UserService is required", func() {
		NewAuthService(AuthServiceOpts{AuthRepo: repo, JWTSecretKey: []byte("k"), BaseURL: "https://x"})
	})
	assert.PanicsWithValue(t, "JWTSecretKey is required", func() {
		NewAuthService(AuthServiceOpts{AuthRepo: repo, UserService: users, BaseURL: "https://x"})
	})
	assert.PanicsWithValue(t, "BaseURL is required", func() {
		NewAuthService(AuthServiceOpts{AuthRepo: repo, UserService: users, JWTSecretKey: []byte("k")})
	})

	svc := newTestAuthService(repo, users)
	assert.Equal(t, "HS256", string(svc.signingAlg))
	assert.Equal(t, 24.0, svc.accessTokenExpiry.Hours())
	assert.Equal(t, 7*24.0, svc.refreshTokenExpiry.Hours())
}

func TestSetUserPassword_HashesBeforeStoring(t *testing.T) {
	repo := newFakeAuthRepo()
	svc := newTestAuthService(repo, &fakeUserService{})
	id := uuid.Must(uuid.NewV7())

	require.NoError(t, svc.SetUserPassword(context.Background(), &models.UserPassword{UserID: id, PasswordHash: "plain-text"}))

	stored := repo.passwords[id]
	assert.NotEqual(t, "plain-text", stored, "plaintext must never reach the repository")
	ok, err := svc.ValidateUserPassword(context.Background(), id, "plain-text")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestSetUserPassword_Validation(t *testing.T) {
	svc := newTestAuthService(newFakeAuthRepo(), &fakeUserService{})
	assert.EqualError(t, svc.SetUserPassword(context.Background(), nil), "password is required")
	assert.EqualError(t, svc.SetUserPassword(context.Background(), &models.UserPassword{}), "password is required")
}

func TestSetUserPassword_RepositoryError(t *testing.T) {
	repo := newFakeAuthRepo()
	repo.err = errors.New("db down")
	svc := newTestAuthService(repo, &fakeUserService{})
	assert.EqualError(t, svc.SetUserPassword(context.Background(), &models.UserPassword{UserID: uuid.Must(uuid.NewV7()), PasswordHash: "x"}), "db down")
}

func TestUpdateUserPassword(t *testing.T) {
	ctx := context.Background()
	setup := func() (*fakeAuthRepo, *AuthService, uuid.UUID) {
		repo := newFakeAuthRepo()
		svc := newTestAuthService(repo, &fakeUserService{})
		id := uuid.Must(uuid.NewV7())
		require.NoError(t, svc.SetUserPassword(ctx, &models.UserPassword{UserID: id, PasswordHash: "old-password"}))
		return repo, svc, id
	}

	t.Run("rejects empty inputs before touching the repository", func(t *testing.T) {
		_, svc, id := setup()
		assert.EqualError(t, svc.UpdateUserPassword(ctx, id, "old-password", ""), "new password is required")
		assert.EqualError(t, svc.UpdateUserPassword(ctx, id, "", "new-password"), "current password is required")
	})
	t.Run("rejects a wrong current password", func(t *testing.T) {
		repo, svc, id := setup()
		assert.EqualError(t, svc.UpdateUserPassword(ctx, id, "wrong", "new-password"), "current password is incorrect")
		assert.Empty(t, repo.updatedPwds)
	})
	t.Run("hashes and stores the new password", func(t *testing.T) {
		repo, svc, id := setup()
		require.NoError(t, svc.UpdateUserPassword(ctx, id, "old-password", "new-password"))
		assert.NotEqual(t, "new-password", repo.updatedPwds[id])
		ok, _ := svc.ValidateUserPassword(ctx, id, "new-password")
		assert.True(t, ok)
		ok, _ = svc.ValidateUserPassword(ctx, id, "old-password")
		assert.False(t, ok)
	})
	t.Run("repository error while validating", func(t *testing.T) {
		repo, svc, id := setup()
		repo.err = errors.New("db down")
		assert.EqualError(t, svc.UpdateUserPassword(ctx, id, "old-password", "new-password"), "db down")
	})
}

func TestValidateUserPassword_EmptyPassword(t *testing.T) {
	svc := newTestAuthService(newFakeAuthRepo(), &fakeUserService{})
	ok, err := svc.ValidateUserPassword(context.Background(), uuid.Must(uuid.NewV7()), "")
	assert.False(t, ok)
	assert.EqualError(t, err, "password is required")
}
