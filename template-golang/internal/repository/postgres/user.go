package postgres

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
	"{{ package_name }}/domain"
	"{{ package_name }}/internal/metrics"
	"{{ package_name }}/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type UserRepository struct {
	Conn    *pgxpool.Pool
	Metrics *metrics.Metrics
}

func NewUserRepository(conn *pgxpool.Pool, metrics *metrics.Metrics) *UserRepository {
	return &UserRepository{
		Conn:    conn,
		Metrics: metrics,
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, error) {
	tracer := otel.Tracer("repo.user")
	ctx, span := tracer.Start(ctx, "UserRepository.CreateUser")
	defer span.End()

	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		span.RecordError(err)
		u.Metrics.UserRepoCalls.WithLabelValues("CreateUser", "error").Inc()
		return nil, domain.NewUserInternalError("create_user", err)
	}

	span.SetAttributes(
		attribute.String("user.email", user.Email),
		attribute.String("user.name", user.Name),
	)

	// Execute query
	var id uuid.UUID
	var createdUser domain.User
	err = u.Conn.QueryRow(ctx, query, user.Name, user.Email, hashedPassword).
		Scan(&id, &createdUser.CreatedAt, &createdUser.UpdatedAt)
	
	if err != nil {
		span.RecordError(err)
		u.Metrics.UserRepoCalls.WithLabelValues("CreateUser", "error").Inc()
		
		slog.ErrorContext(ctx, "Failed to create user in database",
			slog.String("error", err.Error()),
			slog.String("user_email", user.Email),
			slog.String("operation", "create_user"),
		)
		
		// Check for specific database errors (e.g., unique constraint violation)
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, domain.NewUserConflictError(user.Email)
		}
		
		return nil, domain.NewUserDatabaseError("create_user", err)
	}

	createdUser.ID = id.String()
	createdUser.Name = user.Name
	createdUser.Email = user.Email

	u.Metrics.UserRepoCalls.WithLabelValues("CreateUser", "success").Inc()
	
	slog.InfoContext(ctx, "User created successfully",
		slog.String("user_id", createdUser.ID),
		slog.String("user_email", user.Email),
	)

	return &createdUser, nil
}

func (u *UserRepository) GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error) {
	tracer := otel.Tracer("repo.user")
	ctx, span := tracer.Start(ctx, "UserRepository.GetUserList")
	defer span.End()

	query := `
		SELECT
			u.id,
			u.name,
			u.email,
			u.created_at,
			u.updated_at
		FROM users u
		WHERE u.deleted_at IS NULL`

	var args []interface{}
	var conditions []string
	
	if filter != nil && filter.Search != "" {
		conditions = append(conditions, ` AND (u.name ILIKE $1 OR u.email ILIKE $1)`)
		args = append(args, "%"+filter.Search+"%")
		span.SetAttributes(attribute.String("filter.search", filter.Search))
	}

	if len(conditions) > 0 {
		query += strings.Join(conditions, "")
	}

	query += " ORDER BY u.created_at DESC"

	rows, err := u.Conn.Query(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		u.Metrics.UserRepoCalls.WithLabelValues("GetUserList", "error").Inc()
		
		slog.ErrorContext(ctx, "Failed to query user list from database",
			slog.String("error", err.Error()),
			slog.String("operation", "get_user_list"),
		)
		
		return nil, domain.NewUserDatabaseError("get_user_list", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			span.RecordError(err)
			u.Metrics.UserRepoCalls.WithLabelValues("GetUserList", "error").Inc()
			
			slog.ErrorContext(ctx, "Failed to scan user row",
				slog.String("error", err.Error()),
				slog.String("operation", "get_user_list"),
			)
			
			return nil, domain.NewUserDatabaseError("get_user_list", err)
		}
		users = append(users, user)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		span.RecordError(err)
		u.Metrics.UserRepoCalls.WithLabelValues("GetUserList", "error").Inc()
		
		slog.ErrorContext(ctx, "Error iterating over user rows",
			slog.String("error", err.Error()),
			slog.String("operation", "get_user_list"),
		)
		
		return nil, domain.NewUserDatabaseError("get_user_list", err)
	}

	u.Metrics.UserRepoCalls.WithLabelValues("GetUserList", "success").Inc()
	
	slog.InfoContext(ctx, "User list retrieved successfully",
		slog.Int("count", len(users)),
		slog.Bool("filtered", filter != nil && filter.Search != ""),
	)

	return users, nil
}

func (u *UserRepository) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	tracer := otel.Tracer("repo.user")
	ctx, span := tracer.Start(ctx, "UserRepository.GetUser")
	defer span.End()

	query := `
		SELECT
			id,
			name,
			email,
			created_at,
			updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	span.SetAttributes(
		attribute.String("user.id", id.String()),
		attribute.String("query.statement", query),
	)

	row := u.Conn.QueryRow(ctx, query, id)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		span.RecordError(err)
		u.Metrics.UserRepoCalls.WithLabelValues("GetUser", "error").Inc()
		
		if err == sql.ErrNoRows {
			slog.WarnContext(ctx, "User not found",
				slog.String("user_id", id.String()),
				slog.String("operation", "get_user"),
			)
			return nil, domain.NewUserNotFoundError(id.String())
		}
		
		slog.ErrorContext(ctx, "Failed to get user from database",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
			slog.String("operation", "get_user"),
		)
		
		return nil, domain.NewUserDatabaseError("get_user", err)
	}

	u.Metrics.UserRepoCalls.WithLabelValues("GetUser", "success").Inc()
	
	slog.InfoContext(ctx, "User retrieved successfully",
		slog.String("user_id", user.ID),
	)
	
	return &user, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, user *domain.User) (*domain.User, error) {
	tracer := otel.Tracer("repo.user")
	ctx, span := tracer.Start(ctx, "UserRepository.UpdateUser")
	defer span.End()

	// Optimized query to return updated data in single query
	query := `
		UPDATE users
		SET name = $1,
			email = $2,
			updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING id, name, email, created_at, updated_at`

	span.SetAttributes(
		attribute.String("user.id", id.String()),
		attribute.String("user.name", user.Name),
		attribute.String("user.email", user.Email),
	)

	var updatedUser domain.User
	err := u.Conn.QueryRow(ctx, query, user.Name, user.Email, id).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		span.RecordError(err)
		u.Metrics.UserRepoCalls.WithLabelValues("UpdateUser", "error").Inc()
		
		if err == sql.ErrNoRows {
			slog.WarnContext(ctx, "User not found for update",
				slog.String("user_id", id.String()),
				slog.String("operation", "update_user"),
			)
			return nil, domain.NewUserNotFoundError(id.String())
		}
		
		slog.ErrorContext(ctx, "Failed to update user in database",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
			slog.String("operation", "update_user"),
		)
		
		// Check for specific database errors (e.g., unique constraint violation)
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, domain.NewUserConflictError(user.Email)
		}
		
		return nil, domain.NewUserDatabaseError("update_user", err)
	}

	u.Metrics.UserRepoCalls.WithLabelValues("UpdateUser", "success").Inc()
	
	slog.InfoContext(ctx, "User updated successfully",
		slog.String("user_id", updatedUser.ID),
		slog.String("user_email", updatedUser.Email),
	)

	return &updatedUser, nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	tracer := otel.Tracer("repo.user")
	ctx, span := tracer.Start(ctx, "UserRepository.DeleteUser")
	defer span.End()

	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	span.SetAttributes(attribute.String("user.id", id.String()))

	result, err := u.Conn.Exec(ctx, query, id)
	if err != nil {
		span.RecordError(err)
		u.Metrics.UserRepoCalls.WithLabelValues("DeleteUser", "error").Inc()
		
		slog.ErrorContext(ctx, "Failed to delete user from database",
			slog.String("error", err.Error()),
			slog.String("user_id", id.String()),
			slog.String("operation", "delete_user"),
		)
		
		return domain.NewUserDatabaseError("delete_user", err)
	}

	// Check if any rows were affected
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		slog.WarnContext(ctx, "User not found for deletion",
			slog.String("user_id", id.String()),
			slog.String("operation", "delete_user"),
		)
		return domain.NewUserNotFoundError(id.String())
	}

	u.Metrics.UserRepoCalls.WithLabelValues("DeleteUser", "success").Inc()
	
	slog.InfoContext(ctx, "User deleted successfully",
		slog.String("user_id", id.String()),
	)

	return nil
}
