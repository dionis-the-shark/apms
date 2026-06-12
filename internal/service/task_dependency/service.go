package taskdependencyservice

import (
	"github.com/dionis-the-shark/apms-task-tracker/internal/models/taskdependency"
	"github.com/google/uuid"
)

type Repository interface {
	CreateTaskDependency(td taskdependency.TaskDependency) error
	GetTaskDependency(blockedTaskID, dependentOnID uuid.UUID) (taskdependency.TaskDependency, error)
	GetTaskDependencies() ([]taskdependency.TaskDependency, error)
	GetTaskDependenciesByBlockedTask(blockedTaskID uuid.UUID) ([]taskdependency.TaskDependency, error)
	UpdateTaskDependency(oldBlockedTaskID, oldDependentOnID, newBlockedTaskID, newDependentOnID uuid.UUID) error
	DeleteTaskDependency(blockedTaskID, dependentOnID uuid.UUID) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(td taskdependency.TaskDependency) error {
	return s.repo.CreateTaskDependency(td)
}

func (s *Service) GetByIDs(blockedTaskID, dependentOnID uuid.UUID) (taskdependency.TaskDependency, error) {
	return s.repo.GetTaskDependency(blockedTaskID, dependentOnID)
}

func (s *Service) GetAll() ([]taskdependency.TaskDependency, error) {
	return s.repo.GetTaskDependencies()
}

func (s *Service) GetByBlockedTask(blockedTaskID uuid.UUID) ([]taskdependency.TaskDependency, error) {
	return s.repo.GetTaskDependenciesByBlockedTask(blockedTaskID)
}

func (s *Service) Update(oldBlockedTaskID, oldDependentOnID, newBlockedTaskID, newDependentOnID uuid.UUID) error {
	return s.repo.UpdateTaskDependency(oldBlockedTaskID, oldDependentOnID, newBlockedTaskID, newDependentOnID)
}

func (s *Service) Delete(blockedTaskID, dependentOnID uuid.UUID) error {
	return s.repo.DeleteTaskDependency(blockedTaskID, dependentOnID)
}
