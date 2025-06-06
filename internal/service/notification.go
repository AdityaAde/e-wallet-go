package service

import (
	"context"
	"log"

	"adityaad.id/belajar-auth/domain"
	"adityaad.id/belajar-auth/dto"
)

type notificationService struct {
	notificationRepository domain.NotificationRepository
}

func NewNotificationService(notificationRepository domain.NotificationRepository) domain.NotificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
	}
}

// FindByUser implements domain.NotificationService.
func (n notificationService) FindByUser(ctx context.Context, user int64) ([]dto.NotificationData, error) {

	notifications, err := n.notificationRepository.FindByUser(ctx, user)
	if err != nil {
		log.Printf("error finding notifications: %s", err.Error())
		return nil, err
	}

	var result []dto.NotificationData
	for _, notification := range notifications {
		result = append(result, dto.NotificationData{
			ID:        notification.ID,
			Title:     notification.Title,
			Body:      notification.Body,
			Status:    notification.Status,
			IsRead:    notification.IsRead,
			CreatedAt: notification.CreatedAt,
		})
	}

	if result == nil {
		result = make([]dto.NotificationData, 0)
	}

	return result, nil
}
