package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentProgress struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID            string             `bson:"student_id" json:"student_id"`
	ClassID              string             `bson:"class_id" json:"class_id"`
	TotalMaterials       int                `bson:"total_materials" json:"total_materials"`
	CompletedMaterials   int                `bson:"completed_materials" json:"completed_materials"`
	TotalAssignments     int                `bson:"total_assignments" json:"total_assignments"`
	SubmittedAssignments int                `bson:"submitted_assignments" json:"submitted_assignments"`
	GradedAssignments    int                `bson:"graded_assignments" json:"graded_assignments"`
	AverageGrade         float64            `bson:"average_grade" json:"average_grade"`
	LastActivityAt       time.Time          `bson:"last_activity_at" json:"last_activity_at"`
	UpdatedAt            time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedAt            time.Time          `bson:"created_at" json:"created_at"`
}
