package seeders

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"{{ package_name | kebab_case }}/pkg/apputils"
)

type UserSeed struct {
	DisplayName     string
	Email           string
	Username        string
	Password        string
	Metadata        map[string]string
	EmailVerifiedAt *time.Time // pakai time.Time pointer
}

func UserFactory(ctx context.Context, pool *pgxpool.Pool) error {
	slog.Info("Seeding default users...")

	now := time.Now().UTC()

	users := []UserSeed{
		{
			DisplayName:     "Admin Sistem",
			Email:           "admin@example.com",
			Username:        "admin",
			Password:        "secure.password",
			Metadata:        map[string]string{"timezone": "Asia/Jakarta"},
			EmailVerifiedAt: &now, // verified
		},
		{
			DisplayName:     "John Doe",
			Email:           "johndoe@example.com",
			Username:        "johndoe",
			Password:        "secure.password",
			Metadata:        map[string]string{"timezone": "UTC"},
			EmailVerifiedAt: nil, // unverified
		},
	}

	insertUserQuery := `
        INSERT INTO public.users (username, display_name, email, metadata, email_verified_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (username) DO UPDATE
        SET
            display_name = EXCLUDED.display_name,
            email = EXCLUDED.email,
            metadata = EXCLUDED.metadata,
            email_verified_at = EXCLUDED.email_verified_at
        RETURNING id
    `

	insertPasswordQuery := `
        INSERT INTO public.user_passwords (user_id, password_hash) VALUES ($1, $2)
        ON CONFLICT (user_id) DO UPDATE SET password_hash = EXCLUDED.password_hash
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

	passwd := apputils.NewPasswordHasher()

	for _, u := range users {
		metadataJSON, err := json.Marshal(u.Metadata)
		if err != nil {
			slog.Error("Failed to marshal metadata", "username", u.Username, "err", err)
			return err
		}

		var userID string
		err = tx.QueryRow(ctx, insertUserQuery, u.Username, u.DisplayName, u.Email, metadataJSON, u.EmailVerifiedAt).Scan(&userID)
		if err != nil {
			slog.Error("Failed to seed user", "username", u.Username, "err", err)
			return err
		}

		passwordHash, err := passwd.Hash(u.Password)
		if err != nil {
			slog.Error("Failed to hash password", "username", u.Username, "err", err)
			return err
		}

		_, err = tx.Exec(ctx, insertPasswordQuery, userID, passwordHash)
		if err != nil {
			slog.Error("Failed to seed user password", "username", u.Username, "err", err)
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
