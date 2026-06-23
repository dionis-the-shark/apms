package developerskillmodule

import "github.com/google/uuid"

type DeveloperSkill struct {
	DeveloperID uuid.UUID `json:"developer_id"`
	SkillID     uuid.UUID `json:"skill_id"`
}
