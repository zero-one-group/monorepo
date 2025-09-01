package utils

import (
	"log/slog"
	"os"
	"strconv"
	"time"
	"go-app/domain"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, email string) (string, string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	expiryTime, err := strconv.Atoi(os.Getenv("AUTH_TOKEN_EXPIRY_MINUTES"))
	if err != nil {
		expiryTime = 60
	}

	// Access Token
	claims := domain.JwtClaim{
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
	refreshClaims := domain.RefreshClaim{
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

func ValidateToken(tokenString string) (*domain.RefreshClaim, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &domain.RefreshClaim{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		slog.Error("Error parsing token")
		return nil, err
	}

	claims, ok := token.Claims.(*domain.RefreshClaim)
	if !ok {
		slog.Error("Error getting claims")
		return nil, err
	}

	return claims, nil
}
