package ports

import (
	skillmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/skill"
	"github.com/google/uuid"
)

type Repository interface {
	CreateSkill(sk *skillmodule.Skill) error
	GetSkillByID(skillID uuid.UUID) (skillmodule.Skill, error)
	GetSkills() ([]skillmodule.Skill, error)
	UpdateSkill(sk skillmodule.Skill) error
	DeleteSkill(skillID uuid.UUID) error
}
