package consumers

import (
	"context"
	"encoding/json"
	"log"

	"notification-service/models"
	"notification-service/services"

	"github.com/redis/go-redis/v9"
)

const ChannelAssignmentGraded = "lms:assignment:graded"

type AssignmentGradedEvent struct {
	UserID          string  `json:"user_id"`
	StudentEmail    string  `json:"student_email"`
	StudentName     string  `json:"student_name"`
	AssignmentID    string  `json:"assignment_id"`
	AssignmentTitle string  `json:"assignment_title"`
	ClassName       string  `json:"class_name"`
	Grade           float64 `json:"grade"`
}

func StartAssignmentGradedConsumer(ctx context.Context, redisClient *redis.Client, svc *services.NotificationService) {
	pubsub := redisClient.Subscribe(ctx, ChannelAssignmentGraded)

	go func() {
		defer pubsub.Close()

		log.Println("subscribed to Redis channel:", ChannelAssignmentGraded)

		ch := pubsub.Channel()

		for {
			select {
			case <-ctx.Done():
				log.Println("assignment graded consumer stopped")
				return

			case msg, ok := <-ch:
				if !ok {
					log.Println("assignment graded channel closed")
					return
				}

				handleAssignmentGraded(svc, msg.Payload)
			}
		}
	}()
}

func handleAssignmentGraded(svc *services.NotificationService, payload string) {
	log.Println("received assignment graded event:", payload)

	var event AssignmentGradedEvent

	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		log.Println("failed to parse assignment graded event:", err)
		return
	}

	req := models.CreateEmailNotificationRequest{
		UserID: event.UserID,
		Email:  event.StudentEmail,
		Name:   event.StudentName,
		Title:  "Assignment Graded",
		Message: "Your assignment '" + event.AssignmentTitle +
			"' in class " + event.ClassName +
			" has been graded.",
		Type: "assignment_graded",
	}

	if err := svc.CreateAndSendEmail(req); err != nil {
		log.Println("failed to create notification and send email:", err)
		return
	}

	log.Println("assignment graded email sent to:", event.StudentEmail)
}
