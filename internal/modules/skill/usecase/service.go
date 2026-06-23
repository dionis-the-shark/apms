package usecase

import (
	"context"

	skillmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/skill"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/skill/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo    ports.Repository
	newUUID func() uuid.UUID
}

type API interface {
	Create(ctx context.Context, input CreateInput) (skillmodule.Skill, error)
	GetByID(ctx context.Context, skillID uuid.UUID) (skillmodule.Skill, error)
	GetAll(ctx context.Context) ([]skillmodule.Skill, error)
	Update(ctx context.Context, skillID uuid.UUID, input UpdateInput) (skillmodule.Skill, error)
	Delete(ctx context.Context, skillID uuid.UUID) error
}

type CreateInput struct {
	Name string `json:"name"`
}

type UpdateInput struct {
	Name string `json:"name"`
}

func New(repo ports.Repository) *Service {
	return &Service{
		repo:    repo,
		newUUID: uuid.New,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (skillmodule.Skill, error) {
	sk := skillmodule.Skill{
		SkillID: s.newUUID(),
		Name:    input.Name,
	}
	if err := s.repo.CreateSkill(ctx, &sk); err != nil {
		return skillmodule.Skill{}, err
	}
	return sk, nil
}

func (s *Service) GetByID(ctx context.Context, skillID uuid.UUID) (skillmodule.Skill, error) {
	return s.repo.GetSkillByID(ctx, skillID)
}

func (s *Service) GetAll(ctx context.Context) ([]skillmodule.Skill, error) {
	return s.repo.GetSkills(ctx)
}

func (s *Service) Update(ctx context.Context, skillID uuid.UUID, input UpdateInput) (skillmodule.Skill, error) {
	sk := skillmodule.Skill{
		SkillID: skillID,
		Name:    input.Name,
	}
	if err := s.repo.UpdateSkill(ctx, sk); err != nil {
		return skillmodule.Skill{}, err
	}
	return sk, nil
}

func (s *Service) Delete(ctx context.Context, skillID uuid.UUID) error {
	return s.repo.DeleteSkill(ctx, skillID)
}
