package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": false,
		"message": message,
		"data":    nil,
	})
}

func BadRequest(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusBadRequest, message)
}

func Unauthorized(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusUnauthorized, message)
}

func Forbidden(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusForbidden, message)
}

func NotFound(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusNotFound, message)
}

func InternalServerError(c echo.Context, message string) error {
	return ErrorResponse(c, http.StatusInternalServerError, message)
}
