package mocks

import (
	"context"
	"go-app/domain"

	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

type UserService struct {
    mock.Mock
}

func (_m *UserService) CreateUser(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, error) {
	ret := _m.Called(ctx, user)

	var r0 *domain.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.User)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *UserService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.User)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *UserService) UpdateUser(ctx context.Context, id uuid.UUID, u *domain.User) (*domain.User, error) {
	ret := _m.Called(ctx, id, u)

	var r0 *domain.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.User)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *UserService) GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error) {
	ret := _m.Called(ctx, filter)

	var r0 []domain.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]domain.User)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Error(0)
	}

	return r0
}
