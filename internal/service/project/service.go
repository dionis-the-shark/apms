package projectservice

import (
	"github.com/dionis-the-shark/apms-task-tracker/internal/models/project"
	"github.com/dionis-the-shark/apms-task-tracker/internal/models/user"
	"github.com/google/uuid"
)

type Repository interface {
	CreateProject(p *project.Project) error
	GetProjectByID(projectID uuid.UUID) (project.Project, error)
	GetProjects() ([]project.Project, error)
	UpdateProject(p project.Project) error
	DeleteProject(projectID uuid.UUID) error
}

type UserReader interface {
	GetByID(userID uuid.UUID) (user.User, error)
}

type Service struct {
	repo  Repository
	users UserReader
}

func New(repo Repository, users UserReader) *Service {
	return &Service{repo: repo, users: users}
}

func (s *Service) Create(p *project.Project) error {
	return s.repo.CreateProject(p)
}

func (s *Service) GetByID(projectID uuid.UUID) (project.Project, error) {
	return s.repo.GetProjectByID(projectID)
}

func (s *Service) GetAll() ([]project.Project, error) {
	return s.repo.GetProjects()
}

func (s *Service) Update(p project.Project) error {
	return s.repo.UpdateProject(p)
}

func (s *Service) Delete(projectID uuid.UUID) error {
	return s.repo.DeleteProject(projectID)
}
