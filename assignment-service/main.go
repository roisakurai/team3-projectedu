package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"assignment-service/config"
	"assignment-service/handlers"
	appMiddleware "assignment-service/middleware"
	"assignment-service/pkg/redis"
	"assignment-service/repositories"
	"assignment-service/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	// ── MongoDB ──────────────────────────────────────────────────────────────
	mongoCtx, mongoCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoCancel()

	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("error disconnecting MongoDB: %v", err)
		}
	}()

	if err := mongoClient.Ping(mongoCtx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	log.Println("Connected to MongoDB")

	db := mongoClient.Database(cfg.MongoDB)

	// ── Redis ────────────────────────────────────────────────────────────────
	publisher, err := redis.NewPublisher(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	defer func() {
		if err := publisher.Close(); err != nil {
			log.Printf("error closing Redis publisher: %v", err)
		}
	}()
	log.Println("Connected to Redis")

	// ── Dependency wiring ────────────────────────────────────────────────────
	repo, err := repositories.NewMongoRepository(db)
	if err != nil {
		log.Fatalf("failed to initialise repository (index creation): %v", err)
	}

	svc := services.NewAssignmentService(repo, publisher)
	h := handlers.NewAssignmentHandler(svc)

	// ── Echo ─────────────────────────────────────────────────────────────────
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check (no auth)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// All assignment routes require a valid JWT
	api := e.Group("/api/v1/assignments", appMiddleware.JWTMiddleware(cfg.JWTSecret))
	h.RegisterRoutes(api)

	// ── Graceful shutdown ─────────────────────────────────────────────────────
	go func() {
		addr := ":" + cfg.AppPort
		log.Printf("assignment-service listening on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("echo server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
	log.Println("Server exited")
}
