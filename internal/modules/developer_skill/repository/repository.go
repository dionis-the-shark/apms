package repository

import (
	"context"
	"database/sql"

	developerskillmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/developer_skill"
	"github.com/google/uuid"
)

type Repository struct {
	DB *sql.DB
}

type developerSkillScanner interface {
	Scan(dest ...any) error
}

func scanDeveloperSkill(s developerSkillScanner) (developerskillmodule.DeveloperSkill, error) {
	var ds developerskillmodule.DeveloperSkill
	err := s.Scan(&ds.DeveloperID, &ds.SkillID)
	return ds, err
}

func (r *Repository) CreateDeveloperSkill(ctx context.Context, ds developerskillmodule.DeveloperSkill) error {
	query := `INSERT INTO developer_skills (developer_id, skill_id)
			  VALUES ($1, $2)`
	_, err := r.DB.ExecContext(ctx, query, ds.DeveloperID, ds.SkillID)
	return err
}

func (r *Repository) GetDeveloperSkill(ctx context.Context, developerID, skillID uuid.UUID) (developerskillmodule.DeveloperSkill, error) {
	query := `SELECT developer_id, skill_id FROM developer_skills
			  WHERE developer_id = $1 AND skill_id = $2`
	return scanDeveloperSkill(r.DB.QueryRowContext(ctx, query, developerID, skillID))
}

func (r *Repository) GetDeveloperSkills(ctx context.Context) ([]developerskillmodule.DeveloperSkill, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT developer_id, skill_id FROM developer_skills`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []developerskillmodule.DeveloperSkill
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

func (r *Repository) GetDeveloperSkillsByDeveloper(ctx context.Context, developerID uuid.UUID) ([]developerskillmodule.DeveloperSkill, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT developer_id, skill_id FROM developer_skills WHERE developer_id = $1`, developerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []developerskillmodule.DeveloperSkill
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

func (r *Repository) UpdateDeveloperSkill(ctx context.Context, oldDeveloperID, oldSkillID, newDeveloperID, newSkillID uuid.UUID) error {
	query := `UPDATE developer_skills SET developer_id = $1, skill_id = $2
			  WHERE developer_id = $3 AND skill_id = $4`
	result, err := r.DB.ExecContext(ctx, query, newDeveloperID, newSkillID, oldDeveloperID, oldSkillID)
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

func (r *Repository) DeleteDeveloperSkill(ctx context.Context, developerID, skillID uuid.UUID) error {
	result, err := r.DB.ExecContext(
		ctx,
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
