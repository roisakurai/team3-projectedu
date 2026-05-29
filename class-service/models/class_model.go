package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	StudentID string    `json:"student_id" bson:"student_id"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	JoinedAt  time.Time `json:"joined_at" bson:"joined_at"`
}

type Class struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	TeacherID    string             `json:"teacher_id" bson:"teacher_id"`
	TeacherName  string             `json:"teacher_name,omitempty" bson:"teacher_name,omitempty"`
	JoinCode     string             `json:"join_code" bson:"join_code"`
	Students     []Student          `json:"students" bson:"students"`
	StudentCount int                `json:"student_count" bson:"student_count"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

type JoinClassRequest struct {
	JoinCode string `json:"join_code"`
}
