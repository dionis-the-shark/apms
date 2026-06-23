package usecase

import (
	"context"

	developerskillmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/developer_skill"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/developer_skill/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo ports.Repository
}

type API interface {
	Create(ctx context.Context, input CreateInput) (developerskillmodule.DeveloperSkill, error)
	GetByIDs(ctx context.Context, developerID, skillID uuid.UUID) (developerskillmodule.DeveloperSkill, error)
	GetAll(ctx context.Context) ([]developerskillmodule.DeveloperSkill, error)
	GetByDeveloper(ctx context.Context, developerID uuid.UUID) ([]developerskillmodule.DeveloperSkill, error)
	Update(ctx context.Context, oldDeveloperID, oldSkillID uuid.UUID, input UpdateInput) (developerskillmodule.DeveloperSkill, error)
	Delete(ctx context.Context, developerID, skillID uuid.UUID) error
}

type CreateInput struct {
	DeveloperID uuid.UUID `json:"developer_id"`
	SkillID     uuid.UUID `json:"skill_id"`
}

type UpdateInput struct {
	DeveloperID uuid.UUID `json:"developer_id"`
	SkillID     uuid.UUID `json:"skill_id"`
}

func New(repo ports.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (developerskillmodule.DeveloperSkill, error) {
	ds := developerskillmodule.DeveloperSkill{
		DeveloperID: input.DeveloperID,
		SkillID:     input.SkillID,
	}
	if err := s.repo.CreateDeveloperSkill(ctx, ds); err != nil {
		return developerskillmodule.DeveloperSkill{}, err
	}
	return ds, nil
}

func (s *Service) GetByIDs(ctx context.Context, developerID, skillID uuid.UUID) (developerskillmodule.DeveloperSkill, error) {
	return s.repo.GetDeveloperSkill(ctx, developerID, skillID)
}

func (s *Service) GetAll(ctx context.Context) ([]developerskillmodule.DeveloperSkill, error) {
	return s.repo.GetDeveloperSkills(ctx)
}

func (s *Service) GetByDeveloper(ctx context.Context, developerID uuid.UUID) ([]developerskillmodule.DeveloperSkill, error) {
	return s.repo.GetDeveloperSkillsByDeveloper(ctx, developerID)
}

func (s *Service) Update(ctx context.Context, oldDeveloperID, oldSkillID uuid.UUID, input UpdateInput) (developerskillmodule.DeveloperSkill, error) {
	if err := s.repo.UpdateDeveloperSkill(ctx, oldDeveloperID, oldSkillID, input.DeveloperID, input.SkillID); err != nil {
		return developerskillmodule.DeveloperSkill{}, err
	}
	return developerskillmodule.DeveloperSkill{
		DeveloperID: input.DeveloperID,
		SkillID:     input.SkillID,
	}, nil
}

func (s *Service) Delete(ctx context.Context, developerID, skillID uuid.UUID) error {
	return s.repo.DeleteDeveloperSkill(ctx, developerID, skillID)
}
