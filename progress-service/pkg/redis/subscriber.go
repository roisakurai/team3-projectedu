package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// NewClient creates and returns a configured Redis client.
func NewClient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// Subscribe subscribes to a Redis Pub/Sub channel and returns the subscription.
// It pings the connection first; if it fails it logs a fatal error.
func Subscribe(ctx context.Context, client *redis.Client, channel string) *redis.PubSub {
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
	sub := client.Subscribe(ctx, channel)
	log.Printf("subscribed to Redis channel: %s", channel)
	return sub
}
