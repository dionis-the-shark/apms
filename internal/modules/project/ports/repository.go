package ports

import (
	"context"

	projectmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/project"
	"github.com/google/uuid"
)

type Repository interface {
	CreateProject(ctx context.Context, p *projectmodule.Project) error
	GetProjectByID(ctx context.Context, projectID uuid.UUID) (projectmodule.Project, error)
	GetProjects(ctx context.Context) ([]projectmodule.Project, error)
	UpdateProject(ctx context.Context, p projectmodule.Project) error
	DeleteProject(ctx context.Context, projectID uuid.UUID) error
}
