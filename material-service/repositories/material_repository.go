package repositories

import (
	"context"

	"phase3/finalproject/material_service-melvin/config"
	"phase3/finalproject/material_service-melvin/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetMaterialsByClass(classID string) ([]models.Material, error) {
	collection := config.DB.Collection("materials")

	cursor, err := collection.Find(
		context.Background(),
		bson.M{"class_id": classID},
	)

	if err != nil {
		return nil, err
	}

	var materials []models.Material

	err = cursor.All(context.Background(), &materials)

	return materials, err
}

func MarkMaterialAsRead(materialID string, studentID string) error {
	collection := config.DB.Collection("materials")

	objectID, err := primitive.ObjectIDFromHex(materialID)

	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
		bson.M{
			"$addToSet": bson.M{
				"read_by": studentID,
			},
		},
	)

	return err
}
