package usecase

import (
	skillmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/skill"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/skill/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo    ports.Repository
	newUUID func() uuid.UUID
}

type API interface {
	Create(input CreateInput) (skillmodule.Skill, error)
	GetByID(skillID uuid.UUID) (skillmodule.Skill, error)
	GetAll() ([]skillmodule.Skill, error)
	Update(skillID uuid.UUID, input UpdateInput) (skillmodule.Skill, error)
	Delete(skillID uuid.UUID) error
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

func (s *Service) Create(input CreateInput) (skillmodule.Skill, error) {
	sk := skillmodule.Skill{
		SkillID: s.newUUID(),
		Name:    input.Name,
	}
	if err := s.repo.CreateSkill(&sk); err != nil {
		return skillmodule.Skill{}, err
	}
	return sk, nil
}

func (s *Service) GetByID(skillID uuid.UUID) (skillmodule.Skill, error) {
	return s.repo.GetSkillByID(skillID)
}

func (s *Service) GetAll() ([]skillmodule.Skill, error) {
	return s.repo.GetSkills()
}

func (s *Service) Update(skillID uuid.UUID, input UpdateInput) (skillmodule.Skill, error) {
	sk := skillmodule.Skill{
		SkillID: skillID,
		Name:    input.Name,
	}
	if err := s.repo.UpdateSkill(sk); err != nil {
		return skillmodule.Skill{}, err
	}
	return sk, nil
}

func (s *Service) Delete(skillID uuid.UUID) error {
	return s.repo.DeleteSkill(skillID)
}
