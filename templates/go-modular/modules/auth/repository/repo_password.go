package repository

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"{{ package_name | kebab_case }}/modules/auth/models"
	"{{ package_name | kebab_case }}/pkg/apputils"
)

func (r *AuthRepository) SetUserPassword(ctx context.Context, userPassword *models.UserPassword) error {
	if userPassword.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}
	userPassword.CreatedAt = time.Now()

	query := `INSERT INTO ` + models.UserPasswordTable + ` (user_id, password_hash, created_at) VALUES ($1, $2, $3)`
	_, err := r.pgPool.Exec(ctx, query,
		userPassword.UserID,
		[]byte(userPassword.PasswordHash),
		userPassword.CreatedAt,
	)

	if err != nil {
		r.logger.Error("failed to insert user password", "op", "SetUserPassword", "user_id", userPassword.UserID.String(), "error", err.Error())
		return err
	}
	r.logger.Info("user password created", "op", "SetUserPassword", "user_id", userPassword.UserID.String())

	return nil
}

func (r *AuthRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	now := time.Now()
	query := `UPDATE ` + models.UserPasswordTable + ` SET password_hash = $1, updated_at = $2 WHERE user_id = $3`
	cmd, err := r.pgPool.Exec(ctx, query, []byte(newPassword), now, userID)
	if err != nil {
		r.logger.Error("failed to update user password", "op", "UpdateUserPassword", "user_id", userID.String(), "error", err.Error())
		return err
	}

	if cmd.RowsAffected() == 0 {
		r.logger.Warn("user password not found for update", "op", "UpdateUserPassword", "user_id", userID.String())
		return ErrNotFound
	}
	r.logger.Info("user password updated", "op", "UpdateUserPassword", "user_id", userID.String())

	return nil
}

func (r *AuthRepository) ValidateUserPassword(ctx context.Context, userID uuid.UUID, password string) (bool, error) {
	var passwordHash []byte
	query := `SELECT password_hash FROM ` + models.UserPasswordTable + ` WHERE user_id = $1`
	err := r.pgPool.QueryRow(ctx, query, userID).Scan(&passwordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("user password not found", "op", "ValidateUserPassword", "user_id", userID.String())
			return false, ErrNotFound
		}
		r.logger.Error("failed to query user password", "op", "ValidateUserPassword", "user_id", userID.String(), "error", err.Error())
		return false, err
	}

	hasher := apputils.NewPasswordHasher()
	ok, err := hasher.Validate(password, string(passwordHash))
	if err != nil {
		r.logger.Error("failed to validate password hash", "op", "ValidateUserPassword", "user_id", userID.String(), "error", err.Error())
		return false, err
	}
	return ok, nil
}
