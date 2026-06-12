package repository

import (
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

func (r *Repository) CreateProject(p *projectmodule.Project) error {
	if p.ProjectID == uuid.Nil {
		p.ProjectID = uuid.New()
	}

	query := `INSERT INTO projects (project_id, title, description, manager_id, created_at)
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING created_at`
	return r.DB.QueryRow(
		query,
		p.ProjectID,
		p.Title,
		p.Description,
		p.ManagerID,
		p.CreatedAt,
	).Scan(&p.CreatedAt)
}

func (r *Repository) GetProjectByID(projectID uuid.UUID) (projectmodule.Project, error) {
	query := `SELECT project_id, title, description, manager_id, created_at
			  FROM projects WHERE project_id = $1`
	return scanProject(r.DB.QueryRow(query, projectID))
}

func (r *Repository) GetProjects() ([]projectmodule.Project, error) {
	rows, err := r.DB.Query(`SELECT project_id, title, description, manager_id, created_at FROM projects`)
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

func (r *Repository) UpdateProject(p projectmodule.Project) error {
	query := `UPDATE projects SET title = $1, description = $2, manager_id = $3
			  WHERE project_id = $4`
	result, err := r.DB.Exec(query, p.Title, p.Description, p.ManagerID, p.ProjectID)
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

func (r *Repository) DeleteProject(projectID uuid.UUID) error {
	result, err := r.DB.Exec(`DELETE FROM projects WHERE project_id = $1`, projectID)
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
