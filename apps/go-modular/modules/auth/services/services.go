package services

import (
	"context"
	"time"

	"go-modular/modules/auth/models"
	"go-modular/modules/auth/repository"

	"github.com/gofrs/uuid/v5"
	"github.com/lestrrat-go/jwx/jwa"

	user_service "go-modular/modules/user/services"
)

// AuthServiceInterface defines the contract for user business logic.
type AuthServiceInterface interface {
	// User password management
	SetUserPassword(ctx context.Context, userPassword *models.UserPassword) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error
	ValidateUserPassword(ctx context.Context, userID uuid.UUID, password string) (bool, error)

	// Session management
	CreateSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	UpdateSession(ctx context.Context, session *models.Session) error
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	ValidateSession(ctx context.Context, sessionID uuid.UUID) (bool, error)

	// Refresh token management
	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, tokenID uuid.UUID) (*models.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, tokenID uuid.UUID) error
	ValidateRefreshToken(ctx context.Context, tokenID uuid.UUID) (bool, error)

	SignInWithEmail(ctx context.Context, email, password string) (*models.AuthenticatedUser, error)
	SignInWithUsername(ctx context.Context, username, password string) (*models.AuthenticatedUser, error)
}

// Ensure AuthService implements AuthServiceInterface
var _ AuthServiceInterface = (*AuthService)(nil)

// AuthService implements user business logic using a UserRepositoryInterface.
type AuthService struct {
	authRepo           repository.AuthRepositoryInterface
	userService        user_service.UserServiceInterface
	secretKey          []byte                 // Secret key for signing JWTs
	accessTokenExpiry  time.Duration          // Access token expiration duration
	refreshTokenExpiry time.Duration          // Refresh token expiration duration
	signingAlg         jwa.SignatureAlgorithm // Signing algorithm (default: HS256)
}

type AuthServiceOpts struct {
	AuthRepo           repository.AuthRepositoryInterface
	UserService        user_service.UserServiceInterface
	SecretKey          []byte                 // Secret key for signing JWTs
	AccessTokenExpiry  time.Duration          // Access token expiration duration
	RefreshTokenExpiry time.Duration          // Refresh token expiration duration
	SigningAlg         jwa.SignatureAlgorithm // Signing algorithm (default: HS256)
}

// NewAuthService creates a new AuthService.
func NewAuthService(opts AuthServiceOpts) *AuthService {
	return &AuthService{
		authRepo:           opts.AuthRepo,
		userService:        opts.UserService,
		secretKey:          opts.SecretKey,
		accessTokenExpiry:  opts.AccessTokenExpiry,
		refreshTokenExpiry: opts.RefreshTokenExpiry,
		signingAlg:         opts.SigningAlg,
	}
}
