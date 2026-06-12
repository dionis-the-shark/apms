package taskdependencymodule

import "github.com/google/uuid"

type TaskDependency struct {
	BlockedTaskID uuid.UUID `json:"blocked_task_id"`
	DependentOnID uuid.UUID `json:"dependent_on_id"`
}
