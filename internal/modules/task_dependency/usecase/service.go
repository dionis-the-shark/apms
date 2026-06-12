package usecase

import (
	taskdependencymodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task_dependency"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/task_dependency/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo ports.Repository
}

type API interface {
	Create(input CreateInput) (taskdependencymodule.TaskDependency, error)
	GetByIDs(blockedTaskID, dependentOnID uuid.UUID) (taskdependencymodule.TaskDependency, error)
	GetAll() ([]taskdependencymodule.TaskDependency, error)
	GetByBlockedTask(blockedTaskID uuid.UUID) ([]taskdependencymodule.TaskDependency, error)
	Update(oldBlockedTaskID, oldDependentOnID uuid.UUID, input UpdateInput) (taskdependencymodule.TaskDependency, error)
	Delete(blockedTaskID, dependentOnID uuid.UUID) error
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

func (s *Service) Create(input CreateInput) (taskdependencymodule.TaskDependency, error) {
	td := taskdependencymodule.TaskDependency{
		BlockedTaskID: input.BlockedTaskID,
		DependentOnID: input.DependentOnID,
	}
	if err := s.repo.CreateTaskDependency(td); err != nil {
		return taskdependencymodule.TaskDependency{}, err
	}
	return td, nil
}

func (s *Service) GetByIDs(blockedTaskID, dependentOnID uuid.UUID) (taskdependencymodule.TaskDependency, error) {
	return s.repo.GetTaskDependency(blockedTaskID, dependentOnID)
}

func (s *Service) GetAll() ([]taskdependencymodule.TaskDependency, error) {
	return s.repo.GetTaskDependencies()
}

func (s *Service) GetByBlockedTask(blockedTaskID uuid.UUID) ([]taskdependencymodule.TaskDependency, error) {
	return s.repo.GetTaskDependenciesByBlockedTask(blockedTaskID)
}

func (s *Service) Update(oldBlockedTaskID, oldDependentOnID uuid.UUID, input UpdateInput) (taskdependencymodule.TaskDependency, error) {
	if err := s.repo.UpdateTaskDependency(oldBlockedTaskID, oldDependentOnID, input.BlockedTaskID, input.DependentOnID); err != nil {
		return taskdependencymodule.TaskDependency{}, err
	}
	return taskdependencymodule.TaskDependency{
		BlockedTaskID: input.BlockedTaskID,
		DependentOnID: input.DependentOnID,
	}, nil
}

func (s *Service) Delete(blockedTaskID, dependentOnID uuid.UUID) error {
	return s.repo.DeleteTaskDependency(blockedTaskID, dependentOnID)
}
