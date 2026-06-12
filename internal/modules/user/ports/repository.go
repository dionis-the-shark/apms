package ports

import (
	usermodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/user"
	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(u *usermodule.User) error
	GetUserByID(userID uuid.UUID) (usermodule.User, error)
	GetUsers() ([]usermodule.User, error)
	UpdateUser(u usermodule.User) error
	DeleteUser(userID uuid.UUID) error
}
