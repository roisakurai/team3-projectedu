package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"progress-service/config"
	"progress-service/consumer"
	_ "progress-service/docs"
	"progress-service/handler"
	appMiddleware "progress-service/middleware"
	redispkg "progress-service/pkg/redis"
	"progress-service/repository"
	"progress-service/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	// ─── MongoDB ──────────────────────────────────────────────────────────────
	mongoCtx, mongoCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoCancel()

	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("mongodb disconnect error: %v", err)
		}
	}()

	if err := mongoClient.Ping(mongoCtx, nil); err != nil {
		log.Fatalf("mongodb ping failed: %v", err)
	}
	log.Println("connected to MongoDB")

	db := mongoClient.Database(cfg.MongoDB)

	// Create indexes
	idxCtx, idxCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer idxCancel()
	if err := repository.EnsureIndexes(idxCtx, db); err != nil {
		log.Fatalf("failed to create indexes: %v", err)
	}
	log.Println("MongoDB indexes ensured")

	// ─── Redis ────────────────────────────────────────────────────────────────
	redisClient := redispkg.NewClient(cfg.RedisURL, cfg.RedisAddr, cfg.RedisUsername, cfg.RedisPassword, cfg.RedisDB)
	defer redisClient.Close()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := redisClient.Ping(pingCtx).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	log.Println("connected to Redis")

	// ─── Dependency injection ─────────────────────────────────────────────────
	repo := repository.NewProgressRepository(db)
	svc := service.NewProgressService(repo)
	h := handler.NewProgressHandler(svc)

	// ─── Start consumers ──────────────────────────────────────────────────────
	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()
	consumer.StartAll(consumerCtx, redisClient, svc)

	// ─── Echo HTTP server ─────────────────────────────────────────────────────
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// API routes
	api := e.Group("/api/v1")
	progressGroup := api.Group("/progress", appMiddleware.JWTMiddleware(cfg.JWTSecret))

	// Student routes
	student := progressGroup.Group("", appMiddleware.RequireRole(appMiddleware.RoleStudent))
	student.GET("/me", h.GetMyProgress)
	student.GET("/me/class/:class_id", h.GetMyClassProgress)
	student.GET("/me/dashboard", h.GetMyDashboard)

	// Teacher routes
	teacher := progressGroup.Group("", appMiddleware.RequireRole(appMiddleware.RoleTeacher))
	teacher.GET("/class/:class_id", h.GetClassProgress)
	teacher.GET("/class/:class_id/student/:student_id", h.GetStudentClassProgress)

	// Admin routes
	admin := progressGroup.Group("", appMiddleware.RequireRole(appMiddleware.RoleAdmin))
	admin.GET("/summary", h.GetPlatformSummary)

	// ─── Graceful shutdown ────────────────────────────────────────────────────
	go func() {
		addr := ":" + cfg.AppPort
		log.Printf("starting progress-service on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	consumerCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("server stopped")
}
