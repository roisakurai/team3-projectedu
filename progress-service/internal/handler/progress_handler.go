package handler

import (
	"net/http"

	"progress-service/internal/middleware"
	"progress-service/internal/service"

	"github.com/labstack/echo/v4"
)

type ProgressHandler struct {
	svc service.ProgressService
}

func NewProgressHandler(svc service.ProgressService) *ProgressHandler {
	return &ProgressHandler{svc: svc}
}

func success(c echo.Context, statusCode int, data interface{}) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": true,
		"message": "success",
		"data":    data,
	})
}

func fail(c echo.Context, statusCode int, msg string) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": false,
		"message": msg,
		"data":    nil,
	})
}

// GET /api/v1/progress/me
// Role: student
func (h *ProgressHandler) GetMyProgress(c echo.Context) error {
	studentID := middleware.GetUserID(c)

	progress, err := h.svc.GetMyProgress(c.Request().Context(), studentID)
	if err != nil {
		return fail(c, http.StatusInternalServerError, "failed to get progress: "+err.Error())
	}

	return success(c, http.StatusOK, progress)
}

// GET /api/v1/progress/me/class/:class_id
// Role: student
func (h *ProgressHandler) GetMyClassProgress(c echo.Context) error {
	studentID := middleware.GetUserID(c)
	classID := c.Param("class_id")

	if classID == "" {
		return fail(c, http.StatusBadRequest, "class_id is required")
	}

	detail, err := h.svc.GetMyClassProgress(c.Request().Context(), studentID, classID)
	if err != nil {
		return fail(c, http.StatusInternalServerError, "failed to get class progress: "+err.Error())
	}

	return success(c, http.StatusOK, detail)
}

// GET /api/v1/progress/me/dashboard
// Role: student
func (h *ProgressHandler) GetMyDashboard(c echo.Context) error {
	studentID := middleware.GetUserID(c)

	dashboard, err := h.svc.GetMyDashboard(c.Request().Context(), studentID)
	if err != nil {
		return fail(c, http.StatusInternalServerError, "failed to get dashboard: "+err.Error())
	}

	return success(c, http.StatusOK, dashboard)
}

// GET /api/v1/progress/class/:class_id
// Role: teacher
func (h *ProgressHandler) GetClassProgress(c echo.Context) error {
	classID := c.Param("class_id")

	if classID == "" {
		return fail(c, http.StatusBadRequest, "class_id is required")
	}

	progress, err := h.svc.GetClassProgress(c.Request().Context(), classID)
	if err != nil {
		return fail(c, http.StatusInternalServerError, "failed to get class progress: "+err.Error())
	}

	return success(c, http.StatusOK, progress)
}

// GET /api/v1/progress/class/:class_id/student/:student_id
// Role: teacher
func (h *ProgressHandler) GetStudentClassProgress(c echo.Context) error {
	classID := c.Param("class_id")
	studentID := c.Param("student_id")

	if classID == "" || studentID == "" {
		return fail(c, http.StatusBadRequest, "class_id and student_id are required")
	}

	detail, err := h.svc.GetStudentClassProgress(c.Request().Context(), classID, studentID)
	if err != nil {
		return fail(c, http.StatusInternalServerError, "failed to get student class progress: "+err.Error())
	}

	return success(c, http.StatusOK, detail)
}

// GET /api/v1/progress/summary
// Role: admin
func (h *ProgressHandler) GetPlatformSummary(c echo.Context) error {
	summary, err := h.svc.GetPlatformSummary(c.Request().Context())
	if err != nil {
		return fail(c, http.StatusInternalServerError, "failed to get platform summary: "+err.Error())
	}

	return success(c, http.StatusOK, summary)
}
