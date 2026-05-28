package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/models"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (id, skill_id, user_id, start_time, end_time, duration_minutes, notes, session_type, quality_rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.SkillID, session.UserID, session.StartTime, session.EndTime,
		session.DurationMinutes, session.Notes, session.SessionType, session.QualityRating,
		session.CreatedAt, session.UpdatedAt,
	)
	return err
}

func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	session := &models.Session{}
	query := `
		SELECT id, skill_id, user_id, start_time, end_time, duration_minutes, notes, session_type, quality_rating, created_at, updated_at
		FROM sessions WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID, &session.SkillID, &session.UserID, &session.StartTime, &session.EndTime,
		&session.DurationMinutes, &session.Notes, &session.SessionType, &session.QualityRating,
		&session.CreatedAt, &session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

func (r *SessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Session, error) {
	query := `
		SELECT id, skill_id, user_id, start_time, end_time, duration_minutes, notes, session_type, quality_rating, created_at, updated_at
		FROM sessions WHERE user_id = $1 ORDER BY start_time DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		err := rows.Scan(
			&session.ID, &session.SkillID, &session.UserID, &session.StartTime, &session.EndTime,
			&session.DurationMinutes, &session.Notes, &session.SessionType, &session.QualityRating,
			&session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *SessionRepository) GetBySkillID(ctx context.Context, skillID uuid.UUID, limit, offset int) ([]*models.Session, error) {
	query := `
		SELECT id, skill_id, user_id, start_time, end_time, duration_minutes, notes, session_type, quality_rating, created_at, updated_at
		FROM sessions WHERE skill_id = $1 ORDER BY start_time DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, skillID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		err := rows.Scan(
			&session.ID, &session.SkillID, &session.UserID, &session.StartTime, &session.EndTime,
			&session.DurationMinutes, &session.Notes, &session.SessionType, &session.QualityRating,
			&session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *SessionRepository) Update(ctx context.Context, session *models.Session) error {
	query := `
		UPDATE sessions
		SET end_time = $2, duration_minutes = $3, notes = $4, session_type = $5, quality_rating = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.EndTime, session.DurationMinutes, session.Notes,
		session.SessionType, session.QualityRating, time.Now(),
	)
	return err
}

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
