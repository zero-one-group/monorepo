package user

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/zero-one-group/go-modulith/internal/errors"
	"go.opentelemetry.io/otel"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUsers(ctx context.Context, filters UserFilters) (*PaginatedUsersResponse, error) {
	ctx, span := otel.Tracer("user").Start(ctx, "service.get_users")
	defer span.End()

	filters.SetDefaults()

	users, total, err := s.repo.GetUsers(ctx, filters)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	totalPages := int(math.Ceil(float64(total) / float64(filters.Limit)))

	return &PaginatedUsersResponse{
		Users:      userResponses,
		Total:      total,
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
	ctx, span := otel.Tracer("user").Start(ctx, "service.get_user_by_id")
	defer span.End()

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	if user == nil {
		return nil, errors.ErrNotFound
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	ctx, span := otel.Tracer("user").Start(ctx, "service.get_profile")
	defer span.End()

	return s.GetUserByID(ctx, userID)
}

func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserResponse, error) {
	ctx, span := otel.Tracer("user").Start(ctx, "service.update_user")
	defer span.End()

	updates := map[string]interface{}{
		"name": req.Name,
	}

	user, err := s.repo.UpdateUser(ctx, id, updates)
	if err != nil {
		span.RecordError(err)
		return nil, errors.ErrInternal
	}

	if user == nil {
		return nil, errors.ErrNotFound
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ctx, span := otel.Tracer("user").Start(ctx, "service.delete_user")
	defer span.End()

	if err := s.repo.DeleteUser(ctx, id); err != nil {
		if err.Error() == "record not found" {
			return errors.ErrNotFound
		}
		span.RecordError(err)
		return errors.ErrInternal
	}

	return nil
}