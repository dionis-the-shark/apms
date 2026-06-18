package usecase

import (
	"context"
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
	Create(ctx context.Context, input CreateInput) (usermodule.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (usermodule.User, error)
	GetAll(ctx context.Context) ([]usermodule.User, error)
	Update(ctx context.Context, userID uuid.UUID, input UpdateInput) (usermodule.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
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

func (s *Service) Create(ctx context.Context, input CreateInput) (usermodule.User, error) {
	u := usermodule.User{
		UserID:       s.newUUID(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		Role:         input.Role,
		CreatedAt:    s.now(),
	}
	if err := s.repo.CreateUser(ctx, &u); err != nil {
		return usermodule.User{}, err
	}
	return u, nil
}

func (s *Service) GetByID(ctx context.Context, userID uuid.UUID) (usermodule.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) GetAll(ctx context.Context) ([]usermodule.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *Service) Update(ctx context.Context, userID uuid.UUID, input UpdateInput) (usermodule.User, error) {
	u := usermodule.User{
		UserID:       userID,
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		Role:         input.Role,
	}
	if err := s.repo.UpdateUser(ctx, u); err != nil {
		return usermodule.User{}, err
	}
	return u, nil
}

func (s *Service) Delete(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteUser(ctx, userID)
}
