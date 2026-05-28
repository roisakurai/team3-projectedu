package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MaterialProgress struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID   string             `bson:"student_id" json:"student_id"`
	MaterialID  string             `bson:"material_id" json:"material_id"`
	ClassID     string             `bson:"class_id" json:"class_id"`
	IsCompleted bool               `bson:"is_completed" json:"is_completed"`
	CompletedAt *time.Time         `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}
