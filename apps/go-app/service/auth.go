package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"
	"go-app/config"
	"go-app/domain"

	"github.com/golang-jwt/jwt/v5"
)

type AuthRepository interface {
	AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error)
}

type JwtClaim struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type RefreshClaim struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService struct {
	authRepo  AuthRepository
	jwtSecret []byte
	jwtTTL    time.Duration
}

func NewAuthService(repo AuthRepository) *AuthService {
	jwtSecret, jwtTTL := config.LoadJWTConfig()
	return &AuthService{
		authRepo:  repo,
		jwtSecret: []byte(jwtSecret),
		jwtTTL:    jwtTTL,
	}
}

func (as *AuthService) Login(ctx context.Context, email, password string) (*domain.LoginResponse, error) {
	user, err := as.authRepo.AuthenticateUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	token, refreshToken, err := as.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		User:         *user,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (as *AuthService) generateJWT(userID string, email string) (string, string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	expiryTime, err := strconv.Atoi(os.Getenv("AUTH_TOKEN_EXPIRY_MINUTES"))
	if err != nil {
		expiryTime = 60
	}

	// Access Token
	claims := JwtClaim{
		ID:    userID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expiryTime))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(secret)
	if err != nil {
		slog.Error("Error signing access token")
		return "", "", err
	}

	// Refresh Token (24 hours)
	refreshClaims := RefreshClaim{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		slog.Error("Error signing refresh token")
		return "", "", err
	}

	return accessToken, refreshTokenString, nil
}

func (as *AuthService) ValidateToken(tokenString string) (string, error) {
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, errors.New("unexpected signing method")
	// 	}
	// 	return as.jwtSecret, nil
	// })

	// if err != nil {
	// 	return "", err
	// }

	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	// 	userID, ok := claims["user_id"].(string)
	// 	if !ok {
	// 		return "", errors.New("invalid token claims")
	// 	}
	// 	return userID, nil
	// }

	return "", errors.New("invalid token")
}
