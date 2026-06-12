package skillservice

import (
	"github.com/dionis-the-shark/apms-task-tracker/internal/models/skill"
	"github.com/google/uuid"
)

type Repository interface {
	CreateSkill(sk *skill.Skill) error
	GetSkillByID(skillID uuid.UUID) (skill.Skill, error)
	GetSkills() ([]skill.Skill, error)
	UpdateSkill(sk skill.Skill) error
	DeleteSkill(skillID uuid.UUID) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(sk *skill.Skill) error {
	return s.repo.CreateSkill(sk)
}

func (s *Service) GetByID(skillID uuid.UUID) (skill.Skill, error) {
	return s.repo.GetSkillByID(skillID)
}

func (s *Service) GetAll() ([]skill.Skill, error) {
	return s.repo.GetSkills()
}

func (s *Service) Update(sk skill.Skill) error {
	return s.repo.UpdateSkill(sk)
}

func (s *Service) Delete(skillID uuid.UUID) error {
	return s.repo.DeleteSkill(skillID)
}
