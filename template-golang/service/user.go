package service

import (
	"context"
	"fmt"

	"{{ package_name }}/domain"
)


type UserRepository interface {
	GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error)
}

type UserService struct {
	userRepo              UserRepository
}

func NewUserService(u UserRepository) *UserService {
	return &UserService{
		userRepo:       u,
	}
}

func (us *UserService) GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error) {
	users, err := us.userRepo.GetUserList(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return users,  nil
}

