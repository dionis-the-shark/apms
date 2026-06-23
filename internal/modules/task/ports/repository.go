package ports

import (
	"context"

	taskmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task"
	"github.com/google/uuid"
)

type Repository interface {
	CreateTask(ctx context.Context, t *taskmodule.Task) error
	GetTaskByID(ctx context.Context, taskID uuid.UUID) (taskmodule.Task, error)
	GetTasksByProject(ctx context.Context, projectID uuid.UUID) ([]taskmodule.Task, error)
	UpdateTask(ctx context.Context, t taskmodule.Task) error
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
}
