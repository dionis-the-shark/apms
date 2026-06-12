package taskservice

import (
	"github.com/dionis-the-shark/apms-task-tracker/internal/models/task"
	"github.com/google/uuid"
)

type Repository interface {
	CreateTask(t *task.Task) error
	GetTaskByID(taskID uuid.UUID) (task.Task, error)
	GetTasksByProject(projectID uuid.UUID) ([]task.Task, error)
	UpdateTask(t task.Task) error
	DeleteTask(taskID uuid.UUID) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(t *task.Task) error {
	return s.repo.CreateTask(t)
}

func (s *Service) GetByID(taskID uuid.UUID) (task.Task, error) {
	return s.repo.GetTaskByID(taskID)
}

func (s *Service) GetByProject(projectID uuid.UUID) ([]task.Task, error) {
	return s.repo.GetTasksByProject(projectID)
}

func (s *Service) Update(t task.Task) error {
	return s.repo.UpdateTask(t)
}

func (s *Service) Delete(taskID uuid.UUID) error {
	return s.repo.DeleteTask(taskID)
}
