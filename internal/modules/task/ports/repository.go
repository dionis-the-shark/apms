package ports

import (
	taskmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task"
	"github.com/google/uuid"
)

type Repository interface {
	CreateTask(t *taskmodule.Task) error
	GetTaskByID(taskID uuid.UUID) (taskmodule.Task, error)
	GetTasksByProject(projectID uuid.UUID) ([]taskmodule.Task, error)
	UpdateTask(t taskmodule.Task) error
	DeleteTask(taskID uuid.UUID) error
}
