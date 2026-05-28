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
	ChannelAssignmentCreated   = "lms:assignment:created"
	ChannelAssignmentSubmitted = "lms:assignment:submitted"
	ChannelAssignmentGraded    = "lms:assignment:graded"
)

// StartAssignmentCreatedConsumer subscribes to "lms:assignment:created"
// and processes each event in a goroutine loop.
func StartAssignmentCreatedConsumer(ctx context.Context, client *redis.Client, svc service.ProgressService) {
	sub := redispkg.Subscribe(ctx, client, ChannelAssignmentCreated)
	go func() {
		defer sub.Close()
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				log.Printf("[%s] context cancelled, stopping consumer", ChannelAssignmentCreated)
				return
			case msg, ok := <-ch:
				if !ok {
					log.Printf("[%s] channel closed, stopping consumer", ChannelAssignmentCreated)
					return
				}
				handleAssignmentCreated(ctx, svc, msg)
			}
		}
	}()
}

func handleAssignmentCreated(ctx context.Context, svc service.ProgressService, msg *redis.Message) {
	log.Printf("[%s] received: %s", ChannelAssignmentCreated, msg.Payload)

	var event model.AssignmentCreatedEvent
	if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
		log.Printf("[%s] failed to unmarshal event: %v", ChannelAssignmentCreated, err)
		return
	}

	if err := svc.HandleAssignmentCreated(ctx, &event); err != nil {
		log.Printf("[%s] failed to handle event: %v", ChannelAssignmentCreated, err)
	}
}

// StartAssignmentSubmittedConsumer subscribes to "lms:assignment:submitted".
func StartAssignmentSubmittedConsumer(ctx context.Context, client *redis.Client, svc service.ProgressService) {
	sub := redispkg.Subscribe(ctx, client, ChannelAssignmentSubmitted)
	go func() {
		defer sub.Close()
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				log.Printf("[%s] context cancelled, stopping consumer", ChannelAssignmentSubmitted)
				return
			case msg, ok := <-ch:
				if !ok {
					log.Printf("[%s] channel closed, stopping consumer", ChannelAssignmentSubmitted)
					return
				}
				handleAssignmentSubmitted(ctx, svc, msg)
			}
		}
	}()
}

func handleAssignmentSubmitted(ctx context.Context, svc service.ProgressService, msg *redis.Message) {
	log.Printf("[%s] received: %s", ChannelAssignmentSubmitted, msg.Payload)

	var event model.AssignmentSubmittedEvent
	if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
		log.Printf("[%s] failed to unmarshal event: %v", ChannelAssignmentSubmitted, err)
		return
	}

	if err := svc.HandleAssignmentSubmitted(ctx, &event); err != nil {
		log.Printf("[%s] failed to handle event: %v", ChannelAssignmentSubmitted, err)
	}
}

// StartAssignmentGradedConsumer subscribes to "lms:assignment:graded".
func StartAssignmentGradedConsumer(ctx context.Context, client *redis.Client, svc service.ProgressService) {
	sub := redispkg.Subscribe(ctx, client, ChannelAssignmentGraded)
	go func() {
		defer sub.Close()
		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				log.Printf("[%s] context cancelled, stopping consumer", ChannelAssignmentGraded)
				return
			case msg, ok := <-ch:
				if !ok {
					log.Printf("[%s] channel closed, stopping consumer", ChannelAssignmentGraded)
					return
				}
				handleAssignmentGraded(ctx, svc, msg)
			}
		}
	}()
}

func handleAssignmentGraded(ctx context.Context, svc service.ProgressService, msg *redis.Message) {
	log.Printf("[%s] received: %s", ChannelAssignmentGraded, msg.Payload)

	var event model.AssignmentGradedEvent
	if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
		log.Printf("[%s] failed to unmarshal event: %v", ChannelAssignmentGraded, err)
		return
	}

	if err := svc.HandleAssignmentGraded(ctx, &event); err != nil {
		log.Printf("[%s] failed to handle event: %v", ChannelAssignmentGraded, err)
	}
}
