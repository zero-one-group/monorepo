package postgres

import (
	"context"
	"errors"
	"{{ package_name | kebab_case }}/domain"
	"{{ package_name | kebab_case }}/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	Conn *pgxpool.Pool
}

func NewAuthRepository(conn *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		Conn: conn,
	}
}

func (a *AuthRepository) AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error) {

	var (
		id             uuid.UUID
		name, emailDB  string
		hashedPassword string
	)

	query := `
		SELECT id, name, email, password
		FROM users
		WHERE email = $1 AND deleted_at IS NULL`

	err := a.Conn.QueryRow(ctx, query, email).Scan(&id, &name, &emailDB, &hashedPassword)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.ComparePassword(password, hashedPassword) {
		return nil, errors.New("invalid email or password")
	}

	return &domain.User{
		ID:    id.String(),
		Name:  name,
		Email: emailDB,
	}, nil
}
