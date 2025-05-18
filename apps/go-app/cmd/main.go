package main

import (
	"log"
	"os"

	"go-app/config"
	"go-app/cmd/cli"
	"go-app/database"
)

func init() {
	config.LoadEnv()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [migrate|seed] [args...]")
	}

	command := os.Args[1]
	subcommand := ""
	if len(os.Args) >= 3 {
		subcommand = os.Args[2]
	}

	db, err := database.SetupSQLDatabase()
	if err != nil {
		log.Fatal("Failed to set up database: " + err.Error())
	}
	defer db.Close()

	switch command {
	case "migrate":
		dir := "./migrations"
		if err := cli.Migrate(db, dir, subcommand); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "seed":
		target := "all"
		if subcommand != "" {
			target = subcommand
		}
		if err := cli.Seed(db, target); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
	default:
		log.Fatalf("Unknown command: %s", command)
	}

	log.Printf("Command '%s %s' completed successfully", command, subcommand)
}
