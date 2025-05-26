package domain

import "time"


type User struct {
	ID                    string     `json:"id"`
	Name                  string     `json:"name" validate:"required"`
	Email                 string     `json:"email" validate:"required"`
	CreatedAt             *time.Time `json:"created_at"`
    UpdatedAt             *time.Time `json:"updated_at"`
    DeletedAt             *time.Time `json:"deleted_at"`
}

type UserFilter struct {
	Search string `json:"search" query:"search"`
}
