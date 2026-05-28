package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/models"
	"github.com/yourusername/10000hr/internal/repositories"
)

type SocialService struct {
	userRepo     *repositories.UserRepository
	followerRepo *repositories.FollowerRepository
}

func NewSocialService(userRepo *repositories.UserRepository, followerRepo *repositories.FollowerRepository) *SocialService {
	return &SocialService{
		userRepo:     userRepo,
		followerRepo: followerRepo,
	}
}

func (s *SocialService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *SocialService) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return fmt.Errorf("cannot follow yourself")
	}
	return s.followerRepo.Create(ctx, followerID, followingID)
}

func (s *SocialService) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return s.followerRepo.Delete(ctx, followerID, followingID)
}

func (s *SocialService) GetFollowers(ctx context.Context, userID uuid.UUID) ([]*models.User, error) {
	return s.followerRepo.GetFollowers(ctx, userID)
}

func (s *SocialService) GetFollowing(ctx context.Context, userID uuid.UUID) ([]*models.User, error) {
	return s.followerRepo.GetFollowing(ctx, userID)
}
