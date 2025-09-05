package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zero-one-group/go-modulith/internal/config"
	"github.com/zero-one-group/go-modulith/internal/errors"
	"github.com/zero-one-group/go-modulith/internal/middleware"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
	cfg  *config.Config
}

func NewService(repo *Repository, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	ctx, span := otel.Tracer("auth").Start(ctx, "service.register")
	defer span.End()

	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	if existingUser != nil {
		return nil, errors.ErrConflict.WithDetails(map[string]string{
			"email": "Email already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	user := &User{
		ID:       uuid.New(),
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	refreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	return &AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	ctx, span := otel.Tracer("auth").Start(ctx, "service.login")
	defer span.End()

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	if user == nil {
		return nil, errors.ErrUnauthorized.WithDetails(map[string]string{
			"credentials": "Invalid email or password",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.ErrUnauthorized.WithDetails(map[string]string{
			"credentials": "Invalid email or password",
		})
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	refreshToken, err := s.generateRefreshToken(ctx, user.ID)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	return &AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*AuthResponse, error) {
	ctx, span := otel.Tracer("auth").Start(ctx, "service.refresh_token")
	defer span.End()

	refreshToken, err := s.repo.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	if refreshToken == nil || refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.ErrUnauthorized.WithDetails(map[string]string{
			"refresh_token": "Invalid or expired refresh token",
		})
	}

	if err := s.repo.DeleteRefreshToken(ctx, req.RefreshToken); err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	accessToken, err := s.generateAccessToken(&refreshToken.User)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	newRefreshToken, err := s.generateRefreshToken(ctx, refreshToken.User.ID)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	return &AuthResponse{
		User:         refreshToken.User.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, userID string) error {
	ctx, span := otel.Tracer("auth").Start(ctx, "service.logout")
	defer span.End()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.ErrBadRequest
	}

	if err := s.repo.DeleteUserRefreshTokens(ctx, userUUID); err != nil {
		span.RecordError(err)
		return errors.ErrInternal
	}

	return nil
}

func (s *Service) generateAccessToken(user *User) (string, error) {
	claims := &middleware.Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWT.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

func (s *Service) generateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}

	tokenString := hex.EncodeToString(tokenBytes)

	refreshToken := &RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(s.cfg.JWT.RefreshExpiry),
	}

	if err := s.repo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return "", err
	}

	return tokenString, nil
}