package services

import (
	"phase3/finalproject/material_service-melvin/models"
	"phase3/finalproject/material_service-melvin/repositories"
	"time"
)

func CreateMaterial(material models.Material) error {
	material.CreatedAt = time.Now()

	return repositories.CreateMaterial(material)
}

func MarkMaterialAsRead(materialID string, studentID string) error {
	return repositories.MarkMaterialAsRead(materialID, studentID)
}
