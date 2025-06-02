package domain

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name  string    `json:"name" validate:"required"`
	Email string    `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,password"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UserFilter struct {
	Search string `json:"search" query:"search"`
}
