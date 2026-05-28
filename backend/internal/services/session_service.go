package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/models"
	"github.com/yourusername/10000hr/internal/repositories"
)

type SessionService struct {
	sessionRepo *repositories.SessionRepository
}

func NewSessionService(sessionRepo *repositories.SessionRepository) *SessionService {
	return &SessionService{sessionRepo: sessionRepo}
}

func (s *SessionService) CreateSession(ctx context.Context, userID uuid.UUID, req *models.CreateSessionRequest) (*models.Session, error) {
	// Calculate duration if end time is provided
	duration := req.DurationMinutes
	if req.EndTime != nil && duration == 0 {
		duration = int(req.EndTime.Sub(req.StartTime).Minutes())
	}

	session := &models.Session{
		ID:              uuid.New(),
		SkillID:         req.SkillID,
		UserID:          userID,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		DurationMinutes: duration,
		Notes:           req.Notes,
		SessionType:     req.SessionType,
		QualityRating:   req.QualityRating,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) GetSessions(ctx context.Context, userID uuid.UUID, skillID *uuid.UUID, limit, offset int) ([]*models.Session, error) {
	if limit == 0 {
		limit = 50
	}

	if skillID != nil {
		return s.sessionRepo.GetBySkillID(ctx, *skillID, limit, offset)
	}

	return s.sessionRepo.GetByUserID(ctx, userID, limit, offset)
}

func (s *SessionService) GetSession(ctx context.Context, sessionID, userID uuid.UUID) (*models.Session, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}
	if session.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	return session, nil
}

func (s *SessionService) UpdateSession(ctx context.Context, sessionID, userID uuid.UUID, req *models.CreateSessionRequest) (*models.Session, error) {
	session, err := s.GetSession(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}

	// Recalculate duration if end time changed
	if req.EndTime != nil {
		duration := int(req.EndTime.Sub(session.StartTime).Minutes())
		session.DurationMinutes = duration
	} else if req.DurationMinutes > 0 {
		session.DurationMinutes = req.DurationMinutes
	}

	session.EndTime = req.EndTime
	session.Notes = req.Notes
	session.SessionType = req.SessionType
	session.QualityRating = req.QualityRating
	session.UpdatedAt = time.Now()

	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) DeleteSession(ctx context.Context, sessionID, userID uuid.UUID) error {
	session, err := s.GetSession(ctx, sessionID, userID)
	if err != nil {
		return err
	}

	return s.sessionRepo.Delete(ctx, session.ID)
}

// Timer functionality would use Redis to store active timers
func (s *SessionService) StartTimer(ctx context.Context, userID, skillID uuid.UUID) (*models.Session, error) {
	session := &models.Session{
		ID:          uuid.New(),
		SkillID:     skillID,
		UserID:      userID,
		StartTime:   time.Now(),
		SessionType: "focused",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) StopTimer(ctx context.Context, sessionID, userID uuid.UUID) (*models.Session, error) {
	session, err := s.GetSession(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}

	endTime := time.Now()
	duration := int(endTime.Sub(session.StartTime).Minutes())

	session.EndTime = &endTime
	session.DurationMinutes = duration
	session.UpdatedAt = time.Now()

	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}
