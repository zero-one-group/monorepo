package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null;index"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	CategoryID  *uuid.UUID     `json:"category_id,omitempty" gorm:"type:uuid;index"`
	CreatedBy   uuid.UUID      `json:"created_by" gorm:"type:uuid;not null;index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Category struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null;uniqueIndex"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateProductRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=200"`
	Description string     `json:"description" validate:"max=1000"`
	Price       float64    `json:"price" validate:"required,gt=0"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
}

type UpdateProductRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=200"`
	Description string     `json:"description" validate:"max=1000"`
	Price       float64    `json:"price" validate:"required,gt=0"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
}

type ProductResponse struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Category    *CategoryResponse `json:"category,omitempty"`
	CreatedBy   uuid.UUID        `json:"created_by"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type CategoryResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginatedProductsResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

type ProductFilters struct {
	Page       int        `query:"page" validate:"gte=1"`
	Limit      int        `query:"limit" validate:"gte=1,lte=100"`
	Search     string     `query:"search"`
	CategoryID *uuid.UUID `query:"category_id"`
	MinPrice   *float64   `query:"min_price" validate:"omitempty,gte=0"`
	MaxPrice   *float64   `query:"max_price" validate:"omitempty,gte=0"`
	Sort       string     `query:"sort" validate:"omitempty,oneof=name price created_at"`
	Order      string     `query:"order" validate:"omitempty,oneof=asc desc"`
}

func (p Product) ToResponse() ProductResponse {
	response := ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	return response
}

func (c Category) ToResponse() CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func (f *ProductFilters) SetDefaults() {
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

func (f ProductFilters) Offset() int {
	return (f.Page - 1) * f.Limit
}