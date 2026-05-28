package routes

import (
	"phase3/finalproject/class_service-melvin/handlers"
	"phase3/finalproject/class_service-melvin/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	auth := r.Group("/")
	auth.Use(middleware.JWTMiddleware())

	auth.POST("/classes", handlers.CreateClass)
	auth.GET("/classes", handlers.GetAllClasses)
	auth.POST("/classes/join", handlers.JoinClass)

	auth.GET("/classes/:id/students", handlers.GetStudentsByClass)
}
