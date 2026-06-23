package repository

import (
	"context"
	"database/sql"
	"time"

	usermodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/user"
	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

type authUserScanner interface {
	Scan(dest ...any) error
}

func scanAuthUser(s authUserScanner) (usermodule.User, error) {
	var u usermodule.User
	err := s.Scan(
		&u.UserID,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.RoleID,
		&u.CreatedAt,
	)
	return u, err
}

func (r *Repository) CreateUser(ctx context.Context, u *usermodule.User) error {
	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}

	query := `INSERT INTO users (user_id, name, email, password_hash, role_id, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING created_at`
	return r.DB.QueryRowContext(
		ctx,
		query,
		u.UserID,
		u.Name,
		u.Email,
		u.PasswordHash,
		u.RoleID,
		u.CreatedAt,
	).Scan(&u.CreatedAt)
}

func (r *Repository) CreateUserWithSkill(ctx context.Context, u *usermodule.User, skillID uuid.UUID) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}

	userQuery := `INSERT INTO users (user_id, name, email, password_hash, role_id, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING created_at`
	if err = tx.QueryRowContext(
		ctx,
		userQuery,
		u.UserID,
		u.Name,
		u.Email,
		u.PasswordHash,
		u.RoleID,
		u.CreatedAt,
	).Scan(&u.CreatedAt); err != nil {
		_ = tx.Rollback()
		return err
	}

	skillQuery := `INSERT INTO developer_skills (developer_id, skill_id)
			   VALUES ($1, $2)`
	if _, err = tx.ExecContext(ctx, skillQuery, u.UserID, skillID); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (usermodule.User, error) {
	query := `SELECT user_id, name, email, password_hash, role_id, created_at
			  FROM users WHERE email = $1`
	return scanAuthUser(r.DB.QueryRowContext(ctx, query, email))
}
