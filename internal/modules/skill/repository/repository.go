package repository

import (
	"database/sql"

	skillmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/skill"
	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

type skillScanner interface {
	Scan(dest ...any) error
}

func scanSkill(s skillScanner) (skillmodule.Skill, error) {
	var sk skillmodule.Skill
	err := s.Scan(&sk.SkillID, &sk.Name)
	return sk, err
}

func (r *Repository) CreateSkill(sk *skillmodule.Skill) error {
	if sk.SkillID == uuid.Nil {
		sk.SkillID = uuid.New()
	}

	query := `INSERT INTO skills (skill_id, name)
			  VALUES ($1, $2)`
	_, err := r.DB.Exec(query, sk.SkillID, sk.Name)
	return err
}

func (r *Repository) GetSkillByID(skillID uuid.UUID) (skillmodule.Skill, error) {
	query := `SELECT skill_id, name FROM skills WHERE skill_id = $1`
	return scanSkill(r.DB.QueryRow(query, skillID))
}

func (r *Repository) GetSkills() ([]skillmodule.Skill, error) {
	rows, err := r.DB.Query(`SELECT skill_id, name FROM skills`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []skillmodule.Skill
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

func (r *Repository) UpdateSkill(sk skillmodule.Skill) error {
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

func (r *Repository) DeleteSkill(skillID uuid.UUID) error {
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
