package main

import (
	"log"
	"net/http"
	"os"
	"user-service/config"
	handler "user-service/handlers"
	customMiddleware "user-service/middleware"
	repository "user-service/repositories"
	service "user-service/services"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"

	_ "user-service/docs"
)

// @title User Service API
// @version 1.0
// @description API for User Service
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment")
	}

	db := config.ConnectDB()

	userRepo := repository.NewUserRepository(db)
	emailService := service.NewEmailService()
	userService := service.NewUserService(userRepo, emailService)
	userHandler := handler.NewUserHandler(userService)

	e := echo.New()

	e.GET("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "user service is running",
		})
	})

	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)
	e.GET("/verify-email/:token", userHandler.VerifyEmail)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	protected := e.Group("")
	protected.Use(echojwt.WithConfig(customMiddleware.JWTMiddleware()))
	protected.GET("/profile", userHandler.Profile)
	protected.PUT("/users/:id", userHandler.UpdateUser)
	protected.DELETE("/users/:id", userHandler.DeleteUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
