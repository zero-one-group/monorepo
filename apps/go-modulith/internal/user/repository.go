package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/zero-one-group/go-modulith/internal/database"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUsers(ctx context.Context, filters UserFilters) ([]User, int64, error) {
	ctx, span := otel.Tracer("user").Start(ctx, "repository.get_users")
	defer span.End()

	var users []User
	var total int64

	query := r.db.WithContext(ctx).Model(&User{})

	if filters.Search != "" {
		searchPattern := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	orderClause := fmt.Sprintf("%s %s", filters.Sort, strings.ToUpper(filters.Order))
	if err := query.Order(orderClause).
		Offset(filters.Offset()).
		Limit(filters.Limit).
		Find(&users).Error; err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	return users, total, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	ctx, span := otel.Tracer("user").Start(ctx, "repository.get_user_by_id")
	defer span.End()

	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		span.RecordError(err)
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	ctx, span := otel.Tracer("user").Start(ctx, "repository.get_user_by_email")
	defer span.End()

	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		span.RecordError(err)
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*User, error) {
	ctx, span := otel.Tracer("user").Start(ctx, "repository.update_user")
	defer span.End()

	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		span.RecordError(err)
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&user).Updates(updates).Error; err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ctx, span := otel.Tracer("user").Start(ctx, "repository.delete_user")
	defer span.End()

	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&User{})
	if result.Error != nil {
		span.RecordError(result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}