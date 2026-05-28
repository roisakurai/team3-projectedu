package handlers

import (
	"net/http"

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

// Create godoc
// @Summary Create notification
// @Description Create new notification
// @Tags Notifications
// @Accept json
// @Produce json
// @Param request body models.CreateNotificationRequest true "Notification Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /notifications [post]
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

// GetByUserID godoc
// @Summary Get notifications by user ID
// @Description Retrieve all notifications for a user
// @Tags Notifications
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notifications/user/{user_id} [get]
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

// MarkAsRead godoc
// @Summary Mark notification as read
// @Description Update notification status to read
// @Tags Notifications
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /notifications/{id}/read [patch]
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
