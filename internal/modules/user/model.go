package usermodule

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	RoleID       string    `json:"role_id"`
	CreatedAt    time.Time `json:"created_at"`
}
