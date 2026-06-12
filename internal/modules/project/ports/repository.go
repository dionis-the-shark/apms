package ports

import (
	projectmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/project"
	"github.com/google/uuid"
)

type Repository interface {
	CreateProject(p *projectmodule.Project) error
	GetProjectByID(projectID uuid.UUID) (projectmodule.Project, error)
	GetProjects() ([]projectmodule.Project, error)
	UpdateProject(p projectmodule.Project) error
	DeleteProject(projectID uuid.UUID) error
}
