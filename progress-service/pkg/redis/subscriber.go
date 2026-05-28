package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// NewClient creates and returns a configured Redis client.
// If redisURL is provided it will be parsed and used (supports TLS/username),
// otherwise it falls back to Addr/Username/Password/DB options.
func NewClient(redisURL, addr, username, password string, db int) *redis.Client {
	if redisURL != "" {
		opts, err := redis.ParseURL(redisURL)
		if err == nil {
			if username != "" {
				opts.Username = username
			}
			return redis.NewClient(opts)
		}
		// If parse failed, log and fall back
		log.Printf("warning: failed to parse REDIS_URL, falling back to addr: %v", err)
	}

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
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
