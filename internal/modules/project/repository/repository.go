package repository

import (
	"context"
	"database/sql"

	projectmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/project"
	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

type projectScanner interface {
	Scan(dest ...any) error
}

func scanProject(s projectScanner) (projectmodule.Project, error) {
	var p projectmodule.Project
	err := s.Scan(
		&p.ProjectID,
		&p.Title,
		&p.Description,
		&p.ManagerID,
		&p.CreatedAt,
	)
	return p, err
}

func (r *Repository) CreateProject(ctx context.Context, p *projectmodule.Project) error {
	if p.ProjectID == uuid.Nil {
		p.ProjectID = uuid.New()
	}

	query := `INSERT INTO projects (project_id, title, description, manager_id, created_at)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING created_at`
	return r.DB.QueryRowContext(
		ctx,
		query,
		p.ProjectID,
		p.Title,
		p.Description,
		p.ManagerID,
		p.CreatedAt,
	).Scan(&p.CreatedAt)
}

func (r *Repository) GetProjectByID(ctx context.Context, projectID uuid.UUID) (projectmodule.Project, error) {
	query := `SELECT project_id, title, description, manager_id, created_at
			  FROM projects WHERE project_id = $1`
	return scanProject(r.DB.QueryRowContext(ctx, query, projectID))
}

func (r *Repository) GetProjects(ctx context.Context) ([]projectmodule.Project, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT project_id, title, description, manager_id, created_at FROM projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []projectmodule.Project
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *Repository) UpdateProject(ctx context.Context, p projectmodule.Project) error {
	query := `UPDATE projects SET title = $1, description = $2, manager_id = $3
			  WHERE project_id = $4`
	result, err := r.DB.ExecContext(ctx, query, p.Title, p.Description, p.ManagerID, p.ProjectID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	result, err := r.DB.ExecContext(ctx, `DELETE FROM projects WHERE project_id = $1`, projectID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
