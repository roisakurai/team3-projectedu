package repositories

import (
	"context"

	"phase3/finalproject/material_service-melvin/config"
	"phase3/finalproject/material_service-melvin/models"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateMaterial(material models.Material) error {
	collection := config.DB.Collection("materials")

	_, err := collection.InsertOne(context.Background(), material)

	return err
}

func GetAllMaterials() ([]models.Material, error) {
	collection := config.DB.Collection("materials")

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return nil, err
	}

	var materials []models.Material

	err = cursor.All(context.Background(), &materials)

	return materials, err
}

func FindMaterialByJoinCode(code string) (*models.Material, error) {
	collection := config.DB.Collection("materials")

	var material models.Material

	err := collection.FindOne(
		context.Background(),
		bson.M{"join_code": code},
	).Decode(&material)

	return &material, err
}
