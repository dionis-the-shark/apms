package usecase

import (
	"context"
	"time"

	taskmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/task/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo    ports.Repository
	now     func() time.Time
	newUUID func() uuid.UUID
}

type API interface {
	Create(ctx context.Context, input CreateInput) (taskmodule.Task, error)
	GetByID(ctx context.Context, taskID uuid.UUID) (taskmodule.Task, error)
	GetByProject(ctx context.Context, projectID uuid.UUID) ([]taskmodule.Task, error)
	Update(ctx context.Context, taskID uuid.UUID, input UpdateInput) (taskmodule.Task, error)
	Delete(ctx context.Context, taskID uuid.UUID) error
}

type CreateInput struct {
	ProjectID       uuid.UUID  `json:"project_id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description"`
	Status          string     `json:"status"`
	Priority        int        `json:"priority"`
	EstimatedHours  int        `json:"estimated_hours"`
	RequiredSkillID *uuid.UUID `json:"required_skill_id"`
	ExecutorID      *uuid.UUID `json:"executor_id"`
}

type UpdateInput struct {
	ProjectID       uuid.UUID  `json:"project_id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description"`
	Status          string     `json:"status"`
	Priority        int        `json:"priority"`
	EstimatedHours  int        `json:"estimated_hours"`
	RequiredSkillID *uuid.UUID `json:"required_skill_id"`
	ExecutorID      *uuid.UUID `json:"executor_id"`
}

func New(repo ports.Repository) *Service {
	return &Service{
		repo:    repo,
		now:     func() time.Time { return time.Now().UTC() },
		newUUID: uuid.New,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (taskmodule.Task, error) {
	status := input.Status
	if status == "" {
		status = "backlog"
	}
	priority := input.Priority
	if priority == 0 {
		priority = 2
	}

	t := taskmodule.Task{
		TaskID:          s.newUUID(),
		ProjectID:       input.ProjectID,
		Title:           input.Title,
		Description:     input.Description,
		Status:          status,
		Priority:        priority,
		EstimatedHours:  input.EstimatedHours,
		RequiredSkillID: input.RequiredSkillID,
		ExecutorID:      input.ExecutorID,
		CreatedAt:       s.now(),
	}
	if err := s.repo.CreateTask(ctx, &t); err != nil {
		return taskmodule.Task{}, err
	}
	return t, nil
}

func (s *Service) GetByID(ctx context.Context, taskID uuid.UUID) (taskmodule.Task, error) {
	return s.repo.GetTaskByID(ctx, taskID)
}

func (s *Service) GetByProject(ctx context.Context, projectID uuid.UUID) ([]taskmodule.Task, error) {
	return s.repo.GetTasksByProject(ctx, projectID)
}

func (s *Service) Update(ctx context.Context, taskID uuid.UUID, input UpdateInput) (taskmodule.Task, error) {
	status := input.Status
	if status == "" {
		status = "backlog"
	}
	priority := input.Priority
	if priority == 0 {
		priority = 2
	}

	t := taskmodule.Task{
		TaskID:          taskID,
		ProjectID:       input.ProjectID,
		Title:           input.Title,
		Description:     input.Description,
		Status:          status,
		Priority:        priority,
		EstimatedHours:  input.EstimatedHours,
		RequiredSkillID: input.RequiredSkillID,
		ExecutorID:      input.ExecutorID,
	}
	if err := s.repo.UpdateTask(ctx, t); err != nil {
		return taskmodule.Task{}, err
	}
	return t, nil
}

func (s *Service) Delete(ctx context.Context, taskID uuid.UUID) error {
	return s.repo.DeleteTask(ctx, taskID)
}
