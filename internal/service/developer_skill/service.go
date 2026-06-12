package developerskillservice

import (
	"github.com/dionis-the-shark/apms-task-tracker/internal/models/developer_skill"
	"github.com/google/uuid"
)

type Repository interface {
	CreateDeveloperSkill(ds developer_skill.DeveloperSkill) error
	GetDeveloperSkill(developerID, skillID uuid.UUID) (developer_skill.DeveloperSkill, error)
	GetDeveloperSkills() ([]developer_skill.DeveloperSkill, error)
	GetDeveloperSkillsByDeveloper(developerID uuid.UUID) ([]developer_skill.DeveloperSkill, error)
	UpdateDeveloperSkill(oldDeveloperID, oldSkillID, newDeveloperID, newSkillID uuid.UUID) error
	DeleteDeveloperSkill(developerID, skillID uuid.UUID) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ds developer_skill.DeveloperSkill) error {
	return s.repo.CreateDeveloperSkill(ds)
}

func (s *Service) GetByIDs(developerID, skillID uuid.UUID) (developer_skill.DeveloperSkill, error) {
	return s.repo.GetDeveloperSkill(developerID, skillID)
}

func (s *Service) GetAll() ([]developer_skill.DeveloperSkill, error) {
	return s.repo.GetDeveloperSkills()
}

func (s *Service) GetByDeveloper(developerID uuid.UUID) ([]developer_skill.DeveloperSkill, error) {
	return s.repo.GetDeveloperSkillsByDeveloper(developerID)
}

func (s *Service) Update(oldDeveloperID, oldSkillID, newDeveloperID, newSkillID uuid.UUID) error {
	return s.repo.UpdateDeveloperSkill(oldDeveloperID, oldSkillID, newDeveloperID, newSkillID)
}

func (s *Service) Delete(developerID, skillID uuid.UUID) error {
	return s.repo.DeleteDeveloperSkill(developerID, skillID)
}
