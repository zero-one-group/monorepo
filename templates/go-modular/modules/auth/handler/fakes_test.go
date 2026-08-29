package handler

import (
	"context"

	"{{ package_name | kebab_case }}/modules/auth/models"

	"github.com/gofrs/uuid/v5"
)

// crudFake records the last model each CRUD method received and answers from
// canned values, so handler tests can assert the HTTP→model mapping precisely.
type crudFake struct {
	fakeAuthService
	crudErr      error
	session      *models.Session
	token        *models.RefreshToken
	lastSession  *models.Session
	lastToken    *models.RefreshToken
	lastPassword *models.UserPassword
	lastUpdatePw []string // userID, current, new
	deleted      []uuid.UUID
	verifyOK     bool
	verifyErr    error
	lastVerify   []string // email/token, redirectTo
}

func (f *crudFake) SetUserPassword(_ context.Context, p *models.UserPassword) error {
	f.lastPassword = p
	return f.crudErr
}
func (f *crudFake) UpdateUserPassword(_ context.Context, id uuid.UUID, cur, next string) error {
	f.lastUpdatePw = []string{id.String(), cur, next}
	return f.crudErr
}
func (f *crudFake) CreateSession(_ context.Context, s *models.Session) error {
	f.lastSession = s
	return f.crudErr
}
func (f *crudFake) GetSession(context.Context, uuid.UUID) (*models.Session, error) {
	return f.session, f.crudErr
}
func (f *crudFake) UpdateSession(_ context.Context, s *models.Session) error {
	f.lastSession = s
	return f.crudErr
}
func (f *crudFake) DeleteSession(_ context.Context, id uuid.UUID) error {
	f.deleted = append(f.deleted, id)
	return f.crudErr
}
func (f *crudFake) CreateRefreshToken(_ context.Context, t *models.RefreshToken) error {
	f.lastToken = t
	return f.crudErr
}
func (f *crudFake) GetRefreshToken(context.Context, uuid.UUID) (*models.RefreshToken, error) {
	return f.token, f.crudErr
}
func (f *crudFake) UpdateRefreshToken(_ context.Context, t *models.RefreshToken) error {
	f.lastToken = t
	return f.crudErr
}
func (f *crudFake) DeleteRefreshToken(_ context.Context, id uuid.UUID) error {
	f.deleted = append(f.deleted, id)
	return f.crudErr
}
func (f *crudFake) InitiateEmailVerification(_ context.Context, email, redirectTo string) error {
	f.lastVerify = []string{email, redirectTo}
	return f.verifyErr
}
func (f *crudFake) ValidateEmailVerification(_ context.Context, token string) (bool, error) {
	f.lastVerify = []string{token}
	return f.verifyOK, f.verifyErr
}
func (f *crudFake) RevokeEmailVerification(_ context.Context, token string) error {
	f.lastVerify = []string{token}
	return f.verifyErr
}
func (f *crudFake) ResendEmailVerification(_ context.Context, email, redirectTo string) error {
	f.lastVerify = []string{email, redirectTo}
	return f.verifyErr
}
