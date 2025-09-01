package service

import (
	"context"
	"go-app/domain"
	"go-app/utils"
)

type AuthRepository interface {
	AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error)
}

type AuthService struct {
	authRepo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{
		authRepo: repo,
	}
}

func (as *AuthService) Login(ctx context.Context, email, password string) (*domain.LoginResponse, error) {
	user, err := as.authRepo.AuthenticateUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	token, refreshToken, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		User:         *user,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}
