package seeders

import (
	"database/sql"

	"{{ package_name }}/utils"
)

func SeedUsers(db *sql.DB) error {
	// Hash passwords
	password, err := utils.HashPassword("password123")
	if err != nil {
		return err
	}

	// Insert users with password hashes
	_, err = db.Exec(`
        INSERT INTO users (name, email, password) VALUES
        ('Alice', 'alice@example.com', $1),
        ('Bob', 'bob@example.com', $1)
        ON CONFLICT DO NOTHING;
    `, password)

	return err
}
