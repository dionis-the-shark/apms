package usecase

import (
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
	Create(input CreateInput) (projectmodule.Project, error)
	GetByID(projectID uuid.UUID) (projectmodule.Project, error)
	GetAll() ([]projectmodule.Project, error)
	Update(projectID uuid.UUID, input UpdateInput) (projectmodule.Project, error)
	Delete(projectID uuid.UUID) error
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

func (s *Service) Create(input CreateInput) (projectmodule.Project, error) {
	p := projectmodule.Project{
		ProjectID:   s.newUUID(),
		Title:       input.Title,
		Description: input.Description,
		ManagerID:   input.ManagerID,
		CreatedAt:   s.now(),
	}
	if err := s.repo.CreateProject(&p); err != nil {
		return projectmodule.Project{}, err
	}
	return p, nil
}

func (s *Service) GetByID(projectID uuid.UUID) (projectmodule.Project, error) {
	return s.repo.GetProjectByID(projectID)
}

func (s *Service) GetAll() ([]projectmodule.Project, error) {
	return s.repo.GetProjects()
}

func (s *Service) Update(projectID uuid.UUID, input UpdateInput) (projectmodule.Project, error) {
	p := projectmodule.Project{
		ProjectID:   projectID,
		Title:       input.Title,
		Description: input.Description,
		ManagerID:   input.ManagerID,
	}
	if err := s.repo.UpdateProject(p); err != nil {
		return projectmodule.Project{}, err
	}
	return p, nil
}

func (s *Service) Delete(projectID uuid.UUID) error {
	return s.repo.DeleteProject(projectID)
}
