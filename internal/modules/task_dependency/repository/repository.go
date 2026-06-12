package repository

import (
	"database/sql"

	taskdependencymodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task_dependency"
	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

type taskDependencyScanner interface {
	Scan(dest ...any) error
}

func scanTaskDependency(s taskDependencyScanner) (taskdependencymodule.TaskDependency, error) {
	var td taskdependencymodule.TaskDependency
	err := s.Scan(&td.BlockedTaskID, &td.DependentOnID)
	return td, err
}

func (r *Repository) CreateTaskDependency(td taskdependencymodule.TaskDependency) error {
	query := `INSERT INTO task_dependencies (blocked_task_id, dependent_on_id)
			  VALUES ($1, $2)`
	_, err := r.DB.Exec(query, td.BlockedTaskID, td.DependentOnID)
	return err
}

func (r *Repository) GetTaskDependency(blockedTaskID, dependentOnID uuid.UUID) (taskdependencymodule.TaskDependency, error) {
	query := `SELECT blocked_task_id, dependent_on_id FROM task_dependencies
			  WHERE blocked_task_id = $1 AND dependent_on_id = $2`
	return scanTaskDependency(r.DB.QueryRow(query, blockedTaskID, dependentOnID))
}

func (r *Repository) GetTaskDependencies() ([]taskdependencymodule.TaskDependency, error) {
	rows, err := r.DB.Query(`SELECT blocked_task_id, dependent_on_id FROM task_dependencies`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []taskdependencymodule.TaskDependency
	for rows.Next() {
		td, err := scanTaskDependency(rows)
		if err != nil {
			return nil, err
		}
		deps = append(deps, td)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return deps, nil
}

func (r *Repository) GetTaskDependenciesByBlockedTask(blockedTaskID uuid.UUID) ([]taskdependencymodule.TaskDependency, error) {
	rows, err := r.DB.Query(
		`SELECT blocked_task_id, dependent_on_id FROM task_dependencies WHERE blocked_task_id = $1`,
		blockedTaskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []taskdependencymodule.TaskDependency
	for rows.Next() {
		td, err := scanTaskDependency(rows)
		if err != nil {
			return nil, err
		}
		deps = append(deps, td)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return deps, nil
}

func (r *Repository) UpdateTaskDependency(oldBlockedTaskID, oldDependentOnID, newBlockedTaskID, newDependentOnID uuid.UUID) error {
	query := `UPDATE task_dependencies SET blocked_task_id = $1, dependent_on_id = $2
			  WHERE blocked_task_id = $3 AND dependent_on_id = $4`
	result, err := r.DB.Exec(query, newBlockedTaskID, newDependentOnID, oldBlockedTaskID, oldDependentOnID)
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

func (r *Repository) DeleteTaskDependency(blockedTaskID, dependentOnID uuid.UUID) error {
	result, err := r.DB.Exec(
		`DELETE FROM task_dependencies WHERE blocked_task_id = $1 AND dependent_on_id = $2`,
		blockedTaskID,
		dependentOnID,
	)
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
