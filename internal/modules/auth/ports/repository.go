package ports

import (
	"context"

	usermodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/user"
)

type Repository interface {
	CreateUser(ctx context.Context, u *usermodule.User) error
	GetUserByEmail(ctx context.Context, email string) (usermodule.User, error)
}
