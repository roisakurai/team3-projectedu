package services

import (
	"errors"
	"fmt"
	"time"

	"notification-service/models"
	"notification-service/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NotificationService struct {
	Repo         *repositories.NotificationRepository
	EmailService *EmailService
}

func NewNotificationService(repo *repositories.NotificationRepository, emailService *EmailService) *NotificationService {
	return &NotificationService{
		Repo:         repo,
		EmailService: emailService,
	}
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

func (s *NotificationService) CreateAndSendEmail(req models.CreateEmailNotificationRequest) error {
	if req.UserID == "" || req.Email == "" {
		return errors.New("user_id and email are required")
	}

	if req.Title == "" || req.Message == "" || req.Type == "" {
		return errors.New("title, message, and type are required")
	}

	notificationReq := models.CreateNotificationRequest{
		UserIDs: []string{req.UserID},
		Title:   req.Title,
		Message: req.Message,
		Type:    req.Type,
	}

	if err := s.Create(notificationReq); err != nil {
		return err
	}

	html := fmt.Sprintf(`
		<h2>%s</h2>
		<p>Hello %s,</p>
		<p>%s</p>
	`, req.Title, req.Name, req.Message)

	return s.EmailService.SendEmail(
		req.Email,
		req.Name,
		req.Title,
		req.Message,
		html,
	)
}
