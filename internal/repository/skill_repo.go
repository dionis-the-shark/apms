package repository

import (
	"database/sql"

	"github.com/dionis-the-shark/apms-task-tracker/internal/models/skill"
	"github.com/google/uuid"
)

type SkillRepository struct {
	DB *sql.DB
}

type skillScanner interface {
	Scan(dest ...any) error
}

func scanSkill(s skillScanner) (skill.Skill, error) {
	var sk skill.Skill
	err := s.Scan(&sk.SkillID, &sk.Name)
	return sk, err
}

func (r *SkillRepository) CreateSkill(sk *skill.Skill) error {
	if sk.SkillID == uuid.Nil {
		sk.SkillID = uuid.New()
	}

	query := `INSERT INTO skills (skill_id, name)
			  VALUES ($1, $2)`
	_, err := r.DB.Exec(query, sk.SkillID, sk.Name)
	return err
}

func (r *SkillRepository) GetSkillByID(skillID uuid.UUID) (skill.Skill, error) {
	query := `SELECT skill_id, name FROM skills WHERE skill_id = $1`
	return scanSkill(r.DB.QueryRow(query, skillID))
}

func (r *SkillRepository) GetSkills() ([]skill.Skill, error) {
	rows, err := r.DB.Query(`SELECT skill_id, name FROM skills`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []skill.Skill
	for rows.Next() {
		sk, err := scanSkill(rows)
		if err != nil {
			return nil, err
		}
		skills = append(skills, sk)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return skills, nil
}

func (r *SkillRepository) UpdateSkill(sk skill.Skill) error {
	query := `UPDATE skills SET name = $1 WHERE skill_id = $2`
	result, err := r.DB.Exec(query, sk.Name, sk.SkillID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *SkillRepository) DeleteSkill(skillID uuid.UUID) error {
	result, err := r.DB.Exec(`DELETE FROM skills WHERE skill_id = $1`, skillID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
