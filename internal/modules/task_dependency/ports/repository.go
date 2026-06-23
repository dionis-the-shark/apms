package ports

import (
	"context"

	taskdependencymodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task_dependency"
	"github.com/google/uuid"
)

type Repository interface {
	CreateTaskDependency(ctx context.Context, td taskdependencymodule.TaskDependency) error
	GetTaskDependency(ctx context.Context, blockedTaskID, dependentOnID uuid.UUID) (taskdependencymodule.TaskDependency, error)
	GetTaskDependencies(ctx context.Context) ([]taskdependencymodule.TaskDependency, error)
	GetTaskDependenciesByBlockedTask(ctx context.Context, blockedTaskID uuid.UUID) ([]taskdependencymodule.TaskDependency, error)
	UpdateTaskDependency(ctx context.Context, oldBlockedTaskID, oldDependentOnID, newBlockedTaskID, newDependentOnID uuid.UUID) error
	DeleteTaskDependency(ctx context.Context, blockedTaskID, dependentOnID uuid.UUID) error
}
