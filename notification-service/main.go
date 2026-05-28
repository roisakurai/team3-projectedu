package main

import (
	"log"
	"os"

	"notification-service/config"
	"notification-service/handlers"
	"notification-service/repositories"
	"notification-service/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	_ "notification-service/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Notification Service API
// @version 1.0
// @description API for managing notifications
// @host localhost:8086
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment")
	}

	db := config.ConnectDB()

	notificationRepo := repositories.NewNotificationRepository(db)
	notificationService := services.NewNotificationService(notificationRepo)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, echo.Map{
			"message": "notification service is running",
		})
	})

	e.POST("/notifications", notificationHandler.Create)
	e.GET("/notifications/user/:user_id", notificationHandler.GetByUserID)
	e.PATCH("/notifications/:id/read", notificationHandler.MarkAsRead)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
