package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(
	c *gin.Context,
	statusCode int,
	message string,
	data interface{},
) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(
	c *gin.Context,
	statusCode int,
	message string,
) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
	})
}

func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

func InternalServerError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}

func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}
