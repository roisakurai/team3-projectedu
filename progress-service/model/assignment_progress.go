package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AssignmentStatusPending   = "pending"
	AssignmentStatusSubmitted = "submitted"
	AssignmentStatusGraded    = "graded"
)

type AssignmentProgress struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID    string             `bson:"student_id" json:"student_id"`
	AssignmentID string             `bson:"assignment_id" json:"assignment_id"`
	ClassID      string             `bson:"class_id" json:"class_id"`
	Status       string             `bson:"status" json:"status"`
	Grade        *float64           `bson:"grade,omitempty" json:"grade,omitempty"`
	SubmittedAt  *time.Time         `bson:"submitted_at,omitempty" json:"submitted_at,omitempty"`
	GradedAt     *time.Time         `bson:"graded_at,omitempty" json:"graded_at,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
