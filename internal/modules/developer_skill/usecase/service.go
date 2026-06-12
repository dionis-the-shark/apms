package usecase

import (
	developerskillmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/developer_skill"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/developer_skill/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo ports.Repository
}

type API interface {
	Create(input CreateInput) (developerskillmodule.DeveloperSkill, error)
	GetByIDs(developerID, skillID uuid.UUID) (developerskillmodule.DeveloperSkill, error)
	GetAll() ([]developerskillmodule.DeveloperSkill, error)
	GetByDeveloper(developerID uuid.UUID) ([]developerskillmodule.DeveloperSkill, error)
	Update(oldDeveloperID, oldSkillID uuid.UUID, input UpdateInput) (developerskillmodule.DeveloperSkill, error)
	Delete(developerID, skillID uuid.UUID) error
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

func (s *Service) Create(input CreateInput) (developerskillmodule.DeveloperSkill, error) {
	ds := developerskillmodule.DeveloperSkill{
		DeveloperID: input.DeveloperID,
		SkillID:     input.SkillID,
	}
	if err := s.repo.CreateDeveloperSkill(ds); err != nil {
		return developerskillmodule.DeveloperSkill{}, err
	}
	return ds, nil
}

func (s *Service) GetByIDs(developerID, skillID uuid.UUID) (developerskillmodule.DeveloperSkill, error) {
	return s.repo.GetDeveloperSkill(developerID, skillID)
}

func (s *Service) GetAll() ([]developerskillmodule.DeveloperSkill, error) {
	return s.repo.GetDeveloperSkills()
}

func (s *Service) GetByDeveloper(developerID uuid.UUID) ([]developerskillmodule.DeveloperSkill, error) {
	return s.repo.GetDeveloperSkillsByDeveloper(developerID)
}

func (s *Service) Update(oldDeveloperID, oldSkillID uuid.UUID, input UpdateInput) (developerskillmodule.DeveloperSkill, error) {
	if err := s.repo.UpdateDeveloperSkill(oldDeveloperID, oldSkillID, input.DeveloperID, input.SkillID); err != nil {
		return developerskillmodule.DeveloperSkill{}, err
	}
	return developerskillmodule.DeveloperSkill{
		DeveloperID: input.DeveloperID,
		SkillID:     input.SkillID,
	}, nil
}

func (s *Service) Delete(developerID, skillID uuid.UUID) error {
	return s.repo.DeleteDeveloperSkill(developerID, skillID)
}
