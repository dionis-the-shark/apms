package skillmodule

import "github.com/google/uuid"

type Skill struct {
	SkillID uuid.UUID `json:"skill_id"`
	Name    string    `json:"name"`
}
