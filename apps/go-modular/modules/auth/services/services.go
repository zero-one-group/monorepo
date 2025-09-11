package services

import (
	"context"
	"time"

	"go-modular/internal/notification"
	"go-modular/modules/auth/models"
	"go-modular/modules/auth/repository"

	"github.com/gofrs/uuid/v5"
	"github.com/lestrrat-go/jwx/jwa"

	svcUser "go-modular/modules/user/services"
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

	// Authentication
	SignInWithEmail(ctx context.Context, email, password string) (*models.AuthenticatedUser, error)
	SignInWithUsername(ctx context.Context, username, password string) (*models.AuthenticatedUser, error)

	// Account verification (email-based, userID resolved internally)
	// Initiate/Resend now accept an optional redirectTo which will be stored in token metadata.
	InitiateEmailVerification(ctx context.Context, email, redirectTo string) error
	ValidateEmailVerification(ctx context.Context, token string) (bool, error)
	RevokeEmailVerification(ctx context.Context, token string) error
	ResendEmailVerification(ctx context.Context, email, redirectTo string) error
}

// Ensure AuthService implements AuthServiceInterface
var _ AuthServiceInterface = (*AuthService)(nil)

// AuthService implements user business logic using a UserRepositoryInterface.
type AuthService struct {
	authRepo           repository.AuthRepositoryInterface
	userService        svcUser.UserServiceInterface
	secretKey          []byte                 // Secret key for signing JWTs
	accessTokenExpiry  time.Duration          // Access token expiration duration
	refreshTokenExpiry time.Duration          // Refresh token expiration duration
	signingAlg         jwa.SignatureAlgorithm // Signing algorithm (default: HS256)
	mailer             *notification.Mailer
	baseURL            string // Base URL used when constructing verification links
}

type AuthServiceOpts struct {
	AuthRepo           repository.AuthRepositoryInterface
	UserService        svcUser.UserServiceInterface
	JWTSecretKey       []byte                 // Secret key for signing JWTs
	AccessTokenExpiry  time.Duration          // Access token expiration duration
	RefreshTokenExpiry time.Duration          // Refresh token expiration duration
	SigningAlg         jwa.SignatureAlgorithm // Signing algorithm (default: HS256)
	Mailer             *notification.Mailer   // Mailer service for sending emails
	BaseURL            string                 // BaseURL used when constructing verification links (MANDATORY).
}

// NewAuthService creates a new AuthService.
func NewAuthService(opts AuthServiceOpts) *AuthService {
	if opts.AuthRepo == nil {
		panic("AuthRepo is required")
	}
	if opts.UserService == nil {
		panic("UserService is required")
	}
	if len(opts.JWTSecretKey) == 0 {
		panic("JWTSecretKey is required")
	}
	if opts.SigningAlg == "" {
		opts.SigningAlg = jwa.HS256
	}
	if opts.AccessTokenExpiry == 0 {
		opts.AccessTokenExpiry = 24 * time.Hour
	}
	if opts.RefreshTokenExpiry == 0 {
		opts.RefreshTokenExpiry = 7 * 24 * time.Hour
	}

	// BaseURL is mandatory
	if opts.BaseURL == "" {
		// keep behavior consistent with other validations: panic on missing required option
		panic("BaseURL is required")
	}

	return &AuthService{
		authRepo:           opts.AuthRepo,
		userService:        opts.UserService,
		secretKey:          opts.JWTSecretKey,
		accessTokenExpiry:  opts.AccessTokenExpiry,
		refreshTokenExpiry: opts.RefreshTokenExpiry,
		signingAlg:         opts.SigningAlg,
		mailer:             opts.Mailer,
		baseURL:            opts.BaseURL,
	}
}
