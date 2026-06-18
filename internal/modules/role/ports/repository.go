package ports

import (
	"context"

	rolemodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/role"
	"github.com/google/uuid"
)

type Repository interface {
	CreateRole(ctx context.Context, r *rolemodule.Role) error
	GetRoleByID(ctx context.Context, roleID uuid.UUID) (rolemodule.Role, error)
	GetRoles(ctx context.Context) ([]rolemodule.Role, error)
	UpdateRole(ctx context.Context, r rolemodule.Role) error
	DeleteRole(ctx context.Context, roleID uuid.UUID) error
}
