package service

import (
	"context"
	"go-standard/domain"
	"go-standard/internal/logging"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
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
	createdUser, err := us.userRepo.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

// GetUser fetches a user by ID.
func (us *UserService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {
	tracer := otel.Tracer("service.user")
	ctxTrace, span := tracer.Start(ctx, "UserService.GetUser")
	defer span.End()

	user, err := us.userRepo.GetUser(ctxTrace, id)
	if err != nil {
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

	existing, err := us.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrUserNotFound
	}

	existing.Name = u.Name
	existing.Email = u.Email

	_, err = us.userRepo.UpdateUser(ctx, id, existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteUser removes a user by ID.
func (us *UserService) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {

	user, err := us.userRepo.GetUser(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	err = us.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error) {
	users, err := us.userRepo.GetUserList(ctx, filter)
	if err != nil {
		logging.LogError(ctx, err, "get_user_list_service")
		return nil, err
	}

	return users, nil
}
