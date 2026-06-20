package repository

import (
	"context"
	"database/sql"

	rolemodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/role"
	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

type roleScanner interface {
	Scan(dest ...any) error
}

func scanRole(s roleScanner) (rolemodule.Role, error) {
	var r rolemodule.Role
	err := s.Scan(&r.RoleID, &r.Title, &r.Status)
	return r, err
}

func (r *Repository) CreateRole(ctx context.Context, roleModel *rolemodule.Role) error {
	if roleModel.RoleID == uuid.Nil {
		roleModel.RoleID = uuid.New()
	}

	query := `INSERT INTO roles (role_id, title, status)
			  VALUES ($1, $2, $3)`
	_, err := r.DB.ExecContext(ctx, query, roleModel.RoleID, roleModel.Title, roleModel.Status)
	return err
}

func (r *Repository) GetRoleByID(ctx context.Context, roleID uuid.UUID) (rolemodule.Role, error) {
	query := `SELECT role_id, title, status FROM roles WHERE role_id = $1`
	return scanRole(r.DB.QueryRowContext(ctx, query, roleID))
}

func (r *Repository) GetRoles(ctx context.Context) ([]rolemodule.Role, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT role_id, title, status FROM roles WHERE status = 'active'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []rolemodule.Role
	for rows.Next() {
		roleModel, err := scanRole(rows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, roleModel)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *Repository) UpdateRole(ctx context.Context, roleModel rolemodule.Role) error {
	query := `UPDATE roles SET title = $1, status = $2 WHERE role_id = $3`
	result, err := r.DB.ExecContext(ctx, query, roleModel.Title, roleModel.Status, roleModel.RoleID)
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

func (r *Repository) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	result, err := r.DB.ExecContext(ctx, `DELETE FROM roles WHERE role_id = $1`, roleID)
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
