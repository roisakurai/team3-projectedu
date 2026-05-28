package consumer

import (
	"context"
	"log"

	"progress-service/service"

	"github.com/redis/go-redis/v9"
)

// StartAll launches all 6 Redis Pub/Sub consumers as goroutines.
// Each consumer gets its own goroutine and the given context for shutdown.
func StartAll(ctx context.Context, client *redis.Client, svc service.ProgressService) {
	log.Println("starting all Redis consumers...")

	StartStudentJoinedConsumer(ctx, client, svc)
	StartMaterialCreatedConsumer(ctx, client, svc)
	StartMaterialCompletedConsumer(ctx, client, svc)
	StartAssignmentCreatedConsumer(ctx, client, svc)
	StartAssignmentSubmittedConsumer(ctx, client, svc)
	StartAssignmentGradedConsumer(ctx, client, svc)

	log.Println("all consumers started")
}
