package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User         User   `json:"user"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
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
