package postgres

import (
	"context"
	"strings"
	"go-app/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		Conn: conn,
	}
}


func (u *UserRepository) GetUserList(ctx context.Context, filter *domain.UserFilter) ([]domain.User, error) {
	query := `
		SELECT
			u.id,
			u.name,
			u.email,
            u.created_at,
            u.updated_at
		FROM users u
        WHERE u.deleted_at is NULL`

    var args []interface{}
	var conditions []string
    if filter != nil && filter.Search != "" {
        conditions = append(conditions, `(u.name ILIKE $1 OR u.email ILIKE $1)`)
        args = append(args, "%"+filter.Search+"%")
    }

    if len(conditions) > 0 {
		query += strings.Join(conditions, " AND ")
	}
	rows, err := u.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

    var users []domain.User

	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

    return users, nil
}
