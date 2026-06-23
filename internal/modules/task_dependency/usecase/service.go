package usecase

import (
	"context"

	taskdependencymodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task_dependency"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/task_dependency/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo ports.Repository
}

type API interface {
	Create(ctx context.Context, input CreateInput) (taskdependencymodule.TaskDependency, error)
	GetByIDs(ctx context.Context, blockedTaskID, dependentOnID uuid.UUID) (taskdependencymodule.TaskDependency, error)
	GetAll(ctx context.Context) ([]taskdependencymodule.TaskDependency, error)
	GetByBlockedTask(ctx context.Context, blockedTaskID uuid.UUID) ([]taskdependencymodule.TaskDependency, error)
	Update(ctx context.Context, oldBlockedTaskID, oldDependentOnID uuid.UUID, input UpdateInput) (taskdependencymodule.TaskDependency, error)
	Delete(ctx context.Context, blockedTaskID, dependentOnID uuid.UUID) error
}

type CreateInput struct {
	BlockedTaskID uuid.UUID `json:"blocked_task_id"`
	DependentOnID uuid.UUID `json:"dependent_on_id"`
}

type UpdateInput struct {
	BlockedTaskID uuid.UUID `json:"blocked_task_id"`
	DependentOnID uuid.UUID `json:"dependent_on_id"`
}

func New(repo ports.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (taskdependencymodule.TaskDependency, error) {
	td := taskdependencymodule.TaskDependency{
		BlockedTaskID: input.BlockedTaskID,
		DependentOnID: input.DependentOnID,
	}
	if err := s.repo.CreateTaskDependency(ctx, td); err != nil {
		return taskdependencymodule.TaskDependency{}, err
	}
	return td, nil
}

func (s *Service) GetByIDs(ctx context.Context, blockedTaskID, dependentOnID uuid.UUID) (taskdependencymodule.TaskDependency, error) {
	return s.repo.GetTaskDependency(ctx, blockedTaskID, dependentOnID)
}

func (s *Service) GetAll(ctx context.Context) ([]taskdependencymodule.TaskDependency, error) {
	return s.repo.GetTaskDependencies(ctx)
}

func (s *Service) GetByBlockedTask(ctx context.Context, blockedTaskID uuid.UUID) ([]taskdependencymodule.TaskDependency, error) {
	return s.repo.GetTaskDependenciesByBlockedTask(ctx, blockedTaskID)
}

func (s *Service) Update(ctx context.Context, oldBlockedTaskID, oldDependentOnID uuid.UUID, input UpdateInput) (taskdependencymodule.TaskDependency, error) {
	if err := s.repo.UpdateTaskDependency(ctx, oldBlockedTaskID, oldDependentOnID, input.BlockedTaskID, input.DependentOnID); err != nil {
		return taskdependencymodule.TaskDependency{}, err
	}
	return taskdependencymodule.TaskDependency{
		BlockedTaskID: input.BlockedTaskID,
		DependentOnID: input.DependentOnID,
	}, nil
}

func (s *Service) Delete(ctx context.Context, blockedTaskID, dependentOnID uuid.UUID) error {
	return s.repo.DeleteTaskDependency(ctx, blockedTaskID, dependentOnID)
}
