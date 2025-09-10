package repository

import (
	"context"
	"net"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"go-modular/modules/auth/models"
)

// CreateRefreshToken inserts a new refresh token into the database.
func (r *AuthRepository) CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	if refreshToken.ID == uuid.Nil {
		refreshToken.ID = uuid.Must(uuid.NewV7())
	}
	if refreshToken.CreatedAt.IsZero() {
		refreshToken.CreatedAt = time.Now()
	}
	query := `INSERT INTO ` + models.RefreshTokenTable + `
        (id, user_id, session_id, token_hash, ip_address, user_agent, expires_at, created_at, revoked_at, revoked_by)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.pgPool.Exec(ctx, query,
		refreshToken.ID,
		refreshToken.UserID,
		refreshToken.SessionID,
		refreshToken.TokenHash,
		refreshToken.IPAddress,
		refreshToken.UserAgent,
		refreshToken.ExpiresAt,
		refreshToken.CreatedAt,
		refreshToken.RevokedAt,
		refreshToken.RevokedBy,
	)
	if err != nil {
		r.logger.Error("failed to insert refresh token", "op", "CreateRefreshToken", "user_id", refreshToken.UserID.String(), "error", err.Error())
		return err
	}
	r.logger.Info("refresh token created", "op", "CreateRefreshToken", "refresh_token_id", refreshToken.ID.String())
	return nil
}

// GetRefreshToken retrieves a refresh token by its ID.
func (r *AuthRepository) GetRefreshToken(ctx context.Context, tokenID uuid.UUID) (*models.RefreshToken, error) {
	var t models.RefreshToken
	query := `SELECT id, user_id, session_id, token_hash, ip_address, user_agent, expires_at, created_at, revoked_at, revoked_by
        FROM ` + models.RefreshTokenTable + ` WHERE id = $1`
	var ip net.IP
	err := r.pgPool.QueryRow(ctx, query, tokenID).Scan(
		&t.ID,
		&t.UserID,
		&t.SessionID,
		&t.TokenHash,
		&ip,
		&t.UserAgent,
		&t.ExpiresAt,
		&t.CreatedAt,
		&t.RevokedAt,
		&t.RevokedBy,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("refresh token not found", "op", "GetRefreshToken", "refresh_token_id", tokenID.String())
			return nil, ErrNotFound
		}
		r.logger.Error("failed to get refresh token", "op", "GetRefreshToken", "refresh_token_id", tokenID.String(), "error", err.Error())
		return nil, err
	}
	if ip != nil {
		t.IPAddress = &ip
	}
	return &t, nil
}

// UpdateRefreshToken updates an existing refresh token in the database.
func (r *AuthRepository) UpdateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := `UPDATE ` + models.RefreshTokenTable + `
        SET user_id=$2, session_id=$3, token_hash=$4, ip_address=$5, user_agent=$6, expires_at=$7, created_at=$8, revoked_at=$9, revoked_by=$10
        WHERE id=$1`
	cmd, err := r.pgPool.Exec(ctx, query,
		refreshToken.ID,
		refreshToken.UserID,
		refreshToken.SessionID,
		refreshToken.TokenHash,
		refreshToken.IPAddress,
		refreshToken.UserAgent,
		refreshToken.ExpiresAt,
		refreshToken.CreatedAt,
		refreshToken.RevokedAt,
		refreshToken.RevokedBy,
	)
	if err != nil {
		r.logger.Error("failed to update refresh token", "op", "UpdateRefreshToken", "refresh_token_id", refreshToken.ID.String(), "error", err.Error())
		return err
	}
	if cmd.RowsAffected() == 0 {
		r.logger.Warn("refresh token not found for update", "op", "UpdateRefreshToken", "refresh_token_id", refreshToken.ID.String())
		return ErrNotFound
	}
	r.logger.Info("refresh token updated", "op", "UpdateRefreshToken", "refresh_token_id", refreshToken.ID.String())
	return nil
}

// DeleteRefreshToken deletes a refresh token by its ID.
func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	query := `DELETE FROM ` + models.RefreshTokenTable + ` WHERE id = $1`
	cmd, err := r.pgPool.Exec(ctx, query, tokenID)
	if err != nil {
		r.logger.Error("failed to delete refresh token", "op", "DeleteRefreshToken", "refresh_token_id", tokenID.String(), "error", err.Error())
		return err
	}
	if cmd.RowsAffected() == 0 {
		r.logger.Warn("refresh token not found for delete", "op", "DeleteRefreshToken", "refresh_token_id", tokenID.String())
		return ErrNotFound
	}
	r.logger.Info("refresh token deleted", "op", "DeleteRefreshToken", "refresh_token_id", tokenID.String())
	return nil
}

// ValidateRefreshToken checks if a refresh token is valid (not revoked and not expired).
func (r *AuthRepository) ValidateRefreshToken(ctx context.Context, tokenID uuid.UUID) (bool, error) {
	var revokedAt *time.Time
	var expiresAt time.Time
	query := `SELECT revoked_at, expires_at FROM ` + models.RefreshTokenTable + ` WHERE id = $1`
	err := r.pgPool.QueryRow(ctx, query, tokenID).Scan(&revokedAt, &expiresAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("refresh token not found", "op", "ValidateRefreshToken", "refresh_token_id", tokenID.String())
			return false, ErrNotFound
		}
		r.logger.Error("failed to validate refresh token", "op", "ValidateRefreshToken", "refresh_token_id", tokenID.String(), "error", err.Error())
		return false, err
	}
	if revokedAt != nil || time.Now().After(expiresAt) {
		return false, nil
	}
	return true, nil
}
