package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Name      string         `json:"name" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginatedUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type UserFilters struct {
	Page   int    `query:"page" validate:"gte=1"`
	Limit  int    `query:"limit" validate:"gte=1,lte=100"`
	Search string `query:"search"`
	Sort   string `query:"sort" validate:"omitempty,oneof=name email created_at"`
	Order  string `query:"order" validate:"omitempty,oneof=asc desc"`
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (f *UserFilters) SetDefaults() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 10
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	if f.Sort == "" {
		f.Sort = "created_at"
	}
	if f.Order == "" {
		f.Order = "desc"
	}
}

func (f UserFilters) Offset() int {
	return (f.Page - 1) * f.Limit
}