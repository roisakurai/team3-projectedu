package repositories

import (
	"context"

	"phase3/finalproject/class_service-melvin/config"
	"phase3/finalproject/class_service-melvin/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateClass(class models.Class) error {
	collection := config.DB.Collection("classes")

	_, err := collection.InsertOne(context.Background(), class)

	return err
}

func GetAllClasses() ([]models.Class, error) {
	collection := config.DB.Collection("classes")

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return nil, err
	}

	var classes []models.Class

	err = cursor.All(context.Background(), &classes)

	return classes, err
}

func FindClassByJoinCode(code string) (*models.Class, error) {
	collection := config.DB.Collection("classes")

	var class models.Class

	err := collection.FindOne(
		context.Background(),
		bson.M{"join_code": code},
	).Decode(&class)

	return &class, err
}

func FindClassByID(id string) (*models.Class, error) {
	collection := config.DB.Collection("classes")

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	var class models.Class

	err = collection.FindOne(
		context.Background(),
		bson.M{"_id": objectID},
	).Decode(&class)

	return &class, err
}

func JoinClass(classID string, student models.Student) error {
	collection := config.DB.Collection("classes")

	objectID, err := primitive.ObjectIDFromHex(classID)

	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{
			"_id": objectID,
		},
		bson.M{
			"$push": bson.M{
				"students": student,
			},
		},
	)

	return err
}
