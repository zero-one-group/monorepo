package auth

import (
	"fmt"
	"log/slog"
	"strings"

	"go-modular/internal/config"
	appMiddleware "go-modular/internal/middleware"
	"go-modular/internal/notification"
	"go-modular/internal/server"
	"go-modular/modules/auth/handler"
	"go-modular/modules/auth/repository"
	"go-modular/modules/auth/services"
	userServices "go-modular/modules/user/services"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"go.uber.org/fx"
)

// Module wires the Auth feature.
var Module = fx.Module("auth",
	fx.Provide(
		repository.NewAuthRepository,
		newAuthService,
		newAuthHandler,
		newJWTAccess,
	),
	fx.Invoke(registerRoutes),
)

func newAuthService(
	repo *repository.AuthRepository,
	users userServices.UserServiceInterface,
	mailer *notification.Mailer,
	cfg *config.Config,
) (*services.AuthService, error) {
	if err := validateAuthConfig(cfg); err != nil {
		return nil, err
	}
	return services.NewAuthService(services.AuthServiceOpts{
		AuthRepo:     repo,
		UserService:  users,
		JWTSecretKey: []byte(cfg.App.JWTSecretKey),
		SigningAlg:   jwtSigningAlg(cfg),
		Mailer:       mailer,
		BaseURL:      cfg.GetAppBaseURL(),
	}), nil
}

func newAuthHandler(log *slog.Logger, svc *services.AuthService) *handler.Handler {
	return handler.NewHandler(&handler.HandlerOpts{
		Logger:      log,
		AuthService: svc,
	})
}

func newJWTAccess(cfg *config.Config) (appMiddleware.JWTAccess, error) {
	if err := validateAuthConfig(cfg); err != nil {
		return nil, err
	}
	return appMiddleware.JWTAccess(JWTMiddleware([]byte(cfg.App.JWTSecretKey), jwtSigningAlg(cfg))), nil
}

func validateAuthConfig(cfg *config.Config) error {
	if strings.TrimSpace(cfg.App.JWTSecretKey) == "" {
		return fmt.Errorf("JWTSecretKey is required")
	}
	if strings.TrimSpace(cfg.GetAppBaseURL()) == "" {
		return fmt.Errorf("BaseURL is required")
	}
	return nil
}

func jwtSigningAlg(cfg *config.Config) jwa.SignatureAlgorithm {
	alg := strings.ToUpper(strings.TrimSpace(string(cfg.App.JWTAlgorithm)))
	if alg == "" {
		return jwa.HS256
	}
	return jwa.SignatureAlgorithm(alg)
}

func registerRoutes(api server.APIV1, h *handler.Handler, jwt appMiddleware.JWTAccess) {
	publicGroup := api.Root.Group("/auth")
	publicGroup.POST("/signin/email", h.SignInWithEmail)
	publicGroup.POST("/signin/username", h.SignInWithUsername)
	publicGroup.GET("/verify-email", h.ValidateEmailVerificationByLink)
	publicGroup.POST("/refresh-token", h.CreateRefreshToken)
	publicGroup.PUT("/refresh-token", h.UpdateRefreshToken)
	publicGroup.GET("/refresh-token/:tokenId", h.GetRefreshToken)
	publicGroup.DELETE("/refresh-token/:tokenId", h.DeleteRefreshToken)
	publicGroup.POST("/verification/email/initiate", h.InitiateEmailVerification)
	publicGroup.POST("/verification/email/validate", h.ValidateEmailVerification)

	protected := publicGroup.Group("", echo.MiddlewareFunc(jwt))
	protected.POST("/password", h.SetUserPassword)
	protected.PUT("/password/:userId", h.UpdateUserPassword)
	protected.POST("/session", h.CreateSession)
	protected.PUT("/session", h.UpdateSession)
	protected.GET("/session/:sessionId", h.GetSession)
	protected.DELETE("/session/:sessionId", h.DeleteSession)
	protected.POST("/verification/email/revoke", h.RevokeEmailVerification)
	protected.POST("/verification/email/resend", h.ResendEmailVerification)
}
