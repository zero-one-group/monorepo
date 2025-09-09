package seeders

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserSeed struct {
	DisplayName string
	Email       string
	Username    string
	Metadata    map[string]string
}

func UserFactory(ctx context.Context, pool *pgxpool.Pool) error {
	slog.Info("Seeding default users...")

	users := []UserSeed{
		{
			DisplayName: "Admin Sistem",
			Email:       "admin@example.com",
			Username:    "admin",
			Metadata: map[string]string{
				"timezone": "Asia/Jakarta",
			},
		},
	}

	insertUserQuery := `
		INSERT INTO public.users (username, display_name, email, metadata)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (username) DO UPDATE
		SET display_name = EXCLUDED.display_name, email = EXCLUDED.email, metadata = EXCLUDED.metadata
	`

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		slog.Error("Failed to begin transaction", "err", err)
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			slog.Warn("Failed to rollback transaction", "err", err)
		}
	}()

	for _, u := range users {
		metadataJSON, err := json.Marshal(u.Metadata)
		if err != nil {
			slog.Error("Failed to marshal metadata", "username", u.Username, "err", err)
			return err
		}
		_, err = tx.Exec(ctx, insertUserQuery, u.Username, u.DisplayName, u.Email, metadataJSON)
		if err != nil {
			slog.Error("Failed to seed user", "username", u.Username, "err", err)
			return err
		}
		slog.Info("Seeded user", "username", u.Username)
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("Failed to commit transaction", "err", err)
		return err
	}

	return nil
}
