package ports

import (
	developerskillmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/developer_skill"
	"github.com/google/uuid"
)

type Repository interface {
	CreateDeveloperSkill(ds developerskillmodule.DeveloperSkill) error
	GetDeveloperSkill(developerID, skillID uuid.UUID) (developerskillmodule.DeveloperSkill, error)
	GetDeveloperSkills() ([]developerskillmodule.DeveloperSkill, error)
	GetDeveloperSkillsByDeveloper(developerID uuid.UUID) ([]developerskillmodule.DeveloperSkill, error)
	UpdateDeveloperSkill(oldDeveloperID, oldSkillID, newDeveloperID, newSkillID uuid.UUID) error
	DeleteDeveloperSkill(developerID, skillID uuid.UUID) error
}
