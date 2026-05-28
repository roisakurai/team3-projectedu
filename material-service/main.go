package main

import (
	"log"
	"os"

	"phase3/finalproject/material_service-melvin/config"
	"phase3/finalproject/material_service-melvin/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	config.ConnectDB()

	r := gin.Default()

	routes.SetupRoutes(r)

	r.Run(":" + os.Getenv("PORT"))
}
