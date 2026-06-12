package ports

import (
	rolemodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/role"
	"github.com/google/uuid"
)

type Repository interface {
	CreateRole(r *rolemodule.Role) error
	GetRoleByID(roleID uuid.UUID) (rolemodule.Role, error)
	GetRoles() ([]rolemodule.Role, error)
	UpdateRole(r rolemodule.Role) error
	DeleteRole(roleID uuid.UUID) error
}
