package service

import (
	"context"
	"fmt"
	"{{ package_name }}/domain"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, error)
	GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(u UserRepository) *UserService {
	return &UserService{
		userRepo: u,
	}
}

// CreateUser adds a new user.
func (us *UserService) CreateUser(
	ctx context.Context,
	u *domain.CreateUserRequest,
) (*domain.User, error) {
	span, spanCtx := opentracing.StartSpanFromContext(
		ctx,
		"UserService.GetUser",
	)
	defer span.Finish()

	createdUser, err := us.userRepo.CreateUser(spanCtx, u)
	if err != nil {
		span.SetTag("error", true)
		return nil, err
	}
	return createdUser, nil
}

// GetUser fetches a user by ID.
func (us *UserService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {
	span, spanCtx := opentracing.StartSpanFromContext(
		ctx,
		"UserService.GetUser",
	)
	defer span.Finish()

	user, err := us.userRepo.GetUser(spanCtx, id)
	if err != nil {
		span.SetTag("error", true)
		return nil, err
	}

	return user, nil
}

// UpdateUser updates name/email of an existing user.
func (us *UserService) UpdateUser(
	ctx context.Context,
	id uuid.UUID,
	u *domain.User,
) (*domain.User, error) {
	span, spanCtx := opentracing.StartSpanFromContext(
		ctx,
		"UserService.GetUser",
	)
	defer span.Finish()

	existing, err := us.userRepo.GetUser(ctx, id)
	if err != nil {
		span.SetTag("error", true)
		return nil, err
	}
	if existing == nil {
		span.SetTag("error", true)
		return nil, domain.ErrUserNotFound
	}

	existing.Name = u.Name
	existing.Email = u.Email

	_, err = us.userRepo.UpdateUser(spanCtx, id, existing)
	if err != nil {
		span.SetTag("error", true)
		return nil, err
	}

	return existing, nil
}

// DeleteUser removes a user by ID.
func (us *UserService) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	span, spanCtx := opentracing.StartSpanFromContext(
		ctx,
		"UserService.GetUser",
	)
	defer span.Finish()

	user, err := us.userRepo.GetUser(spanCtx, id)
	if err != nil {
		span.SetTag("error", true)
		return err
	}
	if user == nil {
		span.SetTag("error", true)
		return domain.ErrUserNotFound
	}

	err = us.userRepo.DeleteUser(spanCtx, id)
	if err != nil {
		span.SetTag("error", true)
		return err
	}

	return nil
}

func (us *UserService) GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error) {
	span, spanCtx := opentracing.StartSpanFromContext(
		ctx,
		"UserService.GetUser",
	)
	defer span.Finish()

	users, err := us.userRepo.GetUserList(spanCtx, filter)
	if err != nil {
		fmt.Println(err)
		span.SetTag("error", true)
		return nil, err
	}

	return users, nil
}
