package usecase

import (
	rolemodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/role"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/role/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo    ports.Repository
	newUUID func() uuid.UUID
}

type API interface {
	Create(input CreateInput) (rolemodule.Role, error)
	GetByID(roleID uuid.UUID) (rolemodule.Role, error)
	GetAll() ([]rolemodule.Role, error)
	Update(roleID uuid.UUID, input UpdateInput) (rolemodule.Role, error)
	Delete(roleID uuid.UUID) error
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

func (s *Service) Create(input CreateInput) (rolemodule.Role, error) {
	status := input.Status
	if status == "" {
		status = "active"
	}

	r := rolemodule.Role{
		RoleID: s.newUUID(),
		Title:  input.Title,
		Status: status,
	}
	if err := s.repo.CreateRole(&r); err != nil {
		return rolemodule.Role{}, err
	}
	return r, nil
}

func (s *Service) GetByID(roleID uuid.UUID) (rolemodule.Role, error) {
	return s.repo.GetRoleByID(roleID)
}

func (s *Service) GetAll() ([]rolemodule.Role, error) {
	return s.repo.GetRoles()
}

func (s *Service) Update(roleID uuid.UUID, input UpdateInput) (rolemodule.Role, error) {
	status := input.Status
	if status == "" {
		status = "active"
	}

	r := rolemodule.Role{
		RoleID: roleID,
		Title:  input.Title,
		Status: status,
	}
	if err := s.repo.UpdateRole(r); err != nil {
		return rolemodule.Role{}, err
	}
	return r, nil
}

func (s *Service) Delete(roleID uuid.UUID) error {
	return s.repo.DeleteRole(roleID)
}
