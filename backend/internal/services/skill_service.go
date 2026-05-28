package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/models"
	"github.com/yourusername/10000hr/internal/repositories"
)

type SkillService struct {
	skillRepo *repositories.SkillRepository
}

func NewSkillService(skillRepo *repositories.SkillRepository) *SkillService {
	return &SkillService{skillRepo: skillRepo}
}

func (s *SkillService) CreateSkill(ctx context.Context, userID uuid.UUID, req *models.CreateSkillRequest) (*models.Skill, error) {
	targetHours := req.TargetHours
	if targetHours == 0 {
		targetHours = 10000
	}

	skill := &models.Skill{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Category:    req.Category,
		Color:       req.Color,
		Icon:        req.Icon,
		TargetHours: targetHours,
		IsArchived:  false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.skillRepo.Create(ctx, skill); err != nil {
		return nil, err
	}

	return skill, nil
}

func (s *SkillService) GetSkills(ctx context.Context, userID uuid.UUID) ([]*models.Skill, error) {
	return s.skillRepo.GetByUserID(ctx, userID)
}

func (s *SkillService) GetSkill(ctx context.Context, skillID, userID uuid.UUID) (*models.Skill, error) {
	skill, err := s.skillRepo.GetByID(ctx, skillID)
	if err != nil {
		return nil, err
	}
	if skill == nil {
		return nil, fmt.Errorf("skill not found")
	}
	if skill.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	return skill, nil
}

func (s *SkillService) UpdateSkill(ctx context.Context, skillID, userID uuid.UUID, req *models.CreateSkillRequest) (*models.Skill, error) {
	skill, err := s.GetSkill(ctx, skillID, userID)
	if err != nil {
		return nil, err
	}

	skill.Name = req.Name
	skill.Category = req.Category
	skill.Color = req.Color
	skill.Icon = req.Icon
	if req.TargetHours > 0 {
		skill.TargetHours = req.TargetHours
	}
	skill.UpdatedAt = time.Now()

	if err := s.skillRepo.Update(ctx, skill); err != nil {
		return nil, err
	}

	return skill, nil
}

func (s *SkillService) DeleteSkill(ctx context.Context, skillID, userID uuid.UUID) error {
	skill, err := s.GetSkill(ctx, skillID, userID)
	if err != nil {
		return err
	}

	return s.skillRepo.Delete(ctx, skill.ID)
}

func (s *SkillService) GetSkillStats(ctx context.Context, skillID, userID uuid.UUID) (map[string]interface{}, error) {
	skill, err := s.GetSkill(ctx, skillID, userID)
	if err != nil {
		return nil, err
	}

	totalMinutes, err := s.skillRepo.GetTotalMinutes(ctx, skillID)
	if err != nil {
		return nil, err
	}

	totalHours := float64(totalMinutes) / 60.0
	progress := (totalHours / float64(skill.TargetHours)) * 100

	return map[string]interface{}{
		"skill_id":      skill.ID,
		"total_minutes": totalMinutes,
		"total_hours":   totalHours,
		"target_hours":  skill.TargetHours,
		"progress":      progress,
	}, nil
}
