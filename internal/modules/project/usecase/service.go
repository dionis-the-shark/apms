package usecase

import (
	"context"
	"time"

	projectmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/project"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/project/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo    ports.Repository
	now     func() time.Time
	newUUID func() uuid.UUID
}

type API interface {
	Create(ctx context.Context, input CreateInput) (projectmodule.Project, error)
	GetByID(ctx context.Context, projectID uuid.UUID) (projectmodule.Project, error)
	GetAll(ctx context.Context) ([]projectmodule.Project, error)
	Update(ctx context.Context, projectID uuid.UUID, input UpdateInput) (projectmodule.Project, error)
	Delete(ctx context.Context, projectID uuid.UUID) error
}

type CreateInput struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	ManagerID   *uuid.UUID `json:"manager_id"`
}

type UpdateInput struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	ManagerID   *uuid.UUID `json:"manager_id"`
}

func New(repo ports.Repository) *Service {
	return &Service{
		repo:    repo,
		now:     func() time.Time { return time.Now().UTC() },
		newUUID: uuid.New,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (projectmodule.Project, error) {
	p := projectmodule.Project{
		ProjectID:   s.newUUID(),
		Title:       input.Title,
		Description: input.Description,
		ManagerID:   input.ManagerID,
		CreatedAt:   s.now(),
	}
	if err := s.repo.CreateProject(ctx, &p); err != nil {
		return projectmodule.Project{}, err
	}
	return p, nil
}

func (s *Service) GetByID(ctx context.Context, projectID uuid.UUID) (projectmodule.Project, error) {
	return s.repo.GetProjectByID(ctx, projectID)
}

func (s *Service) GetAll(ctx context.Context) ([]projectmodule.Project, error) {
	return s.repo.GetProjects(ctx)
}

func (s *Service) Update(ctx context.Context, projectID uuid.UUID, input UpdateInput) (projectmodule.Project, error) {
	p := projectmodule.Project{
		ProjectID:   projectID,
		Title:       input.Title,
		Description: input.Description,
		ManagerID:   input.ManagerID,
	}
	if err := s.repo.UpdateProject(ctx, p); err != nil {
		return projectmodule.Project{}, err
	}
	return p, nil
}

func (s *Service) Delete(ctx context.Context, projectID uuid.UUID) error {
	return s.repo.DeleteProject(ctx, projectID)
}
