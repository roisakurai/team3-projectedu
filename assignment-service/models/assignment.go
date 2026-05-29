package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Assignment represents a tugas created by a teacher for a class.
type Assignment struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ClassID     string        `bson:"class_id" json:"class_id"`
	TeacherID   string        `bson:"teacher_id" json:"teacher_id"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	FileURL     *string       `bson:"file_url" json:"file_url"`
	Deadline    time.Time     `bson:"deadline" json:"deadline"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
}

// SubmissionStatus enumerates valid submission states.
type SubmissionStatus string

const (
	StatusPending   SubmissionStatus = "pending"
	StatusSubmitted SubmissionStatus = "submitted"
	StatusGraded    SubmissionStatus = "graded"
)

// Submission represents a student's response to an assignment.
type Submission struct {
	ID           bson.ObjectID    `bson:"_id,omitempty" json:"id"`
	AssignmentID string           `bson:"assignment_id" json:"assignment_id"`
	StudentID    string           `bson:"student_id" json:"student_id"`
	FileURL      *string          `bson:"file_url" json:"file_url"`
	Content      *string          `bson:"content" json:"content"`
	Grade        *float64         `bson:"grade" json:"grade"`
	Status       SubmissionStatus `bson:"status" json:"status"`
	SubmittedAt  *time.Time       `bson:"submitted_at" json:"submitted_at"`
	GradedAt     *time.Time       `bson:"graded_at" json:"graded_at"`
	CreatedAt    time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time        `bson:"updated_at" json:"updated_at"`
}
