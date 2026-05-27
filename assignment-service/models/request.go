package models

import "time"

// CreateAssignmentRequest is the body for POST /api/v1/assignments.
type CreateAssignmentRequest struct {
	ClassID     string    `json:"class_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FileURL     *string   `json:"file_url"`
	Deadline    time.Time `json:"deadline"`
}

func (r *CreateAssignmentRequest) Validate() error {
	if r.ClassID == "" {
		return ErrValidation("class_id is required")
	}
	if r.Title == "" {
		return ErrValidation("title is required")
	}
	if r.Description == "" {
		return ErrValidation("description is required")
	}
	if r.Deadline.IsZero() {
		return ErrValidation("deadline is required")
	}
	return nil
}

// UpdateAssignmentRequest is the body for PUT /api/v1/assignments/:id.
type UpdateAssignmentRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FileURL     *string   `json:"file_url"`
	Deadline    time.Time `json:"deadline"`
}

func (r *UpdateAssignmentRequest) Validate() error {
	if r.Title == "" {
		return ErrValidation("title is required")
	}
	if r.Description == "" {
		return ErrValidation("description is required")
	}
	if r.Deadline.IsZero() {
		return ErrValidation("deadline is required")
	}
	return nil
}

// SubmitAssignmentRequest is the body for POST /api/v1/assignments/:id/submit.
type SubmitAssignmentRequest struct {
	FileURL *string `json:"file_url"`
	Content *string `json:"content"`
}

func (r *SubmitAssignmentRequest) Validate() error {
	if r.FileURL == nil && r.Content == nil {
		return ErrValidation("at least one of file_url or content is required")
	}
	return nil
}

// GradeSubmissionRequest is the body for POST /api/v1/assignments/:id/grade/:submission_id.
type GradeSubmissionRequest struct {
	Grade float64 `json:"grade"`
}

func (r *GradeSubmissionRequest) Validate() error {
	if r.Grade < 0 || r.Grade > 100 {
		return ErrValidation("grade must be between 0 and 100")
	}
	return nil
}

// ErrValidation is a simple validation error type.
type ErrValidation string

func (e ErrValidation) Error() string {
	return string(e)
}
