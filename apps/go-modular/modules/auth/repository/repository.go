package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"go-modular/modules/auth/models"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthRepositoryInterface defines the contract for user data access.
type AuthRepositoryInterface interface {
	// User password operations
	SetUserPassword(ctx context.Context, userPassword *models.UserPassword) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	ValidateUserPassword(ctx context.Context, userID uuid.UUID, password string) (bool, error)

	// Session operations
	CreateSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	UpdateSession(ctx context.Context, session *models.Session) error
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	ValidateSession(ctx context.Context, sessionID uuid.UUID) (bool, error)

	// Refresh token operations
	CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, tokenID uuid.UUID) (*models.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, tokenID uuid.UUID) error
	ValidateRefreshToken(ctx context.Context, tokenID uuid.UUID) (bool, error)

	// OneTimeToken operations
	FindAllOneTimeTokens(ctx context.Context) ([]*models.OneTimeToken, error)
	CreateOneTimeToken(ctx context.Context, token *models.OneTimeToken) error
	GetOneTimeTokenByID(ctx context.Context, tokenID uuid.UUID) (*models.OneTimeToken, error)
	GetOneTimeTokenByTokenHash(ctx context.Context, tokenHash string) (*models.OneTimeToken, error)
	DeleteOneTimeToken(ctx context.Context, tokenID uuid.UUID) error
	UpdateOneTimeTokenLastSentAt(ctx context.Context, tokenID uuid.UUID, lastSentAt time.Time) error
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

// Sentinel error for not found
var ErrNotFound = errors.New("not found")
