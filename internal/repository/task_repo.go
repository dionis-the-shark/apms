package repository

import (
	"database/sql"

	"github.com/dionis-the-shark/apms-task-tracker/internal/models/task"

	"github.com/google/uuid"
)

type TaskRepository struct {
	DB *sql.DB
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(s taskScanner) (task.Task, error) {
	var t task.Task
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

func (r *TaskRepository) CreateTask(t *task.Task) error {
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
		executor_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING created_at`

	return r.DB.QueryRow(
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
	).Scan(&t.CreatedAt)
}

func (r *TaskRepository) GetTaskByID(taskID uuid.UUID) (task.Task, error) {
	query := `SELECT task_id, project_id, title, description, status, priority, estimated_hours, required_skill_id, executor_id, created_at
			  FROM tasks WHERE task_id = $1`

	row := r.DB.QueryRow(query, taskID)
	return scanTask(row)
}

func (r *TaskRepository) GetTasksByProject(projectID uuid.UUID) ([]task.Task, error) {
	query := `SELECT task_id, project_id, title, description, status, priority, estimated_hours, required_skill_id, executor_id, created_at
              FROM tasks WHERE project_id = $1`

	rows, err := r.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []task.Task
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

func (r *TaskRepository) UpdateTask(t task.Task) error {
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

	result, err := r.DB.Exec(
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

func (r *TaskRepository) DeleteTask(taskID uuid.UUID) error {
	result, err := r.DB.Exec(`DELETE FROM tasks WHERE task_id = $1`, taskID)
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
