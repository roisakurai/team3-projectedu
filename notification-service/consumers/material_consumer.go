package consumers

import (
	"context"
	"encoding/json"
	"log"

	"notification-service/models"
	"notification-service/services"

	"github.com/redis/go-redis/v9"
)

const ChannelMaterialCreated = "lms:material:created"

type MaterialCreatedEvent struct {
	ClassID       string                         `json:"class_id"`
	ClassName     string                         `json:"class_name"`
	MaterialID    string                         `json:"material_id"`
	MaterialTitle string                         `json:"material_title"`
	Recipients    []models.NotificationRecipient `json:"recipients"`
}

func StartMaterialCreatedConsumer(ctx context.Context, redisClient *redis.Client, svc *services.NotificationService) {
	pubsub := redisClient.Subscribe(ctx, ChannelMaterialCreated)

	go func() {
		defer pubsub.Close()

		log.Println("subscribed to Redis channel:", ChannelMaterialCreated)

		ch := pubsub.Channel()

		for {
			select {
			case <-ctx.Done():
				log.Println("material created consumer stopped")
				return

			case msg, ok := <-ch:
				if !ok {
					log.Println("material created channel closed")
					return
				}

				handleMaterialCreated(svc, msg.Payload)
			}
		}
	}()
}

func handleMaterialCreated(svc *services.NotificationService, payload string) {
	log.Println("received material created event:", payload)

	var event MaterialCreatedEvent

	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		log.Println("failed to parse material created event:", err)
		return
	}

	for _, recipient := range event.Recipients {
		req := models.CreateEmailNotificationRequest{
			UserID:  recipient.UserID,
			Email:   recipient.Email,
			Name:    recipient.Name,
			Title:   "New Material Available",
			Message: "A new material '" + event.MaterialTitle + "' has been posted in class " + event.ClassName,
			Type:    "material_created",
		}

		if err := svc.CreateAndSendEmail(req); err != nil {
			log.Println("failed to send material notification to", recipient.Email, ":", err)
			continue
		}

		log.Println("material notification email sent to:", recipient.Email)
	}
}
