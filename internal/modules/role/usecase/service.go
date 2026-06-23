package usecase

import (
	"context"

	"github.com/dionis-the-shark/apms-task-tracker/internal/authz"
	rolemodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/role"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/role/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo    ports.Repository
	newUUID func() uuid.UUID
}

type API interface {
	Create(ctx context.Context, input CreateInput) (rolemodule.Role, error)
	GetByID(ctx context.Context, roleID uuid.UUID) (rolemodule.Role, error)
	GetAll(ctx context.Context) ([]rolemodule.Role, error)
	Update(ctx context.Context, roleID uuid.UUID, input UpdateInput) (rolemodule.Role, error)
	Delete(ctx context.Context, roleID uuid.UUID) error
	ActionAllowed(ctx context.Context, roleID string, action authz.Action) (bool, error)
}

type CreateInput struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

type UpdateInput struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

func New(repo ports.Repository) *Service {
	return &Service{
		repo:    repo,
		newUUID: uuid.New,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (rolemodule.Role, error) {
	status := input.Status
	if status == "" {
		status = "active"
	}

	r := rolemodule.Role{
		RoleID: s.newUUID(),
		Title:  input.Title,
		Status: status,
	}
	if err := s.repo.CreateRole(ctx, &r); err != nil {
		return rolemodule.Role{}, err
	}
	return r, nil
}

func (s *Service) GetByID(ctx context.Context, roleID uuid.UUID) (rolemodule.Role, error) {
	return s.repo.GetRoleByID(ctx, roleID)
}

func (s *Service) GetAll(ctx context.Context) ([]rolemodule.Role, error) {
	return s.repo.GetRoles(ctx)
}

func (s *Service) Update(ctx context.Context, roleID uuid.UUID, input UpdateInput) (rolemodule.Role, error) {
	status := input.Status
	if status == "" {
		status = "active"
	}

	r := rolemodule.Role{
		RoleID: roleID,
		Title:  input.Title,
		Status: status,
	}
	if err := s.repo.UpdateRole(ctx, r); err != nil {
		return rolemodule.Role{}, err
	}
	return r, nil
}

func (s *Service) Delete(ctx context.Context, roleID uuid.UUID) error {
	return s.repo.DeleteRole(ctx, roleID)
}

func (s *Service) ActionAllowed(ctx context.Context, roleID string, action authz.Action) (bool, error) {
	return s.repo.ActionAllowed(ctx, roleID, action)
}
