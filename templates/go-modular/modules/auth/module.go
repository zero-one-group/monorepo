package auth

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"{{ package_name | kebab_case }}/internal/notification"
	"{{ package_name | kebab_case }}/modules/auth/handler"
	"{{ package_name | kebab_case }}/modules/auth/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwa"

	svcUser "{{ package_name | kebab_case }}/modules/auth/services"
	svcAuth "{{ package_name | kebab_case }}/modules/user/services"
)

type Options struct {
	PgPool             *pgxpool.Pool // PostgreSQL connection pool (required)
	Logger             *slog.Logger  // Slog logger instance (optional)
	UserService        svcAuth.UserServiceInterface
	JWTSecretKey       []byte                 // Secret key for signing JWTs
	AccessTokenExpiry  time.Duration          // Access token expiration duration
	RefreshTokenExpiry time.Duration          // Refresh token expiration duration
	SigningAlg         jwa.SignatureAlgorithm // Signing algorithm (default: HS256)

	// Mailer dependency (optional). Provided mailer will be available to handlers.
	Mailer *notification.Mailer

	// BaseURL used when constructing verification links (MANDATORY).
	// Caller MUST provide a fully qualified base URL (e.g. https://example.com)
	// via Options.BaseURL before creating the module. We no longer read APP_BASE_URL here.
	BaseURL string
}

// AuthModule holds dependencies for auth-related handlers.
type AuthModule struct {
	logger      *slog.Logger
	middlewares []echo.MiddlewareFunc
	handler     *handler.Handler

	// mailer available to the module/handlers
	mailer *notification.Mailer

	// keep JWT config so we can attach middleware to protected routes
	jwtSecretKey []byte
	signingAlg   jwa.SignatureAlgorithm
}

// validateAndSetDefaults validates Options and sets defaults if needed.
func (opts *Options) validateAndSetDefaults() error {
	if opts.PgPool == nil {
		return fmt.Errorf("PgPool is required")
	}
	if opts.UserService == nil {
		return fmt.Errorf("UserService is required")
	}
	if len(opts.JWTSecretKey) == 0 {
		return fmt.Errorf("JWTSecretKey is required")
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

	// BaseURL is mandatory and must be provided via Options.BaseURL
	if opts.BaseURL == "" {
		return fmt.Errorf("BaseURL is required (set Options.BaseURL)")
	}

	return nil
}

// NewModule creates a new AuthModule.
func NewModule(opts *Options) *AuthModule {
	// Validate options and set defaults
	if err := opts.validateAndSetDefaults(); err != nil {
		panic("invalid auth module options: " + err.Error())
	}

	logger := opts.Logger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	if opts.AccessTokenExpiry == 0 {
		opts.AccessTokenExpiry = 24 * time.Hour // Set default access token expiry to 24 hours
	}
	if opts.RefreshTokenExpiry == 0 {
		opts.RefreshTokenExpiry = 7 * 24 * time.Hour // Set default refresh token expiry to 7 days
	}

	authRepo := repository.NewAuthRepository(opts.PgPool, logger)
	authService := svcUser.NewAuthService(svcUser.AuthServiceOpts{
		AuthRepo:           authRepo,
		UserService:        opts.UserService,
		JWTSecretKey:       opts.JWTSecretKey,
		AccessTokenExpiry:  opts.AccessTokenExpiry,
		RefreshTokenExpiry: opts.RefreshTokenExpiry,
		SigningAlg:         opts.SigningAlg,
		Mailer:             opts.Mailer,
		BaseURL:            opts.BaseURL,
	})

	h := handler.NewHandler(&handler.HandlerOpts{
		Logger:      logger,
		AuthService: authService,
	})

	return &AuthModule{
		logger:       logger,
		handler:      h,
		mailer:       opts.Mailer,
		jwtSecretKey: opts.JWTSecretKey,
		signingAlg:   opts.SigningAlg,
	}
}

// Use adds middleware(s) to the AuthModule (grouped).
func (m *AuthModule) Use(mw ...echo.MiddlewareFunc) {
	m.middlewares = append(m.middlewares, mw...)
}

// JWTMiddleware returns an echo.MiddlewareFunc configured with the module's secret and algorithm.
func (m *AuthModule) JWTMiddleware() echo.MiddlewareFunc {
	return JWTMiddleware(m.jwtSecretKey, m.signingAlg)
}

// RegisterRoutes registers auth endpoints to the given Echo group.
func (m *AuthModule) RegisterRoutes(e *echo.Group) {
	// Public routes (no access token required)
	publicGroup := e.Group("/auth", m.middlewares...)
	publicGroup.POST("/signin/email", m.handler.SignInWithEmail)
	publicGroup.POST("/signin/username", m.handler.SignInWithUsername)
	publicGroup.GET("/verify-email", m.handler.ValidateEmailVerificationByLink)
	publicGroup.POST("/refresh-token", m.handler.CreateRefreshToken)
	publicGroup.PUT("/refresh-token", m.handler.UpdateRefreshToken)
	publicGroup.GET("/refresh-token/:tokenId", m.handler.GetRefreshToken)
	publicGroup.DELETE("/refresh-token/:tokenId", m.handler.DeleteRefreshToken)
	publicGroup.POST("/verification/email/initiate", m.handler.InitiateEmailVerification)
	publicGroup.POST("/verification/email/validate", m.handler.ValidateEmailVerification)

	// Protected routes (require access token)
	protected := publicGroup.Group("", m.JWTMiddleware())
	protected.POST("/password", m.handler.SetUserPassword)
	protected.PUT("/password/:userId", m.handler.UpdateUserPassword)
	protected.POST("/session", m.handler.CreateSession)
	protected.PUT("/session", m.handler.UpdateSession)
	protected.GET("/session/:sessionId", m.handler.GetSession)
	protected.DELETE("/session/:sessionId", m.handler.DeleteSession)
	protected.POST("/verification/email/revoke", m.handler.RevokeEmailVerification)
	protected.POST("/verification/email/resend", m.handler.ResendEmailVerification)
}
