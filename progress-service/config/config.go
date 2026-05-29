package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	MongoURI      string
	MongoDB       string
	JWTSecret     string
	RedisAddr     string
	RedisPassword string
	RedisURL      string
	RedisUsername string
	RedisDB       int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, loading from environment")
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		redisDB = 0
	}

	return &Config{
		AppPort:       getEnv("APP_PORT", "8082"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:       getEnv("MONGO_DB", "progress_db"),
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisURL:      getEnv("REDIS_URL", ""),
		RedisUsername: getEnv("REDIS_USERNAME", ""),
		RedisDB:       redisDB,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
