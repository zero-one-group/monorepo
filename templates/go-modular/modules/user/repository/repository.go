package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"{{ package_name | kebab_case }}/modules/user/models"
)

// Sentinel error for not found
var ErrNotFound = errors.New("not found")

// UserRepositoryInterface defines the contract for user data access.
type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListUsers(ctx context.Context, filter *models.FilterUser) ([]*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UsernameExists(ctx context.Context, username string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

// Ensure UserRepository implements UserRepositoryInterface
var _ UserRepositoryInterface = (*UserRepository)(nil)

// UserRepository is an implementation of UserRepositoryInterface using pgxpool.
type UserRepository struct {
	pgPool *pgxpool.Pool
	logger *slog.Logger
}

// NewUserRepository creates a new UserRepository with pgxpool and slog logger.
func NewUserRepository(pgPool *pgxpool.Pool, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		pgPool: pgPool,
		logger: logger,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.Must(uuid.NewV7())
	}
	user.CreatedAt = time.Now()

	tx, err := r.pgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Error("failed to begin transaction", slog.String("op", "CreateUser"), slog.String("error", err.Error()))
		return err
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != pgx.ErrTxClosed {
			r.logger.Warn("failed to rollback transaction", slog.String("op", "CreateUser"), slog.String("error", rbErr.Error()))
		}
	}()

	query := `
		INSERT INTO ` + models.UserTable + ` (
			id, display_name, email, username, avatar_url, metadata, created_at, updated_at, email_verified_at,
            last_login_at, banned_at, ban_expires, ban_reason
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`
	_, err = tx.Exec(ctx, query,
		user.ID,
		user.DisplayName,
		user.Email,
		user.Username,
		user.AvatarURL,
		user.Metadata,
		user.CreatedAt,
		user.UpdatedAt,
		user.EmailVerifiedAt,
		user.LastLoginAt,
		user.BannedAt,
		user.BanExpires,
		user.BanReason,
	)
	if err != nil {
		r.logger.Error("failed to insert user", slog.String("op", "CreateUser"), slog.String("error", err.Error()))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		r.logger.Error("failed to commit transaction", slog.String("op", "CreateUser"), slog.String("error", err.Error()))
		return err
	}
	r.logger.Info("user created", slog.String("op", "CreateUser"), slog.String("user_id", user.ID.String()))
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, display_name, email, username, avatar_url, metadata, created_at, updated_at, email_verified_at,
        last_login_at, banned_at, ban_expires, ban_reason
		FROM ` + models.UserTable + `
		WHERE id = $1
	`
	row := r.pgPool.QueryRow(ctx, query, id)
	var user models.User
	var metadataBytes []byte
	err := row.Scan(
		&user.ID,
		&user.DisplayName,
		&user.Email,
		&user.Username,
		&user.AvatarURL,
		&metadataBytes,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.EmailVerifiedAt,
		&user.LastLoginAt,
		&user.BannedAt,
		&user.BanExpires,
		&user.BanReason,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		r.logger.Error("failed to get user by id", slog.String("op", "GetUserByID"), slog.String("user_id", id.String()), slog.String("error", err.Error()))
		return nil, err
	}
	if len(metadataBytes) > 0 {
		var meta models.UserMetadata
		if err := json.Unmarshal(metadataBytes, &meta); err == nil {
			user.Metadata = &meta
		}
	}
	r.logger.Info("user fetched", slog.String("op", "GetUserByID"), slog.String("user_id", id.String()))
	return &user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, filter *models.FilterUser) ([]*models.User, error) {
	query := `
		SELECT id, display_name, email, username, avatar_url, metadata, created_at, updated_at, email_verified_at,
        last_login_at, banned_at, ban_expires, ban_reason
		FROM ` + models.UserTable
	args := []any{}
	whereClauses := []string{}
	argIdx := 1

	if filter != nil && filter.Search != nil && *filter.Search != "" {
		whereClauses = append(whereClauses, "(display_name ILIKE $"+itoa(argIdx)+" OR username ILIKE $"+itoa(argIdx)+")")
		args = append(args, "%"+*filter.Search+"%")
		argIdx++
	}
	if len(whereClauses) > 0 {
		query += " WHERE " + joinClauses(whereClauses, " AND ")
	}
	query += " ORDER BY created_at DESC"
	if filter != nil && filter.Limit > 0 {
		query += " LIMIT $" + itoa(argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}
	if filter != nil && filter.Offset > 0 {
		query += " OFFSET $" + itoa(argIdx)
		args = append(args, filter.Offset)
	}

	rows, err := r.pgPool.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to list users", slog.String("op", "ListUsers"), slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		var metadataBytes []byte
		err := rows.Scan(
			&user.ID,
			&user.DisplayName,
			&user.Email,
			&user.Username,
			&user.AvatarURL,
			&metadataBytes,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.EmailVerifiedAt,
			&user.LastLoginAt,
			&user.BannedAt,
			&user.BanExpires,
			&user.BanReason,
		)
		if err != nil {
			r.logger.Error("failed to scan user row", slog.String("op", "ListUsers"), slog.String("error", err.Error()))
			return nil, err
		}
		if len(metadataBytes) > 0 {
			var meta models.UserMetadata
			if err := json.Unmarshal(metadataBytes, &meta); err == nil {
				user.Metadata = &meta
			}
		}
		users = append(users, &user)
	}
	r.logger.Info("users listed", slog.String("op", "ListUsers"), slog.Int("count", len(users)))
	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = ptrTime(time.Now())

	tx, err := r.pgPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Error("failed to begin transaction", slog.String("op", "UpdateUser"), slog.String("error", err.Error()))
		return err
	}
	defer func() {
		if rbErr := tx.Rollback(ctx); rbErr != nil && rbErr != pgx.ErrTxClosed {
			r.logger.Warn("failed to rollback transaction", slog.String("op", "UpdateUser"), slog.String("error", rbErr.Error()))
		}
	}()

	query := `
		UPDATE ` + models.UserTable + `
		SET display_name = $1, email = $2, username = $3, avatar_url = $4, metadata = $5, updated_at = $6,
        email_verified_at = $7, last_login_at = $8, banned_at = $9, ban_expires = $10, ban_reason = $11
		WHERE id = $12
	`
	cmd, err := tx.Exec(ctx, query,
		user.DisplayName,
		user.Email,
		user.Username,
		user.AvatarURL,
		user.Metadata,
		user.UpdatedAt,
		user.EmailVerifiedAt,
		user.LastLoginAt,
		user.BannedAt,
		user.BanExpires,
		user.BanReason,
		user.ID,
	)
	if err != nil {
		r.logger.Error("failed to update user", slog.String("op", "UpdateUser"), slog.String("user_id", user.ID.String()), slog.String("error", err.Error()))
		return err
	}
	if cmd.RowsAffected() == 0 {
		r.logger.Warn("user not found for update", slog.String("op", "UpdateUser"), slog.String("user_id", user.ID.String()))
		return ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		r.logger.Error("failed to commit transaction", slog.String("op", "UpdateUser"), slog.String("user_id", user.ID.String()), slog.String("error", err.Error()))
		return err
	}
	r.logger.Info("user updated", slog.String("op", "UpdateUser"), slog.String("user_id", user.ID.String()))
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM ` + models.UserTable + ` WHERE id = $1`
	cmd, err := r.pgPool.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete user", slog.String("op", "DeleteUser"), slog.String("user_id", id.String()), slog.String("error", err.Error()))
		return err
	}
	if cmd.RowsAffected() == 0 {
		r.logger.Warn("user not found for delete", slog.String("op", "DeleteUser"), slog.String("user_id", id.String()))
		return ErrNotFound
	}
	r.logger.Info("user deleted", slog.String("op", "DeleteUser"), slog.String("user_id", id.String()))
	return nil
}

// Check if a username exists (case-insensitive)
func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	query := `SELECT 1 FROM ` + models.UserTable + ` WHERE LOWER(username) = LOWER($1) LIMIT 1`
	var exists int
	err := r.pgPool.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		r.logger.Error("failed to check username exists", slog.String("op", "UsernameExists"), slog.String("error", err.Error()))
		return false, err
	}
	return true, nil
}

// Check if an email exists (case-insensitive)
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT 1 FROM ` + models.UserTable + ` WHERE LOWER(email) = LOWER($1) LIMIT 1`
	var exists int
	err := r.pgPool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		r.logger.Error("failed to check email exists", slog.String("op", "EmailExists"), slog.String("error", err.Error()))
		return false, err
	}
	return true, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, display_name, email, username, avatar_url, metadata, created_at, updated_at,
        email_verified_at, last_login_at, banned_at, ban_expires, ban_reason
        FROM ` + models.UserTable + `
        WHERE LOWER(email) = LOWER($1)
        LIMIT 1
    `
	row := r.pgPool.QueryRow(ctx, query, email)
	var user models.User
	var metadataBytes []byte
	err := row.Scan(
		&user.ID,
		&user.DisplayName,
		&user.Email,
		&user.Username,
		&user.AvatarURL,
		&metadataBytes,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.EmailVerifiedAt,
		&user.LastLoginAt,
		&user.BannedAt,
		&user.BanExpires,
		&user.BanReason,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		r.logger.Error("failed to get user by email", slog.String("op", "GetUserByEmail"), slog.String("email", email), slog.String("error", err.Error()))
		return nil, err
	}
	if len(metadataBytes) > 0 {
		var meta models.UserMetadata
		if err := json.Unmarshal(metadataBytes, &meta); err == nil {
			user.Metadata = &meta
		}
	}
	r.logger.Info("user fetched", slog.String("op", "GetUserByEmail"), slog.String("user_id", user.ID.String()))
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
        SELECT id, display_name, email, username, avatar_url, metadata, created_at, updated_at,
        email_verified_at, last_login_at, banned_at, ban_expires, ban_reason
        FROM ` + models.UserTable + `
        WHERE LOWER(username) = LOWER($1)
        LIMIT 1
    `
	row := r.pgPool.QueryRow(ctx, query, username)
	var user models.User
	var metadataBytes []byte
	err := row.Scan(
		&user.ID,
		&user.DisplayName,
		&user.Email,
		&user.Username,
		&user.AvatarURL,
		&metadataBytes,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.EmailVerifiedAt,
		&user.LastLoginAt,
		&user.BannedAt,
		&user.BanExpires,
		&user.BanReason,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		r.logger.Error("failed to get user by username", slog.String("op", "GetUserByUsername"), slog.String("username", username), slog.String("error", err.Error()))
		return nil, err
	}
	if len(metadataBytes) > 0 {
		var meta models.UserMetadata
		if err := json.Unmarshal(metadataBytes, &meta); err == nil {
			user.Metadata = &meta
		}
	}
	r.logger.Info("user fetched", slog.String("op", "GetUserByUsername"), slog.String("user_id", user.ID.String()))
	return &user, nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

// Helper functions for dynamic query building
func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func joinClauses(clauses []string, sep string) string {
	if len(clauses) == 0 {
		return ""
	}
	result := clauses[0]
	for i := 1; i < len(clauses); i++ {
		result += sep + clauses[i]
	}
	return result
}
