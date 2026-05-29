package handlers

import (
	"net/http"

	middlewares "notification-service/middleware"
	"notification-service/models"
	"notification-service/services"

	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	Service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{Service: service}
}

func (h *NotificationHandler) Create(c echo.Context) error {
	var req models.CreateNotificationRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	if err := h.Service.Create(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "notification created successfully",
	})
}

func (h *NotificationHandler) GetByUserID(c echo.Context) error {
	userID := c.Param("user_id")

	notifications, err := h.Service.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to get notifications",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "notifications retrieved successfully",
		"data":    notifications,
	})
}

func (h *NotificationHandler) MarkAsRead(c echo.Context) error {
	id := c.Param("id")

	if err := h.Service.MarkAsRead(id); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "failed to mark notification as read",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "notification marked as read",
	})
}

func (h *NotificationHandler) CreateAndSendEmail(c echo.Context) error {
	var req models.CreateEmailNotificationRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	if err := h.Service.CreateAndSendEmail(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "notification created and email sent",
	})
}

func (h *NotificationHandler) GetMyNotifications(c echo.Context) error {
	userID := middlewares.GetUserID(c)

	notifications, err := h.Service.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, notifications)
}
