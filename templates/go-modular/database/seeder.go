package database

import (
	"context"
	"fmt"
	"log/slog"

	"{{ package_name | kebab_case }}/database/seeders"
)

func (m *Migrator) SeedInitialData(ctx context.Context) error {
	// Call UserFactory to seed default users
	if err := seeders.UserFactory(ctx, m.pool); err != nil {
		return fmt.Errorf("failed to seed default users: %w", err)
	}

	slog.Info("Initial data seeded successfully")

	return nil
}
