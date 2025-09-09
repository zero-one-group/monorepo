package repository

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-modular/modules/auth/models"
)

// AuthRepositoryInterface defines the contract for user data access.
type AuthRepositoryInterface interface {
	SetUserPassword(ctx context.Context, userPassword *models.UserPassword) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	ValidateUserPassword(ctx context.Context, userID uuid.UUID, password string) (bool, error)
}

// Ensure AuthRepository implements AuthRepositoryInterface
var _ AuthRepositoryInterface = (*AuthRepository)(nil)

// AuthRepository is an implementation of AuthRepositoryInterface using pgxpool.
type AuthRepository struct {
	pgPool *pgxpool.Pool
	logger *slog.Logger
}

// NewAuthRepository creates a new AuthRepository with pgxpool and slog logger.
func NewAuthRepository(pgPool *pgxpool.Pool, logger *slog.Logger) *AuthRepository {
	return &AuthRepository{
		pgPool: pgPool,
		logger: logger,
	}
}
