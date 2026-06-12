package repository

import (
	"database/sql"

	usermodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/user"
	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(s userScanner) (usermodule.User, error) {
	var u usermodule.User
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

func (r *Repository) CreateUser(u *usermodule.User) error {
	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}

	query := `INSERT INTO users (user_id, name, email, password_hash, role, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING created_at`
	return r.DB.QueryRow(
		query,
		u.UserID,
		u.Name,
		u.Email,
		u.PasswordHash,
		u.Role,
		u.CreatedAt,
	).Scan(&u.CreatedAt)
}

func (r *Repository) GetUserByID(userID uuid.UUID) (usermodule.User, error) {
	query := `SELECT user_id, name, email, password_hash, role, created_at
			  FROM users WHERE user_id = $1`
	return scanUser(r.DB.QueryRow(query, userID))
}

func (r *Repository) GetUsers() ([]usermodule.User, error) {
	rows, err := r.DB.Query(`SELECT user_id, name, email, password_hash, role, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []usermodule.User
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

func (r *Repository) UpdateUser(u usermodule.User) error {
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

func (r *Repository) DeleteUser(userID uuid.UUID) error {
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
