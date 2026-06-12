package ports

import (
	taskdependencymodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task_dependency"
	"github.com/google/uuid"
)

type Repository interface {
	CreateTaskDependency(td taskdependencymodule.TaskDependency) error
	GetTaskDependency(blockedTaskID, dependentOnID uuid.UUID) (taskdependencymodule.TaskDependency, error)
	GetTaskDependencies() ([]taskdependencymodule.TaskDependency, error)
	GetTaskDependenciesByBlockedTask(blockedTaskID uuid.UUID) ([]taskdependencymodule.TaskDependency, error)
	UpdateTaskDependency(oldBlockedTaskID, oldDependentOnID, newBlockedTaskID, newDependentOnID uuid.UUID) error
	DeleteTaskDependency(blockedTaskID, dependentOnID uuid.UUID) error
}
