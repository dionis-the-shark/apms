package ports

import (
	"context"

	skillmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/skill"
	"github.com/google/uuid"
)

type Repository interface {
	CreateSkill(ctx context.Context, sk *skillmodule.Skill) error
	GetSkillByID(ctx context.Context, skillID uuid.UUID) (skillmodule.Skill, error)
	GetSkills(ctx context.Context) ([]skillmodule.Skill, error)
	UpdateSkill(ctx context.Context, sk skillmodule.Skill) error
	DeleteSkill(ctx context.Context, skillID uuid.UUID) error
}
