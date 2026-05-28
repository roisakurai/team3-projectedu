package handlers

import (
	"phase3/finalproject/material_service-melvin/models"
	"phase3/finalproject/material_service-melvin/repositories"
	"phase3/finalproject/material_service-melvin/services"
	"phase3/finalproject/material_service-melvin/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

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
