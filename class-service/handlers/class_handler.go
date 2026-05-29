package handlers

import (
	"phase3/finalproject/class_service-melvin/models"
	"phase3/finalproject/class_service-melvin/repositories"
	"phase3/finalproject/class_service-melvin/services"
	"phase3/finalproject/class_service-melvin/utils"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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

	// ensure students slice is initialized and timestamps/count set
	class.Students = []models.Student{}
	class.StudentCount = 0
	class.CreatedAt = time.Now()

	// try to fetch full profile from user-service to sync teacher name
	if rawTokenIfc, exists := c.Get("raw_token"); exists {
		if rawToken, ok := rawTokenIfc.(string); ok && rawToken != "" {
			if profile, err := services.GetUserProfile(rawToken); err == nil {
				if name, ok := profile["name"].(string); ok {
					class.TeacherName = name
				}
			}
		}
	}

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

	rawToken := ""

	if rawTokenIfc, exists := c.Get("raw_token"); exists {
		if token, ok := rawTokenIfc.(string); ok {
			rawToken = token
		}
	}

	err := services.JoinClass(req.JoinCode, studentID, rawToken)

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

func GetStudentsByClass(c *gin.Context) {
	id := c.Param("id")

	class, err := repositories.FindClassByID(id)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	recipients := []gin.H{}

	for _, student := range class.Students {
		recipients = append(recipients, gin.H{
			"user_id": student.StudentID,
			"email":   student.Email,
			"name":    student.Name,
		})
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Success get students by class",
		gin.H{
			"class_id":   class.ID.Hex(),
			"class_name": class.Name,
			"students":   recipients,
		},
	)
}
