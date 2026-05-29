package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"notification-service/config"
	"notification-service/consumers"
	_ "notification-service/docs"
	"notification-service/handlers"
	middlewares "notification-service/middleware"
	"notification-service/repositories"
	"notification-service/services"

	redispkg "notification-service/pkg/redis"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment")
	}

	db := config.ConnectDB()

	notificationRepo := repositories.NewNotificationRepository(db)
	emailService := services.NewEmailService()
	notificationService := services.NewNotificationService(notificationRepo, emailService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	redisClient := redispkg.NewClient()
	defer redisClient.Close()

	consumerCtx := context.Background()
	consumers.StartAll(consumerCtx, redisClient, notificationService)

	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "notification service is running",
		})
	})

	// Swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/notifications", notificationHandler.Create)
	e.POST("/notifications/email", notificationHandler.CreateAndSendEmail)

	notificationGroup := e.Group("/notifications")
	notificationGroup.Use(middlewares.JWTMiddleware(os.Getenv("JWT_SECRET")))

	notificationGroup.GET("/me", notificationHandler.GetMyNotifications)
	notificationGroup.PATCH("/:id/read", notificationHandler.MarkAsRead)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
