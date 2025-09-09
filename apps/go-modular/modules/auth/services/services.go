package services

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"go-modular/modules/auth/models"
	"go-modular/modules/auth/repository"
)

// AuthServiceInterface defines the contract for user business logic.
type AuthServiceInterface interface {
	SetUserPassword(ctx context.Context, userPassword *models.UserPassword) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error
	ValidateUserPassword(ctx context.Context, userID uuid.UUID, password string) (bool, error)
}

// Ensure AuthService implements AuthServiceInterface
var _ AuthServiceInterface = (*AuthService)(nil)

// AuthService implements user business logic using a UserRepositoryInterface.
type AuthService struct {
	authRepo repository.AuthRepositoryInterface
}

type AuthServiceOpts struct {
	AuthRepo repository.AuthRepositoryInterface
}

// NewAuthService creates a new AuthService.
func NewAuthService(opts AuthServiceOpts) *AuthService {
	return &AuthService{
		authRepo: opts.AuthRepo,
	}
}
