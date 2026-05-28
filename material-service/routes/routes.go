package routes

import (
	"phase3/finalproject/material_service-melvin/handlers"
	"phase3/finalproject/material_service-melvin/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	auth := r.Group("/")
	auth.Use(middleware.JWTMiddleware())

	auth.POST("/materials", handlers.CreateMaterial)
	r.GET("/materials/class/:class_id", handlers.GetMaterialsByClass)
	r.GET("/materials", handlers.GetAllMaterials)
	auth.PUT("/materials/:id/read", handlers.MarkMaterialAsRead)
}
