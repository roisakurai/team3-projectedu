package consumers

import (
	"context"
	"encoding/json"
	"log"

	"notification-service/models"
	"notification-service/services"

	"github.com/redis/go-redis/v9"
)

const ChannelAssignmentCreated = "lms:assignment:created"

type AssignmentCreatedEvent struct {
	ClassID         string                         `json:"class_id"`
	ClassName       string                         `json:"class_name"`
	AssignmentID    string                         `json:"assignment_id"`
	AssignmentTitle string                         `json:"assignment_title"`
	Recipients      []models.NotificationRecipient `json:"recipients"`
}

func StartAssignmentCreatedConsumer(ctx context.Context, redisClient *redis.Client, svc *services.NotificationService) {
	pubsub := redisClient.Subscribe(ctx, ChannelAssignmentCreated)

	go func() {
		defer pubsub.Close()

		log.Println("subscribed to Redis channel:", ChannelAssignmentCreated)

		ch := pubsub.Channel()

		for {
			select {
			case <-ctx.Done():
				log.Println("assignment created consumer stopped")
				return

			case msg, ok := <-ch:
				if !ok {
					log.Println("assignment created channel closed")
					return
				}

				handleAssignmentCreated(svc, msg.Payload)
			}
		}
	}()
}

func handleAssignmentCreated(svc *services.NotificationService, payload string) {
	log.Println("received assignment created event:", payload)

	var event AssignmentCreatedEvent

	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		log.Println("failed to parse assignment created event:", err)
		return
	}

	for _, recipient := range event.Recipients {
		req := models.CreateEmailNotificationRequest{
			UserID:  recipient.UserID,
			Email:   recipient.Email,
			Name:    recipient.Name,
			Title:   "New Assignment Available",
			Message: "A new assignment '" + event.AssignmentTitle + "' has been posted in class " + event.ClassName,
			Type:    "assignment_created",
		}

		if err := svc.CreateAndSendEmail(req); err != nil {
			log.Println("failed to send assignment notification to", recipient.Email, ":", err)
			continue
		}

		log.Println("assignment notification email sent to:", recipient.Email)
	}
}
