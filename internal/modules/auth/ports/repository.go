package ports

import (
	"context"

	usermodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/user"
	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, u *usermodule.User) error
	CreateUserWithSkill(ctx context.Context, u *usermodule.User, skillID uuid.UUID) error
	GetUserByEmail(ctx context.Context, email string) (usermodule.User, error)
}
