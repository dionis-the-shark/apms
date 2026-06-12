package repository

import (
	"database/sql"

	"github.com/dionis-the-shark/apms-task-tracker/internal/models/user"
	"github.com/google/uuid"
)

type UserRepository struct {
	DB *sql.DB
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(s userScanner) (user.User, error) {
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

func (r *UserRepository) CreateUser(u *user.User) error {
	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}

	query := `INSERT INTO users (user_id, name, email, password_hash, role)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING created_at`
	return r.DB.QueryRow(
		query,
		u.UserID,
		u.Name,
		u.Email,
		u.PasswordHash,
		u.Role,
	).Scan(&u.CreatedAt)
}

func (r *UserRepository) GetUserByID(userID uuid.UUID) (user.User, error) {
	query := `SELECT user_id, name, email, password_hash, role, created_at
			  FROM users WHERE user_id = $1`
	return scanUser(r.DB.QueryRow(query, userID))
}

func (r *UserRepository) GetUsers() ([]user.User, error) {
	rows, err := r.DB.Query(`SELECT user_id, name, email, password_hash, role, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(u user.User) error {
	query := `UPDATE users SET name = $1, email = $2, password_hash = $3, role = $4
			  WHERE user_id = $5`
	result, err := r.DB.Exec(query, u.Name, u.Email, u.PasswordHash, u.Role, u.UserID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) DeleteUser(userID uuid.UUID) error {
	result, err := r.DB.Exec(`DELETE FROM users WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
