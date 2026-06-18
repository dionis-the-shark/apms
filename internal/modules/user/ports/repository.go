package ports

import (
	"context"

	usermodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/user"
	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, u *usermodule.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (usermodule.User, error)
	GetUsers(ctx context.Context) ([]usermodule.User, error)
	UpdateUser(ctx context.Context, u usermodule.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}
