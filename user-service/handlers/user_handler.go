package handler

import (
	"encoding/json"
	"net/http"
	model "user-service/models"
	service "user-service/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) Register(c echo.Context) error {
	var req model.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	err := h.Service.Register(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "user registered successfully",
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	var req model.LoginRequest

	decoder := json.NewDecoder(c.Request().Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "email and password are required",
		})
	}

	token, err := h.Service.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "login success",
		"token":   token,
	})
}

func (h *UserHandler) Profile(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(string)

	user, err := h.Service.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "user not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "profile retrieved successfully",
		"data":    user,
	})
}

func (h *UserHandler) VerifyEmail(c echo.Context) error {
	token := c.Param("token")

	err := h.Service.VerifyEmail(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "email verified successfully",
	})
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	targetUserID := c.Param("id")

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	requesterID := claims["user_id"].(string)
	requesterRole := claims["role"].(string)

	var req model.UpdateUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	err := h.Service.UpdateUser(targetUserID, requesterID, requesterRole, req)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	targetUserID := c.Param("id")

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	requesterID := claims["user_id"].(string)
	requesterRole := claims["role"].(string)

	err := h.Service.DeleteUser(targetUserID, requesterID, requesterRole)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user deleted successfully",
	})
}
