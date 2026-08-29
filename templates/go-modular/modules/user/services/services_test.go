package services

import (
	"context"
	"errors"
	"testing"

	"{{ package_name | kebab_case }}/modules/user/models"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeUserRepo is an in-memory UserRepositoryInterface that records calls and
// lets each test decide which usernames already exist and which calls fail.
type fakeUserRepo struct {
	users        map[uuid.UUID]*models.User
	taken        map[string]bool
	err          error
	created      []*models.User
	updated      []*models.User
	deleted      []uuid.UUID
	existsCalled []string
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: map[uuid.UUID]*models.User{}, taken: map[string]bool{}}
}

func (f *fakeUserRepo) CreateUser(_ context.Context, u *models.User) error {
	if f.err != nil {
		return f.err
	}
	f.created = append(f.created, u)
	f.users[u.ID] = u
	return nil
}
func (f *fakeUserRepo) GetUserByID(_ context.Context, id uuid.UUID) (*models.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.users[id], nil
}
func (f *fakeUserRepo) ListUsers(_ context.Context, _ *models.FilterUser) ([]*models.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	out := []*models.User{}
	for _, u := range f.users {
		out = append(out, u)
	}
	return out, nil
}
func (f *fakeUserRepo) UpdateUser(_ context.Context, u *models.User) error {
	if f.err != nil {
		return f.err
	}
	f.updated = append(f.updated, u)
	return nil
}
func (f *fakeUserRepo) DeleteUser(_ context.Context, id uuid.UUID) error {
	if f.err != nil {
		return f.err
	}
	f.deleted = append(f.deleted, id)
	return nil
}
func (f *fakeUserRepo) UsernameExists(_ context.Context, username string) (bool, error) {
	f.existsCalled = append(f.existsCalled, username)
	if f.err != nil {
		return false, f.err
	}
	return f.taken[username], nil
}
func (f *fakeUserRepo) EmailExists(_ context.Context, _ string) (bool, error) { return false, f.err }
func (f *fakeUserRepo) GetUserByEmail(_ context.Context, email string) (*models.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	for _, u := range f.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}
func (f *fakeUserRepo) GetUserByUsername(_ context.Context, username string) (*models.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	for _, u := range f.users {
		if u.Username != nil && *u.Username == username {
			return u, nil
		}
	}
	return nil, nil
}

func newSvc(repo *fakeUserRepo) *UserService {
	return NewUserService(UserServiceOpts{UserRepo: repo})
}

func TestCreateUser_DerivesUsernameFromEmailAndDefaults(t *testing.T) {
	repo := newFakeUserRepo()
	u := &models.User{Email: "Jane.Doe+shop@example.com", DisplayName: "Jane"}

	require.NoError(t, newSvc(repo).CreateUser(context.Background(), u))

	assert.NotEqual(t, uuid.Nil, u.ID, "an id is minted when none is given")
	require.NotNil(t, u.Username)
	assert.Equal(t, "janedoeshop", *u.Username, "lowercased, non [a-z0-9_] stripped, domain dropped")
	require.NotNil(t, u.Metadata)
	assert.Equal(t, "UTC", u.Metadata.Timezone)
	require.Len(t, repo.created, 1)
}

func TestCreateUser_SuffixesUsernameUntilUnique(t *testing.T) {
	repo := newFakeUserRepo()
	repo.taken["jane"] = true
	repo.taken["jane_1"] = true
	u := &models.User{Email: "jane@example.com"}

	require.NoError(t, newSvc(repo).CreateUser(context.Background(), u))

	assert.Equal(t, "jane_2", *u.Username)
	assert.Equal(t, []string{"jane", "jane_1", "jane_2"}, repo.existsCalled)
}

func TestCreateUser_FallsBackToUserWhenLocalPartHasNoValidChars(t *testing.T) {
	repo := newFakeUserRepo()
	u := &models.User{Email: "!!!@example.com"}
	require.NoError(t, newSvc(repo).CreateUser(context.Background(), u))
	assert.Equal(t, "user", *u.Username)
}

func TestCreateUser_KeepsExplicitUsernameIDAndMetadata(t *testing.T) {
	repo := newFakeUserRepo()
	id := uuid.Must(uuid.NewV7())
	name := "custom_name"
	u := &models.User{ID: id, Email: "x@example.com", Username: &name, Metadata: &models.UserMetadata{Timezone: "Asia/Jakarta"}}

	require.NoError(t, newSvc(repo).CreateUser(context.Background(), u))

	assert.Equal(t, id, u.ID)
	assert.Equal(t, "custom_name", *u.Username)
	assert.Equal(t, "Asia/Jakarta", u.Metadata.Timezone)
	assert.Empty(t, repo.existsCalled, "no uniqueness lookup when the username is explicit")
}

func TestCreateUser_PropagatesRepositoryErrors(t *testing.T) {
	repo := newFakeUserRepo()
	repo.err = errors.New("db down")
	err := newSvc(repo).CreateUser(context.Background(), &models.User{Email: "a@b.c"})
	assert.EqualError(t, err, "db down")
	assert.Empty(t, repo.created)
}

func TestPassThroughs(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newSvc(repo)
	ctx := context.Background()
	id := uuid.Must(uuid.NewV7())
	uname := "amy"
	repo.users[id] = &models.User{ID: id, Email: "amy@example.com", Username: &uname}

	got, err := svc.GetUserByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "amy@example.com", got.Email)

	byEmail, err := svc.GetUserByEmail(ctx, "amy@example.com")
	require.NoError(t, err)
	assert.Equal(t, id, byEmail.ID)

	byUsername, err := svc.GetUserByUsername(ctx, "amy")
	require.NoError(t, err)
	assert.Equal(t, id, byUsername.ID)

	list, err := svc.ListUsers(ctx, &models.FilterUser{})
	require.NoError(t, err)
	assert.Len(t, list, 1)

	require.NoError(t, svc.UpdateUser(ctx, got))
	assert.Len(t, repo.updated, 1)

	require.NoError(t, svc.DeleteUser(ctx, id))
	assert.Equal(t, []uuid.UUID{id}, repo.deleted)
}

func TestMarkEmailVerified(t *testing.T) {
	t.Run("stamps email_verified_at and saves", func(t *testing.T) {
		repo := newFakeUserRepo()
		id := uuid.Must(uuid.NewV7())
		repo.users[id] = &models.User{ID: id, Email: "v@example.com"}

		require.NoError(t, newSvc(repo).MarkEmailVerified(context.Background(), id))

		require.Len(t, repo.updated, 1)
		assert.NotNil(t, repo.updated[0].EmailVerifiedAt)
	})
	t.Run("unknown user", func(t *testing.T) {
		repo := newFakeUserRepo()
		err := newSvc(repo).MarkEmailVerified(context.Background(), uuid.Must(uuid.NewV7()))
		assert.EqualError(t, err, "user not found")
		assert.Empty(t, repo.updated)
	})
	t.Run("repository error", func(t *testing.T) {
		repo := newFakeUserRepo()
		repo.err = errors.New("boom")
		assert.EqualError(t, newSvc(repo).MarkEmailVerified(context.Background(), uuid.Must(uuid.NewV7())), "boom")
	})
}
