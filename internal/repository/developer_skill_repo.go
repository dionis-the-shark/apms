package repository

import (
	"database/sql"

	"github.com/dionis-the-shark/apms-task-tracker/internal/models/developer_skill"
	"github.com/google/uuid"
)

type DeveloperSkillRepository struct {
	DB *sql.DB
}

type developerSkillScanner interface {
	Scan(dest ...any) error
}

func scanDeveloperSkill(s developerSkillScanner) (developer_skill.DeveloperSkill, error) {
	var ds developer_skill.DeveloperSkill
	err := s.Scan(&ds.DeveloperID, &ds.SkillID)
	return ds, err
}

func (r *DeveloperSkillRepository) CreateDeveloperSkill(ds developer_skill.DeveloperSkill) error {
	query := `INSERT INTO developer_skills (developer_id, skill_id)
			  VALUES ($1, $2)`
	_, err := r.DB.Exec(query, ds.DeveloperID, ds.SkillID)
	return err
}

func (r *DeveloperSkillRepository) GetDeveloperSkill(developerID, skillID uuid.UUID) (developer_skill.DeveloperSkill, error) {
	query := `SELECT developer_id, skill_id FROM developer_skills
			  WHERE developer_id = $1 AND skill_id = $2`
	return scanDeveloperSkill(r.DB.QueryRow(query, developerID, skillID))
}

func (r *DeveloperSkillRepository) GetDeveloperSkills() ([]developer_skill.DeveloperSkill, error) {
	rows, err := r.DB.Query(`SELECT developer_id, skill_id FROM developer_skills`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []developer_skill.DeveloperSkill
	for rows.Next() {
		ds, err := scanDeveloperSkill(rows)
		if err != nil {
			return nil, err
		}
		skills = append(skills, ds)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return skills, nil
}

func (r *DeveloperSkillRepository) GetDeveloperSkillsByDeveloper(developerID uuid.UUID) ([]developer_skill.DeveloperSkill, error) {
	rows, err := r.DB.Query(`SELECT developer_id, skill_id FROM developer_skills WHERE developer_id = $1`, developerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []developer_skill.DeveloperSkill
	for rows.Next() {
		ds, err := scanDeveloperSkill(rows)
		if err != nil {
			return nil, err
		}
		skills = append(skills, ds)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return skills, nil
}

func (r *DeveloperSkillRepository) UpdateDeveloperSkill(oldDeveloperID, oldSkillID, newDeveloperID, newSkillID uuid.UUID) error {
	query := `UPDATE developer_skills SET developer_id = $1, skill_id = $2
			  WHERE developer_id = $3 AND skill_id = $4`
	result, err := r.DB.Exec(query, newDeveloperID, newSkillID, oldDeveloperID, oldSkillID)
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

func (r *DeveloperSkillRepository) DeleteDeveloperSkill(developerID, skillID uuid.UUID) error {
	result, err := r.DB.Exec(
		`DELETE FROM developer_skills WHERE developer_id = $1 AND skill_id = $2`,
		developerID,
		skillID,
	)
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
