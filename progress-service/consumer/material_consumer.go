package consumer

import (
	"context"
	"encoding/json"
	"log"

	"progress-service/model"
	"progress-service/service"
	redispkg "progress-service/pkg/redis"

	"github.com/redis/go-redis/v9"
)

const (
	ChannelMaterialCreated   = "lms:material:created"
	ChannelMaterialCompleted = "lms:material:completed"
	ChannelStudentJoined     = "lms:class:student_joined"
)

// StartMaterialCreatedConsumer subscribes to "lms:material:created".
func StartMaterialCreatedConsumer(ctx context.Context, client *redis.Client, svc service.ProgressService) {
	sub := redispkg.Subscribe(ctx, client, ChannelMaterialCreated)
	go func() {
		defer sub.Close()
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				log.Printf("[%s] context cancelled, stopping consumer", ChannelMaterialCreated)
				return
			case msg, ok := <-ch:
				if !ok {
					log.Printf("[%s] channel closed, stopping consumer", ChannelMaterialCreated)
					return
				}
				handleMaterialCreated(ctx, svc, msg)
			}
		}
	}()
}

func handleMaterialCreated(ctx context.Context, svc service.ProgressService, msg *redis.Message) {
	log.Printf("[%s] received: %s", ChannelMaterialCreated, msg.Payload)

	var event model.MaterialCreatedEvent
	if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
		log.Printf("[%s] failed to unmarshal event: %v", ChannelMaterialCreated, err)
		return
	}

	if err := svc.HandleMaterialCreated(ctx, &event); err != nil {
		log.Printf("[%s] failed to handle event: %v", ChannelMaterialCreated, err)
	}
}

// StartMaterialCompletedConsumer subscribes to "lms:material:completed".
func StartMaterialCompletedConsumer(ctx context.Context, client *redis.Client, svc service.ProgressService) {
	sub := redispkg.Subscribe(ctx, client, ChannelMaterialCompleted)
	go func() {
		defer sub.Close()
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				log.Printf("[%s] context cancelled, stopping consumer", ChannelMaterialCompleted)
				return
			case msg, ok := <-ch:
				if !ok {
					log.Printf("[%s] channel closed, stopping consumer", ChannelMaterialCompleted)
					return
				}
				handleMaterialCompleted(ctx, svc, msg)
			}
		}
	}()
}

func handleMaterialCompleted(ctx context.Context, svc service.ProgressService, msg *redis.Message) {
	log.Printf("[%s] received: %s", ChannelMaterialCompleted, msg.Payload)

	var event model.MaterialCompletedEvent
	if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
		log.Printf("[%s] failed to unmarshal event: %v", ChannelMaterialCompleted, err)
		return
	}

	if err := svc.HandleMaterialCompleted(ctx, &event); err != nil {
		log.Printf("[%s] failed to handle event: %v", ChannelMaterialCompleted, err)
	}
}

// StartStudentJoinedConsumer subscribes to "lms:class:student_joined".
func StartStudentJoinedConsumer(ctx context.Context, client *redis.Client, svc service.ProgressService) {
	sub := redispkg.Subscribe(ctx, client, ChannelStudentJoined)
	go func() {
		defer sub.Close()
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				log.Printf("[%s] context cancelled, stopping consumer", ChannelStudentJoined)
				return
			case msg, ok := <-ch:
				if !ok {
					log.Printf("[%s] channel closed, stopping consumer", ChannelStudentJoined)
					return
				}
				handleStudentJoined(ctx, svc, msg)
			}
		}
	}()
}

func handleStudentJoined(ctx context.Context, svc service.ProgressService, msg *redis.Message) {
	log.Printf("[%s] received: %s", ChannelStudentJoined, msg.Payload)

	var event model.StudentJoinedClassEvent
	if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
		log.Printf("[%s] failed to unmarshal event: %v", ChannelStudentJoined, err)
		return
	}

	if err := svc.HandleStudentJoined(ctx, &event); err != nil {
		log.Printf("[%s] failed to handle event: %v", ChannelStudentJoined, err)
	}
}
