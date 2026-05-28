package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GradeRecord struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID    string             `bson:"student_id" json:"student_id"`
	AssignmentID string             `bson:"assignment_id" json:"assignment_id"`
	ClassID      string             `bson:"class_id" json:"class_id"`
	Grade        float64            `bson:"grade" json:"grade"`
	GradedAt     time.Time          `bson:"graded_at" json:"graded_at"`
}
