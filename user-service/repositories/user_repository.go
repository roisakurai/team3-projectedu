package repository

import (
	"context"
	model "user-service/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(user model.User) error {
	_, err := r.Collection.InsertOne(context.Background(), user)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.Collection.FindOne(context.Background(), bson.M{
		"email": email,
	}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user model.User

	err = r.Collection.FindOne(context.Background(), bson.M{
		"_id": objectID,
	}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
