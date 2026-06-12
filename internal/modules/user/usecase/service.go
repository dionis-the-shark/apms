package usecase

import (
	"time"

	usermodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/user"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/user/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo    ports.Repository
	now     func() time.Time
	newUUID func() uuid.UUID
}

type API interface {
	Create(input CreateInput) (usermodule.User, error)
	GetByID(userID uuid.UUID) (usermodule.User, error)
	GetAll() ([]usermodule.User, error)
	Update(userID uuid.UUID, input UpdateInput) (usermodule.User, error)
	Delete(userID uuid.UUID) error
}

type CreateInput struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

type UpdateInput struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
}

func New(repo ports.Repository) *Service {
	return &Service{
		repo:    repo,
		now:     func() time.Time { return time.Now().UTC() },
		newUUID: uuid.New,
	}
}

func (s *Service) Create(input CreateInput) (usermodule.User, error) {
	u := usermodule.User{
		UserID:       s.newUUID(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		Role:         input.Role,
		CreatedAt:    s.now(),
	}
	if err := s.repo.CreateUser(&u); err != nil {
		return usermodule.User{}, err
	}
	return u, nil
}

func (s *Service) GetByID(userID uuid.UUID) (usermodule.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *Service) GetAll() ([]usermodule.User, error) {
	return s.repo.GetUsers()
}

func (s *Service) Update(userID uuid.UUID, input UpdateInput) (usermodule.User, error) {
	u := usermodule.User{
		UserID:       userID,
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		Role:         input.Role,
	}
	if err := s.repo.UpdateUser(u); err != nil {
		return usermodule.User{}, err
	}
	return u, nil
}

func (s *Service) Delete(userID uuid.UUID) error {
	return s.repo.DeleteUser(userID)
}
