package handlers

import (
	"phase3/finalproject/class_service-melvin/models"
	"phase3/finalproject/class_service-melvin/repositories"
	"phase3/finalproject/class_service-melvin/services"
	"phase3/finalproject/class_service-melvin/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateClass godoc
// @Summary Create class
// @Description Teacher creates class
// @Tags Classes
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body models.Class true "Class Data"
// @Success 201 {object} map[string]interface{}
// @Router /classes [post]
func CreateClass(c *gin.Context) {
	var class models.Class

	if err := c.ShouldBindJSON(&class); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	role := c.MustGet("role").(string)

	if role != "teacher" {
		utils.ErrorResponse(c, http.StatusForbidden, "only teacher can create class")
		return
	}

	class.TeacherID = c.MustGet("user_id").(string)

	err := services.CreateClass(class)

	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusCreated,
		"Class created successfully",
		class,
	)
}

// GetAllClasses godoc
// @Summary Get all classes
// @Tags Classes
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Router /classes [get]
func GetAllClasses(c *gin.Context) {
	classes, err := repositories.GetAllClasses()

	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Success get all classes",
		classes,
	)
}

// JoinClass godoc
// @Summary Student join class
// @Tags Classes
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body models.JoinClassRequest true "Join Class"
// @Success 200 {object} map[string]interface{}
// @Router /classes/join [post]
func JoinClass(c *gin.Context) {
	var req models.JoinClassRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	role := c.MustGet("role").(string)

	if role != "student" {
		utils.ErrorResponse(c, http.StatusForbidden, "only student can join class")
		return
	}

	studentID := c.MustGet("user_id").(string)

	err := services.JoinClass(req.JoinCode, studentID)

	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Joined class successfully",
		gin.H{
			"student_id": studentID,
			"name":       "Melvin Student",
			"join_code":  req.JoinCode,
		},
	)
}

// GetStudentsByClass godoc
// @Summary Get students by class
// @Tags Classes
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Class ID"
// @Success 200 {object} map[string]interface{}
// @Router /classes/{id}/students [get]
func GetStudentsByClass(c *gin.Context) {
	id := c.Param("id")

	class, err := repositories.FindClassByID(id)

	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Success get students by class",
		class.Students,
	)
}
