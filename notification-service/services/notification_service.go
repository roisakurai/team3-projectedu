package services

import (
	"errors"
	"time"

	"notification-service/models"
	"notification-service/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NotificationService struct {
	Repo *repositories.NotificationRepository
}

func NewNotificationService(repo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{Repo: repo}
}

func (s *NotificationService) Create(req models.CreateNotificationRequest) error {
	if len(req.UserIDs) == 0 {
		return errors.New("user_ids is required")
	}

	if req.Title == "" || req.Message == "" || req.Type == "" {
		return errors.New("title, message, and type are required")
	}

	var notifications []interface{}

	for _, userID := range req.UserIDs {
		notification := models.Notification{
			ID:        bson.NewObjectID(),
			UserID:    userID,
			Title:     req.Title,
			Message:   req.Message,
			Type:      req.Type,
			IsRead:    false,
			CreatedAt: time.Now(),
		}

		notifications = append(notifications, notification)
	}

	return s.Repo.CreateMany(notifications)
}

func (s *NotificationService) GetByUserID(userID string) ([]models.Notification, error) {
	return s.Repo.FindByUserID(userID)
}

func (s *NotificationService) MarkAsRead(id string) error {
	return s.Repo.MarkAsRead(id)
}
