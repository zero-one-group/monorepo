package service

import (
	"context"
	"log/slog"
	"{{ package_name }}/domain"
	apperrors "{{ package_name }}/internal/errors"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, error)
	GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserService struct {
	userRepo UserRepository
}

func NewUserService(u UserRepository) *UserService {
	return &UserService{
		userRepo: u,
	}
}

// CreateUser adds a new user.
func (us *UserService) CreateUser(
	ctx context.Context,
	u *domain.CreateUserRequest,
) (*domain.User, error) {
	tracer := otel.Tracer("service.user")
	ctx, span := tracer.Start(ctx, "UserService.CreateUser")
	defer span.End()

	slog.InfoContext(ctx, "Creating new user",
		slog.String("user_email", u.Email),
		slog.String("user_name", u.Name),
	)

	createdUser, err := us.userRepo.CreateUser(ctx, u)
	if err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to create user",
			slog.String("error", err.Error()),
			slog.String("user_email", u.Email),
		)
		
		// If it's already an AppError, just return it with trace info
		if appErr, ok := apperrors.GetAppError(err); ok {
			return nil, appErr.WithTrace(span)
		}
		
		// Wrap unexpected errors
		return nil, apperrors.NewInternalError("User creation failed", err).WithTrace(span)
	}

	slog.InfoContext(ctx, "User created successfully",
		slog.String("user_id", createdUser.ID),
		slog.String("user_email", createdUser.Email),
	)

	return createdUser, nil
}

// GetUser fetches a user by ID.
func (us *UserService) GetUser(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {
	tracer := otel.Tracer("service.user")
	ctx, span := tracer.Start(ctx, "UserService.GetUser")
	defer span.End()

	slog.InfoContext(ctx, "Fetching user",
		slog.String("user_id", id.String()),
	)

	user, err := us.userRepo.GetUser(ctx, id)
	if err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to get user",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		
		// If it's already an AppError, just return it with trace info
		if appErr, ok := apperrors.GetAppError(err); ok {
			return nil, appErr.WithTrace(span)
		}
		
		// Wrap unexpected errors
		return nil, apperrors.NewInternalError("User retrieval failed", err).WithTrace(span)
	}

	slog.InfoContext(ctx, "User fetched successfully",
		slog.String("user_id", user.ID),
	)

	return user, nil
}

// UpdateUser updates name/email of an existing user.
func (us *UserService) UpdateUser(
	ctx context.Context,
	id uuid.UUID,
	u *domain.User,
) (*domain.User, error) {
	tracer := otel.Tracer("service.user")
	ctx, span := tracer.Start(ctx, "UserService.UpdateUser")
	defer span.End()

	slog.InfoContext(ctx, "Updating user",
		slog.String("user_id", id.String()),
		slog.String("new_name", u.Name),
		slog.String("new_email", u.Email),
	)

	// Get existing user first to validate it exists
	existing, err := us.userRepo.GetUser(ctx, id)
	if err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to get user for update",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		
		// If it's already an AppError, just return it with trace info
		if appErr, ok := apperrors.GetAppError(err); ok {
			return nil, appErr.WithTrace(span)
		}
		
		// Wrap unexpected errors
		return nil, apperrors.NewInternalError("User update failed", err).WithTrace(span)
	}

	if existing == nil {
		return nil, domain.NewUserNotFoundError(id.String()).WithTrace(span)
	}

	// Update fields
	existing.Name = u.Name
	existing.Email = u.Email

	updatedUser, err := us.userRepo.UpdateUser(ctx, id, existing)
	if err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to update user",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		
		// If it's already an AppError, just return it with trace info
		if appErr, ok := apperrors.GetAppError(err); ok {
			return nil, appErr.WithTrace(span)
		}
		
		// Wrap unexpected errors
		return nil, apperrors.NewInternalError("User update failed", err).WithTrace(span)
	}

	slog.InfoContext(ctx, "User updated successfully",
		slog.String("user_id", updatedUser.ID),
		slog.String("user_email", updatedUser.Email),
	)

	return updatedUser, nil
}

// DeleteUser removes a user by ID.
func (us *UserService) DeleteUser(
	ctx context.Context,
	id uuid.UUID,
) error {
	tracer := otel.Tracer("service.user")
	ctx, span := tracer.Start(ctx, "UserService.DeleteUser")
	defer span.End()

	slog.InfoContext(ctx, "Deleting user",
		slog.String("user_id", id.String()),
	)

	// Check if user exists first
	user, err := us.userRepo.GetUser(ctx, id)
	if err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to get user for deletion",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		
		// If it's already an AppError, just return it with trace info
		if appErr, ok := apperrors.GetAppError(err); ok {
			return appErr.WithTrace(span)
		}
		
		// Wrap unexpected errors
		return apperrors.NewInternalError("User deletion failed", err).WithTrace(span)
	}

	if user == nil {
		return domain.NewUserNotFoundError(id.String()).WithTrace(span)
	}

	// Perform deletion
	err = us.userRepo.DeleteUser(ctx, id)
	if err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to delete user",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
		)
		
		// If it's already an AppError, just return it with trace info
		if appErr, ok := apperrors.GetAppError(err); ok {
			return appErr.WithTrace(span)
		}
		
		// Wrap unexpected errors
		return apperrors.NewInternalError("User deletion failed", err).WithTrace(span)
	}

	slog.InfoContext(ctx, "User deleted successfully",
		slog.String("user_id", id.String()),
	)

	return nil
}

func (us *UserService) GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error) {
	tracer := otel.Tracer("service.user")
	ctx, span := tracer.Start(ctx, "UserService.GetUserList")
	defer span.End()

	slog.InfoContext(ctx, "Fetching user list",
		slog.Bool("filtered", filter != nil && filter.Search != ""),
	)

	users, err := us.userRepo.GetUserList(ctx, filter)
	if err != nil {
		span.RecordError(err)
		
		slog.ErrorContext(ctx, "Failed to get user list",
			slog.String("error", err.Error()),
		)
		
		// If it's already an AppError, just return it with trace info
		if appErr, ok := apperrors.GetAppError(err); ok {
			return nil, appErr.WithTrace(span)
		}
		
		// Wrap unexpected errors
		return nil, apperrors.NewInternalError("User list retrieval failed", err).WithTrace(span)
	}

	slog.InfoContext(ctx, "User list fetched successfully",
		slog.Int("count", len(users)),
	)

	return users, nil
}
