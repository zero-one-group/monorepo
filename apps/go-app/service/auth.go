package service

import (
	"context"
	"errors"
	"go-app/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthRepository interface {
	AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error)
}

type AuthService struct {
	authRepo  AuthRepository
	jwtSecret []byte
	jwtTTL    time.Duration
}

func NewAuthService(repo AuthRepository, jwtSecret string, jwtTTL time.Duration) *AuthService {
	return &AuthService{
		authRepo:  repo,
		jwtSecret: []byte(jwtSecret),
		jwtTTL:    jwtTTL,
	}
}

func (as *AuthService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := as.authRepo.AuthenticateUser(ctx, email, password)
	if err != nil {
		return "", nil, err
	}

	token, err := as.generateJWT(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (as *AuthService) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(as.jwtTTL).Unix(),
		// "iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(as.jwtSecret)
}

func (as *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return as.jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid token claims")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}
