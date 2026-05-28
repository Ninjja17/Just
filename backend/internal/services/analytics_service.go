package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/repositories"
)

type AnalyticsService struct {
	sessionRepo *repositories.SessionRepository
	skillRepo   *repositories.SkillRepository
}

func NewAnalyticsService(sessionRepo *repositories.SessionRepository, skillRepo *repositories.SkillRepository) *AnalyticsService {
	return &AnalyticsService{
		sessionRepo: sessionRepo,
		skillRepo:   skillRepo,
	}
}

func (s *AnalyticsService) GetOverview(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	skills, err := s.skillRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var totalMinutes int
	skillStats := make([]map[string]interface{}, 0)

	for _, skill := range skills {
		minutes, _ := s.skillRepo.GetTotalMinutes(ctx, skill.ID)
		totalMinutes += minutes

		skillStats = append(skillStats, map[string]interface{}{
			"skill_id":      skill.ID,
			"skill_name":    skill.Name,
			"total_minutes": minutes,
			"total_hours":   float64(minutes) / 60.0,
		})
	}

	return map[string]interface{}{
		"total_skills":  len(skills),
		"total_minutes": totalMinutes,
		"total_hours":   float64(totalMinutes) / 60.0,
		"skills":        skillStats,
	}, nil
}

func (s *AnalyticsService) GetTrends(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// This would implement trend analysis based on historical data
	// For now, return placeholder
	return map[string]interface{}{
		"message": "Trend analysis coming soon",
	}, nil
}

func (s *AnalyticsService) GetPredictions(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// This would implement prediction algorithms
	// For now, return placeholder
	return map[string]interface{}{
		"message": "Predictions coming soon",
	}, nil
}
