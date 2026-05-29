package redis

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewClient() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")

	if redisURL != "" {
		opts, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Fatalf("failed to parse REDIS_URL: %v", err)
		}

		client := redis.NewClient(opts)

		if err := client.Ping(context.Background()).Err(); err != nil {
			log.Fatalf("failed to connect to Redis: %v", err)
		}

		log.Println("connected to Redis using REDIS_URL")
		return client
	}

	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	log.Println("connected to Redis using REDIS_ADDR")
	return client
}
