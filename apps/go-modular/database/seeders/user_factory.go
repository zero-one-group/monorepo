package seeders

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserSeed struct {
	Username    string
	DisplayName string
	Metadata    map[string]string
}

func UserFactory(ctx context.Context, pool *pgxpool.Pool) error {
	slog.Info("Seeding default users...")

	users := []UserSeed{
		{
			Username:    "admin",
			DisplayName: "Admin Sistem",
			Metadata: map[string]string{
				"timezone": "Asia/Jakarta",
			},
		},
	}

	insertUserQuery := `
		INSERT INTO public.users (username, display_name, metadata)
		VALUES ($1, $2, $3)
		ON CONFLICT (username) DO UPDATE
		SET display_name = EXCLUDED.display_name, metadata = EXCLUDED.metadata
	`

	for _, u := range users {
		metadataJSON, err := json.Marshal(u.Metadata)
		if err != nil {
			slog.Error("Failed to marshal metadata", "username", u.Username, "err", err)
			return err
		}
		_, err = pool.Exec(ctx, insertUserQuery, u.Username, u.DisplayName, metadataJSON)
		if err != nil {
			slog.Error("Failed to seed user", "username", u.Username, "err", err)
			return err
		}
		slog.Info("Seeded user", "username", u.Username)
	}

	return nil
}
