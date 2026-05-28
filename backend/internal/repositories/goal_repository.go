package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/models"
)

type GoalRepository struct {
	db *sql.DB
}

func NewGoalRepository(db *sql.DB) *GoalRepository {
	return &GoalRepository{db: db}
}

func (r *GoalRepository) Create(ctx context.Context, goal *models.Goal) error {
	query := `
		INSERT INTO goals (id, user_id, skill_id, type, target_hours, deadline, is_completed, created_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		goal.ID, goal.UserID, goal.SkillID, goal.Type, goal.TargetHours,
		goal.Deadline, goal.IsCompleted, goal.CreatedAt, goal.CompletedAt,
	)
	return err
}

func (r *GoalRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Goal, error) {
	query := `
		SELECT id, user_id, skill_id, type, target_hours, deadline, is_completed, created_at, completed_at
		FROM goals WHERE user_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []*models.Goal
	for rows.Next() {
		goal := &models.Goal{}
		err := rows.Scan(
			&goal.ID, &goal.UserID, &goal.SkillID, &goal.Type, &goal.TargetHours,
			&goal.Deadline, &goal.IsCompleted, &goal.CreatedAt, &goal.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		goals = append(goals, goal)
	}
	return goals, nil
}

func (r *GoalRepository) Update(ctx context.Context, goal *models.Goal) error {
	query := `
		UPDATE goals
		SET target_hours = $2, deadline = $3, is_completed = $4, completed_at = $5
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		goal.ID, goal.TargetHours, goal.Deadline, goal.IsCompleted, goal.CompletedAt,
	)
	return err
}

func (r *GoalRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM goals WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
