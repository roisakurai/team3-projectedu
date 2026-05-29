package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() *mongo.Database {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerSelectionTimeout(30 * time.Second)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("failed to connect MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("failed to ping MongoDB:", err)
	}

	log.Println("MongoDB connected")
	return client.Database(dbName)
}
