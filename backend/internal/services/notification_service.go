package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/10000hr/internal/models"
	"github.com/yourusername/10000hr/internal/repositories"
)

type NotificationService struct {
	notificationRepo *repositories.NotificationRepository
}

func NewNotificationService(notificationRepo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{notificationRepo: notificationRepo}
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID uuid.UUID, notifType, message string) error {
	notification := &models.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notifType,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}
	return s.notificationRepo.Create(ctx, notification)
}

func (s *NotificationService) GetNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	if limit == 0 {
		limit = 20
	}
	return s.notificationRepo.GetByUserID(ctx, userID, limit, offset)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(ctx, notificationID)
}

func (s *NotificationService) Delete(ctx context.Context, notificationID uuid.UUID) error {
	return s.notificationRepo.Delete(ctx, notificationID)
}
