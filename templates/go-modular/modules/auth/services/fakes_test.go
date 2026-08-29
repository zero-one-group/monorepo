package services

import (
	"context"
	"errors"
	"time"

	"{{ package_name | kebab_case }}/modules/auth/models"
	user_models "{{ package_name | kebab_case }}/modules/user/models"
	"{{ package_name | kebab_case }}/pkg/apputils"

	"github.com/gofrs/uuid/v5"
)

// fakeAuthRepo is an in-memory AuthRepositoryInterface. Passwords are stored
// hashed exactly as the service hands them over, and validated with the real hasher,
// so the tests exercise the same argon2 path production uses.
type fakeAuthRepo struct {
	err            error
	passwords      map[uuid.UUID]string
	sessions       map[uuid.UUID]*models.Session
	refreshTokens  map[uuid.UUID]*models.RefreshToken
	updatedPwds    map[uuid.UUID]string
	validSessions  map[uuid.UUID]bool
	validTokens    map[uuid.UUID]bool
	deletedSession []uuid.UUID
	deletedTokens  []uuid.UUID
	oneTimeTokens  map[uuid.UUID]*models.OneTimeToken
	deletedOTT     []uuid.UUID
	resent         []uuid.UUID
}

func newFakeAuthRepo() *fakeAuthRepo {
	return &fakeAuthRepo{
		passwords:     map[uuid.UUID]string{},
		sessions:      map[uuid.UUID]*models.Session{},
		refreshTokens: map[uuid.UUID]*models.RefreshToken{},
		updatedPwds:   map[uuid.UUID]string{},
		validSessions: map[uuid.UUID]bool{},
		validTokens:   map[uuid.UUID]bool{},
		oneTimeTokens: map[uuid.UUID]*models.OneTimeToken{},
	}
}

func (f *fakeAuthRepo) SetUserPassword(_ context.Context, p *models.UserPassword) error {
	if f.err != nil {
		return f.err
	}
	f.passwords[p.UserID] = p.PasswordHash
	return nil
}
func (f *fakeAuthRepo) UpdateUserPassword(_ context.Context, id uuid.UUID, hash string) error {
	if f.err != nil {
		return f.err
	}
	f.updatedPwds[id] = hash
	f.passwords[id] = hash
	return nil
}
func (f *fakeAuthRepo) ValidateUserPassword(_ context.Context, id uuid.UUID, password string) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	hash, ok := f.passwords[id]
	if !ok {
		return false, nil
	}
	return apputils.NewPasswordHasher().Validate(password, hash)
}
func (f *fakeAuthRepo) CreateSession(_ context.Context, s *models.Session) error {
	if f.err != nil {
		return f.err
	}
	if s.ID == uuid.Nil {
		s.ID = uuid.Must(uuid.NewV7())
	}
	f.sessions[s.ID] = s
	return nil
}
func (f *fakeAuthRepo) GetSession(_ context.Context, id uuid.UUID) (*models.Session, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.sessions[id], nil
}
func (f *fakeAuthRepo) UpdateSession(_ context.Context, s *models.Session) error {
	if f.err != nil {
		return f.err
	}
	f.sessions[s.ID] = s
	return nil
}
func (f *fakeAuthRepo) DeleteSession(_ context.Context, id uuid.UUID) error {
	if f.err != nil {
		return f.err
	}
	f.deletedSession = append(f.deletedSession, id)
	delete(f.sessions, id)
	return nil
}
func (f *fakeAuthRepo) ValidateSession(_ context.Context, id uuid.UUID) (bool, error) {
	return f.validSessions[id], f.err
}
func (f *fakeAuthRepo) CreateRefreshToken(_ context.Context, t *models.RefreshToken) error {
	if f.err != nil {
		return f.err
	}
	f.refreshTokens[t.ID] = t
	return nil
}
func (f *fakeAuthRepo) GetRefreshToken(_ context.Context, id uuid.UUID) (*models.RefreshToken, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.refreshTokens[id], nil
}
func (f *fakeAuthRepo) UpdateRefreshToken(_ context.Context, t *models.RefreshToken) error {
	if f.err != nil {
		return f.err
	}
	f.refreshTokens[t.ID] = t
	return nil
}
func (f *fakeAuthRepo) DeleteRefreshToken(_ context.Context, id uuid.UUID) error {
	if f.err != nil {
		return f.err
	}
	f.deletedTokens = append(f.deletedTokens, id)
	return nil
}
func (f *fakeAuthRepo) ValidateRefreshToken(_ context.Context, id uuid.UUID) (bool, error) {
	return f.validTokens[id], f.err
}
func (f *fakeAuthRepo) FindAllOneTimeTokens(context.Context) ([]*models.OneTimeToken, error) {
	if f.err != nil {
		return nil, f.err
	}
	out := []*models.OneTimeToken{}
	for _, t := range f.oneTimeTokens {
		out = append(out, t)
	}
	return out, nil
}
func (f *fakeAuthRepo) CreateOneTimeToken(_ context.Context, t *models.OneTimeToken) error {
	if f.err != nil {
		return f.err
	}
	if t.ID == uuid.Nil {
		t.ID = uuid.Must(uuid.NewV7())
	}
	f.oneTimeTokens[t.ID] = t
	return nil
}
func (f *fakeAuthRepo) GetOneTimeTokenByID(_ context.Context, id uuid.UUID) (*models.OneTimeToken, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.oneTimeTokens[id], nil
}
func (f *fakeAuthRepo) GetOneTimeTokenByTokenHash(_ context.Context, hash string) (*models.OneTimeToken, error) {
	if f.err != nil {
		return nil, f.err
	}
	for _, t := range f.oneTimeTokens {
		if t.TokenHash == hash {
			return t, nil
		}
	}
	return nil, errors.New("not found")
}
func (f *fakeAuthRepo) DeleteOneTimeToken(_ context.Context, id uuid.UUID) error {
	if f.err != nil {
		return f.err
	}
	f.deletedOTT = append(f.deletedOTT, id)
	delete(f.oneTimeTokens, id)
	return nil
}
func (f *fakeAuthRepo) UpdateOneTimeTokenLastSentAt(_ context.Context, id uuid.UUID, at time.Time) error {
	if f.err != nil {
		return f.err
	}
	if t, ok := f.oneTimeTokens[id]; ok {
		t.LastSentAt = &at
	}
	f.resent = append(f.resent, id)
	return nil
}

// fakeUserService resolves users by email/username from a fixed set.
type fakeUserService struct {
	users    []*user_models.User
	err      error
	verified []uuid.UUID
}

func (f *fakeUserService) CreateUser(context.Context, *user_models.User) error { return f.err }
func (f *fakeUserService) GetUserByID(_ context.Context, id uuid.UUID) (*user_models.User, error) {
	for _, u := range f.users {
		if u.ID == id {
			return u, f.err
		}
	}
	return nil, f.err
}
func (f *fakeUserService) ListUsers(context.Context, *user_models.FilterUser) ([]*user_models.User, error) {
	return f.users, f.err
}
func (f *fakeUserService) UpdateUser(context.Context, *user_models.User) error { return f.err }
func (f *fakeUserService) DeleteUser(context.Context, uuid.UUID) error         { return f.err }
func (f *fakeUserService) GetUserByEmail(_ context.Context, email string) (*user_models.User, error) {
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
func (f *fakeUserService) GetUserByUsername(_ context.Context, username string) (*user_models.User, error) {
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
func (f *fakeUserService) MarkEmailVerified(_ context.Context, id uuid.UUID) error {
	f.verified = append(f.verified, id)
	return f.err
}

func newTestAuthService(repo *fakeAuthRepo, users *fakeUserService) *AuthService {
	return NewAuthService(AuthServiceOpts{
		AuthRepo:     repo,
		UserService:  users,
		JWTSecretKey: []byte("test-secret-key-for-unit-tests-only"),
		BaseURL:      "https://api.example.com",
	})
}
