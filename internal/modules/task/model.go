package taskmodule

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	TaskID          uuid.UUID  `json:"task_id"`
	ProjectID       uuid.UUID  `json:"project_id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description"`
	Status          string     `json:"status"`
	Priority        int        `json:"priority"`
	EstimatedHours  int        `json:"estimated_hours"`
	RequiredSkillID *uuid.UUID `json:"required_skill_id"`
	ExecutorID      *uuid.UUID `json:"executor_id"`
	CreatedAt       time.Time  `json:"created_at"`
}
