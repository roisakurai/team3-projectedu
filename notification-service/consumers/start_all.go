package consumers

import (
	"context"
	"log"

	"notification-service/services"

	"github.com/redis/go-redis/v9"
)

func StartAll(ctx context.Context, redisClient *redis.Client, svc *services.NotificationService) {
	log.Println("starting all Redis consumers...")

	StartAssignmentCreatedConsumer(ctx, redisClient, svc)
	StartMaterialCreatedConsumer(ctx, redisClient, svc)
	StartAssignmentGradedConsumer(ctx, redisClient, svc)

	log.Println("all Redis consumers started")
}
