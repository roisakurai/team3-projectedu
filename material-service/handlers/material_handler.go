package handlers

import (
	"phase3/finalproject/material_service-melvin/models"
	"phase3/finalproject/material_service-melvin/repositories"
	"phase3/finalproject/material_service-melvin/services"
	"phase3/finalproject/material_service-melvin/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateMaterial godoc
// @Summary Create material
// @Description Teacher creates learning material
// @Tags Materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param request body models.Material true "Material Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /materials [post]
func CreateMaterial(c *gin.Context) {
	var material models.Material

	if err := c.ShouldBindJSON(&material); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	role := c.MustGet("role").(string)

	if role != "teacher" {
		utils.ErrorResponse(c, http.StatusForbidden, "only teacher can create material")
		return
	}

	material.TeacherID = c.MustGet("user_id").(string)

	err := services.CreateMaterial(material)

	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusCreated,
		"Material created successfully",
		material,
	)
}

// GetAllMaterials godoc
// @Summary Get all materials
// @Description Get all learning materials
// @Tags Materials
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /materials [get]
func GetAllMaterials(c *gin.Context) {
	materials, err := repositories.GetAllMaterials()

	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Success get all materials",
		materials,
	)
}

// GetMaterialsByClass godoc
// @Summary Get materials by class
// @Description Get materials by class ID
// @Tags Materials
// @Produce json
// @Param class_id path string true "Class ID"
// @Success 200 {object} map[string]interface{}
// @Router /materials/class/{class_id} [get]
func GetMaterialsByClass(c *gin.Context) {
	classID := c.Param("class_id")

	materials, err := repositories.GetMaterialsByClass(classID)

	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Success get materials by class",
		materials,
	)
}

// MarkMaterialAsRead godoc
// @Summary Mark material as read
// @Description Student marks material as read
// @Tags Materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Material ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /materials/{id}/read [put]
func MarkMaterialAsRead(c *gin.Context) {
	role := c.MustGet("role").(string)

	if role != "student" {
		utils.ErrorResponse(c, http.StatusForbidden, "only student can read material")
		return
	}

	materialID := c.Param("id")
	studentID := c.MustGet("user_id").(string)

	err := services.MarkMaterialAsRead(materialID, studentID)

	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessResponse(
		c,
		http.StatusOK,
		"Material marked as read",
		gin.H{
			"material_id": materialID,
			"student_id":  studentID,
		},
	)
}
