package main

import (
	"log"
	"os"

	"phase3/finalproject/class_service-melvin/config"
	"phase3/finalproject/class_service-melvin/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "phase3/finalproject/class_service-melvin/docs"
)

// @title Class Service API
// @version 1.0
// @description API for managing classes
// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	config.ConnectDB()

	r := gin.Default()

	routes.SetupRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + os.Getenv("PORT"))
}
