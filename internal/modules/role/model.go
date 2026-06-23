package rolemodule

import "github.com/google/uuid"

type Role struct {
	RoleID uuid.UUID `json:"role_id"`
	Title  string    `json:"title"`
	Status string    `json:"status"`
}
