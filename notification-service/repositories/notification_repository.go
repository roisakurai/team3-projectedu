package repositories

import (
	"context"

	"notification-service/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type NotificationRepository struct {
	Collection *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	return &NotificationRepository{
		Collection: db.Collection("notifications"),
	}
}

func (r *NotificationRepository) CreateMany(notifications []interface{}) error {
	_, err := r.Collection.InsertMany(context.Background(), notifications)
	return err
}

func (r *NotificationRepository) FindByUserID(userID string) ([]models.Notification, error) {
	filter := bson.M{
		"user_id": userID,
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.Collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var notifications []models.Notification

	if err := cursor.All(context.Background(), &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.Collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"is_read": true,
			},
		},
	)

	return err
}
