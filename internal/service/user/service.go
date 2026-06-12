package userservice

import (
	"github.com/dionis-the-shark/apms-task-tracker/internal/models/user"
	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(u *user.User) error
	GetUserByID(userID uuid.UUID) (user.User, error)
	GetUsers() ([]user.User, error)
	UpdateUser(u user.User) error
	DeleteUser(userID uuid.UUID) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(u *user.User) error {
	return s.repo.CreateUser(u)
}

func (s *Service) GetByID(userID uuid.UUID) (user.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *Service) GetAll() ([]user.User, error) {
	return s.repo.GetUsers()
}

func (s *Service) Update(u user.User) error {
	return s.repo.UpdateUser(u)
}

func (s *Service) Delete(userID uuid.UUID) error {
	return s.repo.DeleteUser(userID)
}
