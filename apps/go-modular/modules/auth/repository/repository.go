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
	// User password management
	SetUserPassword(ctx context.Context, userPassword *models.UserPassword) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	ValidateUserPassword(ctx context.Context, userID uuid.UUID, password string) (bool, error)

	// Session management
	CreateSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	UpdateSession(ctx context.Context, session *models.Session) error
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	ValidateSession(ctx context.Context, sessionID uuid.UUID) (bool, error)

	// Refresh token management
	CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, tokenID uuid.UUID) (*models.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, tokenID uuid.UUID) error
	ValidateRefreshToken(ctx context.Context, tokenID uuid.UUID) (bool, error)
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
