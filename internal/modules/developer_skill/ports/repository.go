package ports

import (
	"context"

	developerskillmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/developer_skill"
	"github.com/google/uuid"
)

type Repository interface {
	CreateDeveloperSkill(ctx context.Context, ds developerskillmodule.DeveloperSkill) error
	GetDeveloperSkill(ctx context.Context, developerID, skillID uuid.UUID) (developerskillmodule.DeveloperSkill, error)
	GetDeveloperSkills(ctx context.Context) ([]developerskillmodule.DeveloperSkill, error)
	GetDeveloperSkillsByDeveloper(ctx context.Context, developerID uuid.UUID) ([]developerskillmodule.DeveloperSkill, error)
	UpdateDeveloperSkill(ctx context.Context, oldDeveloperID, oldSkillID, newDeveloperID, newSkillID uuid.UUID) error
	DeleteDeveloperSkill(ctx context.Context, developerID, skillID uuid.UUID) error
}
