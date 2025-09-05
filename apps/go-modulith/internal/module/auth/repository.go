package auth

import (
	"context"
	"errors"

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

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	ctx, span := otel.Tracer("auth").Start(ctx, "repository.create_user")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	ctx, span := otel.Tracer("auth").Start(ctx, "repository.get_user_by_email")
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

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	ctx, span := otel.Tracer("auth").Start(ctx, "repository.get_user_by_id")
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

func (r *Repository) CreateRefreshToken(ctx context.Context, refreshToken *RefreshToken) error {
	ctx, span := otel.Tracer("auth").Start(ctx, "repository.create_refresh_token")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(refreshToken).Error; err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *Repository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	ctx, span := otel.Tracer("auth").Start(ctx, "repository.get_refresh_token")
	defer span.End()

	var refreshToken RefreshToken
	if err := r.db.WithContext(ctx).Preload("User").Where("token = ?", token).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		span.RecordError(err)
		return nil, err
	}

	return &refreshToken, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, token string) error {
	ctx, span := otel.Tracer("auth").Start(ctx, "repository.delete_refresh_token")
	defer span.End()

	if err := r.db.WithContext(ctx).Where("token = ?", token).Delete(&RefreshToken{}).Error; err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *Repository) DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	ctx, span := otel.Tracer("auth").Start(ctx, "repository.delete_user_refresh_tokens")
	defer span.End()

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&RefreshToken{}).Error; err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}