package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Material struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Content   string             `json:"content" bson:"content"`
	ClassID   string             `json:"class_id" bson:"class_id"`
	TeacherID string             `json:"teacher_id" bson:"teacher_id"`
	ReadBy    []string           `json:"read_by" bson:"read_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
