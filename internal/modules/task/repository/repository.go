package repository

import (
	"context"
	"database/sql"

	taskmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task"
	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(s taskScanner) (taskmodule.Task, error) {
	var t taskmodule.Task
	err := s.Scan(
		&t.TaskID,
		&t.ProjectID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Priority,
		&t.EstimatedHours,
		&t.RequiredSkillID,
		&t.ExecutorID,
		&t.CreatedAt,
	)
	return t, err
}

func (r *Repository) CreateTask(ctx context.Context, t *taskmodule.Task) error {
	if t.TaskID == uuid.Nil {
		t.TaskID = uuid.New()
	}

	query := `INSERT INTO tasks (
		task_id,
		project_id,
		title,
		description,
		status,
		priority,
		estimated_hours,
		required_skill_id,
		executor_id,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING created_at`

	return r.DB.QueryRowContext(
		ctx,
		query,
		t.TaskID,
		t.ProjectID,
		t.Title,
		t.Description,
		t.Status,
		t.Priority,
		t.EstimatedHours,
		t.RequiredSkillID,
		t.ExecutorID,
		t.CreatedAt,
	).Scan(&t.CreatedAt)
}

func (r *Repository) GetTaskByID(ctx context.Context, taskID uuid.UUID) (taskmodule.Task, error) {
	query := `SELECT task_id, project_id, title, description, status, priority, estimated_hours, required_skill_id, executor_id, created_at
			  FROM tasks WHERE task_id = $1`

	row := r.DB.QueryRowContext(ctx, query, taskID)
	return scanTask(row)
}

func (r *Repository) GetTasksByProject(ctx context.Context, projectID uuid.UUID) ([]taskmodule.Task, error) {
	query := `SELECT task_id, project_id, title, description, status, priority, estimated_hours, required_skill_id, executor_id, created_at
              FROM tasks WHERE project_id = $1`

	rows, err := r.DB.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []taskmodule.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) UpdateTask(ctx context.Context, t taskmodule.Task) error {
	query := `UPDATE tasks SET
		project_id = $1,
		title = $2,
		description = $3,
		status = $4,
		priority = $5,
		estimated_hours = $6,
		required_skill_id = $7,
		executor_id = $8
	WHERE task_id = $9`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		t.ProjectID,
		t.Title,
		t.Description,
		t.Status,
		t.Priority,
		t.EstimatedHours,
		t.RequiredSkillID,
		t.ExecutorID,
		t.TaskID,
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

func (r *Repository) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	result, err := r.DB.ExecContext(ctx, `DELETE FROM tasks WHERE task_id = $1`, taskID)
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
