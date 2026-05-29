package main

import (
	"log"
	"os"

	"phase3/finalproject/material_service-melvin/config"
	_ "phase3/finalproject/material_service-melvin/docs"
	"phase3/finalproject/material_service-melvin/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	config.ConnectDB()

	r := gin.Default()

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRoutes(r)

	r.Run(":" + os.Getenv("PORT"))
}
