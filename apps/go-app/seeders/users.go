package seeders

import (
	"database/sql"

	"go-app/cmd/bcrypt"
)

func SeedUsers(db *sql.DB) error {
    // Hash passwords
    password, err := bcrypt.Hash("password123")
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
