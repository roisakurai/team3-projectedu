package handlers

import (
	"errors"
	"net/http"

	"assignment-service/middleware"
	"assignment-service/models"
	"assignment-service/services"
	"assignment-service/utils"

	"github.com/labstack/echo/v4"
)

// AssignmentHandler wires HTTP routes to service calls.
type AssignmentHandler struct {
	svc services.AssignmentService
}

// NewAssignmentHandler creates a new handler.
func NewAssignmentHandler(svc services.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{svc: svc}
}

// RegisterRoutes attaches all routes to the provided Echo group.
// The group must already have JWTMiddleware applied.
func (h *AssignmentHandler) RegisterRoutes(g *echo.Group) {
	// Both roles.
	g.GET("", h.ListAssignments)
	g.GET("/:id", h.GetAssignment)

	// Teacher only.
	g.POST("", h.CreateAssignment, middleware.RequireRole(middleware.RoleTeacher))
	g.PUT("/:id", h.UpdateAssignment, middleware.RequireRole(middleware.RoleTeacher))
	g.DELETE("/:id", h.DeleteAssignment, middleware.RequireRole(middleware.RoleTeacher))
	g.GET("/:id/submissions", h.ListSubmissions, middleware.RequireRole(middleware.RoleTeacher))
	g.POST("/:id/grade/:submission_id", h.GradeSubmission, middleware.RequireRole(middleware.RoleTeacher))

	// Student only.
	g.POST("/:id/submit", h.SubmitAssignment, middleware.RequireRole(middleware.RoleStudent))
	g.GET("/:id/my-submission", h.GetMySubmission, middleware.RequireRole(middleware.RoleStudent))
}

func userID(c echo.Context) string {
	v, _ := c.Get(middleware.ContextKeyUserID).(string)
	return v
}

func handleServiceError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, services.ErrNotFound):
		return utils.NotFound(c, "resource not found")
	case errors.Is(err, services.ErrForbidden):
		return utils.Forbidden(c, "you are not allowed to perform this action")
	case errors.Is(err, services.ErrDeadlinePassed):
		return utils.ErrorResponse(c, http.StatusUnprocessableEntity, "submission deadline has passed")
	case errors.Is(err, services.ErrAlreadySubmitted):
		return utils.ErrorResponse(c, http.StatusConflict, "you have already submitted this assignment")
	default:
		return utils.InternalServerError(c, "internal server error")
	}
}

// CreateAssignment handles POST /api/v1/assignments
func (h *AssignmentHandler) CreateAssignment(c echo.Context) error {
	var req models.CreateAssignmentRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	assignment, err := h.svc.CreateAssignment(c.Request().Context(), &req, userID(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, http.StatusCreated, "success", assignment)
}

// GetAssignment handles GET /api/v1/assignments/:id
func (h *AssignmentHandler) GetAssignment(c echo.Context) error {
	id := c.Param("id")

	assignment, err := h.svc.GetAssignment(c.Request().Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, http.StatusOK, "success", assignment)
}

// ListAssignments handles GET /api/v1/assignments?class_id=xxx
func (h *AssignmentHandler) ListAssignments(c echo.Context) error {
	classID := c.QueryParam("class_id")
	if classID == "" {
		return utils.BadRequest(c, "class_id query parameter is required")
	}

	assignments, err := h.svc.ListAssignments(c.Request().Context(), classID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, http.StatusOK, "success", assignments)
}

// UpdateAssignment handles PUT /api/v1/assignments/:id
func (h *AssignmentHandler) UpdateAssignment(c echo.Context) error {
	id := c.Param("id")

	var req models.UpdateAssignmentRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	assignment, err := h.svc.UpdateAssignment(c.Request().Context(), id, &req, userID(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, http.StatusOK, "success", assignment)
}

// DeleteAssignment handles DELETE /api/v1/assignments/:id
func (h *AssignmentHandler) DeleteAssignment(c echo.Context) error {
	id := c.Param("id")

	if err := h.svc.DeleteAssignment(c.Request().Context(), id, userID(c)); err != nil {
		return handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, http.StatusOK, "success", nil)
}

// SubmitAssignment handles POST /api/v1/assignments/:id/submit
func (h *AssignmentHandler) SubmitAssignment(c echo.Context) error {
	id := c.Param("id")

	var req models.SubmitAssignmentRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	submission, err := h.svc.SubmitAssignment(c.Request().Context(), id, &req, userID(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, http.StatusCreated, "success", submission)
}

// GetMySubmission handles GET /api/v1/assignments/:id/my-submission
func (h *AssignmentHandler) GetMySubmission(c echo.Context) error {
	id := c.Param("id")

	submission, err := h.svc.GetMySubmission(c.Request().Context(), id, userID(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, http.StatusOK, "success", submission)
}

// ListSubmissions handles GET /api/v1/assignments/:id/submissions
func (h *AssignmentHandler) ListSubmissions(c echo.Context) error {
	id := c.Param("id")

	submissions, err := h.svc.ListSubmissions(c.Request().Context(), id, userID(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, http.StatusOK, "success", submissions)
}

// GradeSubmission handles POST /api/v1/assignments/:id/grade/:submission_id
func (h *AssignmentHandler) GradeSubmission(c echo.Context) error {
	assignmentID := c.Param("id")
	submissionID := c.Param("submission_id")

	var req models.GradeSubmissionRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	if err := req.Validate(); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	submission, err := h.svc.GradeSubmission(c.Request().Context(), assignmentID, submissionID, &req, userID(c))
	if err != nil {
		return handleServiceError(c, err)
	}

	return utils.SuccessResponse(c, http.StatusOK, "success", submission)
}
