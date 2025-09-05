package main

import (
	"flag"
	"log"
	"os"

	"github.com/zero-one-group/go-modulith/internal/config"
	"github.com/zero-one-group/go-modulith/internal/database"
	"github.com/zero-one-group/go-modulith/internal/migration"
)

func main() {
	var (
		up     = flag.Bool("up", false, "Run migrations up")
		down   = flag.Bool("down", false, "Run migrations down")
		status = flag.Bool("status", false, "Show migration status")
		reset  = flag.Bool("reset", false, "Reset all migrations")
	)
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	migrator := migration.NewMigrator(db, cfg)

	switch {
	case *up:
		if err := migrator.Up(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations completed successfully")
	case *down:
		if err := migrator.Down(); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		log.Println("Migration rolled back successfully")
	case *status:
		if err := migrator.Status(); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	case *reset:
		if err := migrator.Reset(); err != nil {
			log.Fatalf("Failed to reset migrations: %v", err)
		}
		log.Println("Migrations reset successfully")
	default:
		flag.Usage()
		os.Exit(1)
	}
}