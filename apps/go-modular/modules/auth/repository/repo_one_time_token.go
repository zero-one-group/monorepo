package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"go-modular/modules/auth/models"
)

// FindAllOneTimeTokens returns all one time tokens in the database.
func (r *AuthRepository) FindAllOneTimeTokens(ctx context.Context) ([]*models.OneTimeToken, error) {
	query := `SELECT id, user_id, subject, token_hash, relates_to, metadata, created_at, expires_at, last_sent_at FROM ` + models.OneTimeTokenTable
	rows, err := r.pgPool.Query(ctx, query)
	if err != nil {
		r.logger.Error("failed to query all one time tokens", "op", "FindAllOneTimeTokens", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var tokens []*models.OneTimeToken
	for rows.Next() {
		var t models.OneTimeToken
		var metaBytes []byte
		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Subject,
			&t.TokenHash,
			&t.RelatesTo,
			&metaBytes,
			&t.CreatedAt,
			&t.ExpiresAt,
			&t.LastSentAt,
		)
		if err != nil {
			r.logger.Error("failed to scan one time token", "op", "FindAllOneTimeTokens", "error", err.Error())
			return nil, err
		}
		if len(metaBytes) > 0 {
			var m map[string]any
			if err := json.Unmarshal(metaBytes, &m); err != nil {
				r.logger.Warn("failed to unmarshal metadata for one time token", "token_id", t.ID.String(), "err", err.Error())
			} else {
				t.Metadata = m
			}
		}
		tokens = append(tokens, &t)
	}
	return tokens, nil
}

// CreateOneTimeToken inserts a new one time token into the database.
func (r *AuthRepository) CreateOneTimeToken(ctx context.Context, token *models.OneTimeToken) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.Must(uuid.NewV7())
	}
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now()
	}

	var metaArg interface{}
	if token.Metadata != nil {
		b, err := json.Marshal(token.Metadata)
		if err != nil {
			r.logger.Error("failed to marshal metadata for one time token", "op", "CreateOneTimeToken", "token_id", token.ID.String(), "error", err.Error())
			return err
		}
		metaArg = b
	} else {
		metaArg = nil
	}

	query := `INSERT INTO ` + models.OneTimeTokenTable + `
        (id, user_id, subject, token_hash, relates_to, metadata, created_at, expires_at, last_sent_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pgPool.Exec(ctx, query,
		token.ID,
		token.UserID,
		token.Subject,
		token.TokenHash,
		token.RelatesTo,
		metaArg,
		token.CreatedAt,
		token.ExpiresAt,
		token.LastSentAt,
	)
	if err != nil {
		r.logger.Error("failed to insert one time token", "op", "CreateOneTimeToken", "token_id", token.ID.String(), "error", err.Error())
		return err
	}
	r.logger.Info("one time token created", "op", "CreateOneTimeToken", "token_id", token.ID.String())
	return nil
}

// GetOneTimeTokenByID retrieves a one time token by its ID.
func (r *AuthRepository) GetOneTimeTokenByID(ctx context.Context, tokenID uuid.UUID) (*models.OneTimeToken, error) {
	query := `SELECT id, user_id, subject, token_hash, relates_to, metadata, created_at, expires_at, last_sent_at FROM ` + models.OneTimeTokenTable + ` WHERE id = $1`
	var t models.OneTimeToken
	var metaBytes []byte
	err := r.pgPool.QueryRow(ctx, query, tokenID).Scan(
		&t.ID,
		&t.UserID,
		&t.Subject,
		&t.TokenHash,
		&t.RelatesTo,
		&metaBytes,
		&t.CreatedAt,
		&t.ExpiresAt,
		&t.LastSentAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("one time token not found", "op", "GetOneTimeTokenByID", "token_id", tokenID.String())
			return nil, ErrNotFound
		}
		r.logger.Error("failed to get one time token", "op", "GetOneTimeTokenByID", "token_id", tokenID.String(), "error", err.Error())
		return nil, err
	}
	if len(metaBytes) > 0 {
		var m map[string]any
		if err := json.Unmarshal(metaBytes, &m); err != nil {
			r.logger.Warn("failed to unmarshal metadata for one time token", "token_id", t.ID.String(), "err", err.Error())
		} else {
			t.Metadata = m
		}
	}
	return &t, nil
}

// GetOneTimeTokenByTokenHash retrieves a one time token by its token_hash.
func (r *AuthRepository) GetOneTimeTokenByTokenHash(ctx context.Context, tokenHash string) (*models.OneTimeToken, error) {
	query := `SELECT id, user_id, subject, token_hash, relates_to, metadata, created_at, expires_at, last_sent_at FROM ` + models.OneTimeTokenTable + ` WHERE token_hash = $1`
	var t models.OneTimeToken
	var metaBytes []byte
	err := r.pgPool.QueryRow(ctx, query, tokenHash).Scan(
		&t.ID,
		&t.UserID,
		&t.Subject,
		&t.TokenHash,
		&t.RelatesTo,
		&metaBytes,
		&t.CreatedAt,
		&t.ExpiresAt,
		&t.LastSentAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("one time token not found", "op", "GetOneTimeTokenByTokenHash", "token_hash", tokenHash)
			return nil, ErrNotFound
		}
		r.logger.Error("failed to get one time token", "op", "GetOneTimeTokenByTokenHash", "token_hash", tokenHash, "error", err.Error())
		return nil, err
	}
	if len(metaBytes) > 0 {
		var m map[string]any
		if err := json.Unmarshal(metaBytes, &m); err != nil {
			r.logger.Warn("failed to unmarshal metadata for one time token", "token_id", t.ID.String(), "err", err.Error())
		} else {
			t.Metadata = m
		}
	}
	return &t, nil
}

// DeleteOneTimeToken deletes a one time token by its ID.
func (r *AuthRepository) DeleteOneTimeToken(ctx context.Context, tokenID uuid.UUID) error {
	query := `DELETE FROM ` + models.OneTimeTokenTable + ` WHERE id = $1`
	cmd, err := r.pgPool.Exec(ctx, query, tokenID)
	if err != nil {
		r.logger.Error("failed to delete one time token", "op", "DeleteOneTimeToken", "token_id", tokenID.String(), "error", err.Error())
		return err
	}
	if cmd.RowsAffected() == 0 {
		r.logger.Warn("one time token not found for delete", "op", "DeleteOneTimeToken", "token_id", tokenID.String())
		return ErrNotFound
	}
	r.logger.Info("one time token deleted", "op", "DeleteOneTimeToken", "token_id", tokenID.String())
	return nil
}

// UpdateOneTimeTokenLastSentAt updates the last_sent_at field of a one time token by its ID.
func (r *AuthRepository) UpdateOneTimeTokenLastSentAt(ctx context.Context, tokenID uuid.UUID, lastSentAt time.Time) error {
	query := `UPDATE ` + models.OneTimeTokenTable + ` SET last_sent_at = $1 WHERE id = $2`
	cmd, err := r.pgPool.Exec(ctx, query, lastSentAt, tokenID)
	if err != nil {
		r.logger.Error("failed to update last_sent_at", "op", "UpdateOneTimeTokenLastSentAt", "token_id", tokenID.String(), "error", err.Error())
		return err
	}
	if cmd.RowsAffected() == 0 {
		r.logger.Warn("one time token not found for update", "op", "UpdateOneTimeTokenLastSentAt", "token_id", tokenID.String())
		return ErrNotFound
	}
	r.logger.Info("one time token last_sent_at updated", "op", "UpdateOneTimeTokenLastSentAt", "token_id", tokenID.String())
	return nil
}
