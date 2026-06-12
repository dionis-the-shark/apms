package repository

import (
	"context"
	"database/sql"

	"github.com/dionis-the-shark/apms-task-tracker/internal/models/user"
	"github.com/google/uuid"
)

type AuthRepository struct {
	DB *sql.DB
}

type authUserScanner interface {
	Scan(dest ...any) error
}

func scanAuthUser(s authUserScanner) (user.User, error) {
	var u user.User
	err := s.Scan(
		&u.UserID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.CreatedAt,
	)
	return u, err
}

func (r *AuthRepository) CreateUser(ctx context.Context, u *user.User) error {
	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}

	query := `INSERT INTO users (user_id, name, email, password_hash, role)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING created_at`
	return r.DB.QueryRowContext(
		ctx,
		query,
		u.UserID,
		u.Name,
		u.Email,
		u.PasswordHash,
		u.Role,
	).Scan(&u.CreatedAt)
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (user.User, error) {
	query := `SELECT user_id, name, email, password_hash, role, created_at
			  FROM users WHERE email = $1`
	return scanAuthUser(r.DB.QueryRowContext(ctx, query, email))
}
