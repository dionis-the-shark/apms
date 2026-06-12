package project

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ProjectID   uuid.UUID  `json:"project_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	ManagerID   *uuid.UUID `json:"manager_id"`
	CreatedAt   time.Time  `json:"created_at"`
}
