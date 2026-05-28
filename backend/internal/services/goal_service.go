package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/models"
	"github.com/yourusername/10000hr/internal/repositories"
)

type GoalService struct {
	goalRepo *repositories.GoalRepository
}

func NewGoalService(goalRepo *repositories.GoalRepository) *GoalService {
	return &GoalService{goalRepo: goalRepo}
}

func (s *GoalService) CreateGoal(ctx context.Context, userID uuid.UUID, req *models.CreateGoalRequest) (*models.Goal, error) {
	goal := &models.Goal{
		ID:          uuid.New(),
		UserID:      userID,
		SkillID:     req.SkillID,
		Type:        req.Type,
		TargetHours: req.TargetHours,
		Deadline:    req.Deadline,
		IsCompleted: false,
	}

	if err := s.goalRepo.Create(ctx, goal); err != nil {
		return nil, err
	}

	return goal, nil
}

func (s *GoalService) GetGoals(ctx context.Context, userID uuid.UUID) ([]*models.Goal, error) {
	return s.goalRepo.GetByUserID(ctx, userID)
}

func (s *GoalService) UpdateGoal(ctx context.Context, goalID, userID uuid.UUID, req *models.CreateGoalRequest) (*models.Goal, error) {
	goals, err := s.goalRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var goal *models.Goal
	for _, g := range goals {
		if g.ID == goalID {
			goal = g
			break
		}
	}

	if goal == nil {
		return nil, fmt.Errorf("goal not found")
	}

	goal.TargetHours = req.TargetHours
	goal.Deadline = req.Deadline

	if err := s.goalRepo.Update(ctx, goal); err != nil {
		return nil, err
	}

	return goal, nil
}

func (s *GoalService) DeleteGoal(ctx context.Context, goalID, userID uuid.UUID) error {
	return s.goalRepo.Delete(ctx, goalID)
}
