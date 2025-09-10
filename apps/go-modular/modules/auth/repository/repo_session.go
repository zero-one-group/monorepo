package repository

import (
	"context"
	"net"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"go-modular/modules/auth/models"
)

// CreateSession inserts a new session into the database.
func (r *AuthRepository) CreateSession(ctx context.Context, session *models.Session) error {
	if session.ID == uuid.Nil {
		session.ID = uuid.Must(uuid.NewV7())
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}
	query := `INSERT INTO ` + models.SessionTable + `
        (id, user_id, token_hash, user_agent, device_name, device_fingerprint, ip_address, expires_at, created_at, refreshed_at, revoked_at, revoked_by)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.pgPool.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.TokenHash,
		session.UserAgent,
		session.DeviceName,
		session.DeviceFingerprint,
		session.IPAddress,
		session.ExpiresAt,
		session.CreatedAt,
		session.RefreshedAt,
		session.RevokedAt,
		session.RevokedBy,
	)
	if err != nil {
		r.logger.Error("failed to insert session", "op", "CreateSession", "user_id", session.UserID.String(), "error", err.Error())
		return err
	}
	r.logger.Info("session created", "op", "CreateSession", "session_id", session.ID.String())
	return nil
}

// GetSession retrieves a session by its ID.
func (r *AuthRepository) GetSession(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	var s models.Session
	query := `SELECT id, user_id, token_hash, user_agent, device_name, device_fingerprint, ip_address, expires_at, created_at, refreshed_at, revoked_at, revoked_by
        FROM ` + models.SessionTable + ` WHERE id = $1`
	var ip net.IP
	err := r.pgPool.QueryRow(ctx, query, sessionID).Scan(
		&s.ID,
		&s.UserID,
		&s.TokenHash,
		&s.UserAgent,
		&s.DeviceName,
		&s.DeviceFingerprint,
		&ip,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.RefreshedAt,
		&s.RevokedAt,
		&s.RevokedBy,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("session not found", "op", "GetSession", "session_id", sessionID.String())
			return nil, nil
		}
		r.logger.Error("failed to get session", "op", "GetSession", "session_id", sessionID.String(), "error", err.Error())
		return nil, err
	}
	if ip != nil {
		s.IPAddress = &ip
	}
	return &s, nil
}

// UpdateSession updates an existing session in the database.
func (r *AuthRepository) UpdateSession(ctx context.Context, session *models.Session) error {
	query := `UPDATE ` + models.SessionTable + `
        SET user_id=$2, token_hash=$3, user_agent=$4, device_name=$5, device_fingerprint=$6, ip_address=$7, expires_at=$8, created_at=$9, refreshed_at=$10, revoked_at=$11, revoked_by=$12
        WHERE id=$1`
	cmd, err := r.pgPool.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.TokenHash,
		session.UserAgent,
		session.DeviceName,
		session.DeviceFingerprint,
		session.IPAddress,
		session.ExpiresAt,
		session.CreatedAt,
		session.RefreshedAt,
		session.RevokedAt,
		session.RevokedBy,
	)
	if err != nil {
		r.logger.Error("failed to update session", "op", "UpdateSession", "session_id", session.ID.String(), "error", err.Error())
		return err
	}
	if cmd.RowsAffected() == 0 {
		r.logger.Warn("session not found for update", "op", "UpdateSession", "session_id", session.ID.String())
		return pgx.ErrNoRows
	}
	r.logger.Info("session updated", "op", "UpdateSession", "session_id", session.ID.String())
	return nil
}

// DeleteSession deletes a session by its ID.
func (r *AuthRepository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	query := `DELETE FROM ` + models.SessionTable + ` WHERE id = $1`
	cmd, err := r.pgPool.Exec(ctx, query, sessionID)
	if err != nil {
		r.logger.Error("failed to delete session", "op", "DeleteSession", "session_id", sessionID.String(), "error", err.Error())
		return err
	}
	if cmd.RowsAffected() == 0 {
		r.logger.Warn("session not found for delete", "op", "DeleteSession", "session_id", sessionID.String())
		return pgx.ErrNoRows
	}
	r.logger.Info("session deleted", "op", "DeleteSession", "session_id", sessionID.String())
	return nil
}

// ValidateSession checks if a session is valid (not revoked and not expired).
func (r *AuthRepository) ValidateSession(ctx context.Context, sessionID uuid.UUID) (bool, error) {
	var revokedAt *time.Time
	var expiresAt time.Time
	query := `SELECT revoked_at, expires_at FROM ` + models.SessionTable + ` WHERE id = $1`
	err := r.pgPool.QueryRow(ctx, query, sessionID).Scan(&revokedAt, &expiresAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("session not found", "op", "ValidateSession", "session_id", sessionID.String())
			return false, nil
		}
		r.logger.Error("failed to validate session", "op", "ValidateSession", "session_id", sessionID.String(), "error", err.Error())
		return false, err
	}
	if revokedAt != nil || time.Now().After(expiresAt) {
		return false, nil
	}
	return true, nil
}
