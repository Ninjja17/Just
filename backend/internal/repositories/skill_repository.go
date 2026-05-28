package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/models"
)

type SkillRepository struct {
	db *sql.DB
}

func NewSkillRepository(db *sql.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

func (r *SkillRepository) Create(ctx context.Context, skill *models.Skill) error {
	query := `
		INSERT INTO skills (id, user_id, name, category, color, icon, target_hours, is_archived, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		skill.ID, skill.UserID, skill.Name, skill.Category, skill.Color,
		skill.Icon, skill.TargetHours, skill.IsArchived, skill.CreatedAt, skill.UpdatedAt,
	)
	return err
}

func (r *SkillRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Skill, error) {
	skill := &models.Skill{}
	query := `
		SELECT id, user_id, name, category, color, icon, target_hours, is_archived, created_at, updated_at
		FROM skills WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&skill.ID, &skill.UserID, &skill.Name, &skill.Category, &skill.Color,
		&skill.Icon, &skill.TargetHours, &skill.IsArchived, &skill.CreatedAt, &skill.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return skill, err
}

func (r *SkillRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Skill, error) {
	query := `
		SELECT id, user_id, name, category, color, icon, target_hours, is_archived, created_at, updated_at
		FROM skills WHERE user_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*models.Skill
	for rows.Next() {
		skill := &models.Skill{}
		err := rows.Scan(
			&skill.ID, &skill.UserID, &skill.Name, &skill.Category, &skill.Color,
			&skill.Icon, &skill.TargetHours, &skill.IsArchived, &skill.CreatedAt, &skill.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}
	return skills, nil
}

func (r *SkillRepository) Update(ctx context.Context, skill *models.Skill) error {
	query := `
		UPDATE skills
		SET name = $2, category = $3, color = $4, icon = $5, target_hours = $6, is_archived = $7, updated_at = $8
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		skill.ID, skill.Name, skill.Category, skill.Color, skill.Icon,
		skill.TargetHours, skill.IsArchived, time.Now(),
	)
	return err
}

func (r *SkillRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM skills WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *SkillRepository) GetTotalMinutes(ctx context.Context, skillID uuid.UUID) (int, error) {
	var total int
	query := `SELECT COALESCE(SUM(duration_minutes), 0) FROM sessions WHERE skill_id = $1`
	err := r.db.QueryRowContext(ctx, query, skillID).Scan(&total)
	return total, err
}
