package services

import (
	"context"
	"errors"
	"time"

	"assignment-service/models"
	redispkg "assignment-service/pkg/redis"
	"assignment-service/repositories"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrForbidden        = errors.New("forbidden")
	ErrDeadlinePassed   = errors.New("submission deadline has passed")
	ErrAlreadySubmitted = errors.New("you have already submitted this assignment")
	ErrInvalidGrade     = errors.New("grade must be between 0 and 100")
)

// AssignmentService defines the business-logic interface.
type AssignmentService interface {
	CreateAssignment(ctx context.Context, req *models.CreateAssignmentRequest, teacherID string) (*models.Assignment, error)
	GetAssignment(ctx context.Context, id string) (*models.Assignment, error)
	ListAssignments(ctx context.Context, classID string) ([]*models.Assignment, error)
	UpdateAssignment(ctx context.Context, id string, req *models.UpdateAssignmentRequest, teacherID string) (*models.Assignment, error)
	DeleteAssignment(ctx context.Context, id string, teacherID string) error

	SubmitAssignment(ctx context.Context, assignmentID string, req *models.SubmitAssignmentRequest, studentID string) (*models.Submission, error)
	GetMySubmission(ctx context.Context, assignmentID, studentID string) (*models.Submission, error)
	ListSubmissions(ctx context.Context, assignmentID, teacherID string) ([]*models.Submission, error)
	GradeSubmission(ctx context.Context, assignmentID, submissionID string, req *models.GradeSubmissionRequest, teacherID string) (*models.Submission, error)
}

type assignmentService struct {
	repo      repositories.AssignmentRepository
	publisher *redispkg.Publisher
}

// NewAssignmentService creates a service with the given repository and Redis publisher.
func NewAssignmentService(repo repositories.AssignmentRepository, publisher *redispkg.Publisher) AssignmentService {
	return &assignmentService{
		repo:      repo,
		publisher: publisher,
	}
}

// ── Assignment CRUD ───────────────────────────────────────────────────────────

func (s *assignmentService) CreateAssignment(ctx context.Context, req *models.CreateAssignmentRequest, teacherID string) (*models.Assignment, error) {
	a := &models.Assignment{
		ClassID:     req.ClassID,
		TeacherID:   teacherID,
		Title:       req.Title,
		Description: req.Description,
		FileURL:     req.FileURL,
		Deadline:    req.Deadline.UTC(),
	}

	if err := s.repo.CreateAssignment(ctx, a); err != nil {
		return nil, err
	}

	// Publish event - errors are logged inside publisher, not propagated.
	s.publisher.PublishAssignmentCreated(ctx, redispkg.AssignmentCreatedEvent{
		Event:        "ASSIGNMENT_CREATED",
		AssignmentID: a.ID.Hex(),
		ClassID:      a.ClassID,
		Title:        a.Title,
		Deadline:     a.Deadline.Format(time.RFC3339),
		TeacherID:    teacherID,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	})

	return a, nil
}

func (s *assignmentService) GetAssignment(ctx context.Context, id string) (*models.Assignment, error) {
	a, err := s.repo.GetAssignmentByID(ctx, id)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrNotFound
	}
	return a, err
}

func (s *assignmentService) ListAssignments(ctx context.Context, classID string) ([]*models.Assignment, error) {
	return s.repo.ListAssignmentsByClassID(ctx, classID)
}

func (s *assignmentService) UpdateAssignment(ctx context.Context, id string, req *models.UpdateAssignmentRequest, teacherID string) (*models.Assignment, error) {
	a, err := s.repo.GetAssignmentByID(ctx, id)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if a.TeacherID != teacherID {
		return nil, ErrForbidden
	}

	update := bson.M{
		"title":       req.Title,
		"description": req.Description,
		"file_url":    req.FileURL,
		"deadline":    req.Deadline.UTC(),
	}

	if err := s.repo.UpdateAssignment(ctx, id, update); err != nil {
		return nil, err
	}

	// Return the refreshed document.
	return s.repo.GetAssignmentByID(ctx, id)
}

func (s *assignmentService) DeleteAssignment(ctx context.Context, id string, teacherID string) error {
	a, err := s.repo.GetAssignmentByID(ctx, id)
	if errors.Is(err, repositories.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	if a.TeacherID != teacherID {
		return ErrForbidden
	}

	return s.repo.DeleteAssignment(ctx, id)
}

// ── Submission operations ─────────────────────────────────────────────────────

func (s *assignmentService) SubmitAssignment(ctx context.Context, assignmentID string, req *models.SubmitAssignmentRequest, studentID string) (*models.Submission, error) {
	a, err := s.repo.GetAssignmentByID(ctx, assignmentID)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Enforce deadline.
	if time.Now().UTC().After(a.Deadline) {
		return nil, ErrDeadlinePassed
	}

	// Prevent duplicate submissions.
	existing, err := s.repo.GetSubmissionByAssignmentAndStudent(ctx, assignmentID, studentID)
	if err != nil && !errors.Is(err, repositories.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAlreadySubmitted
	}

	now := time.Now().UTC()
	sub := &models.Submission{
		AssignmentID: assignmentID,
		StudentID:    studentID,
		FileURL:      req.FileURL,
		Content:      req.Content,
		Status:       models.StatusSubmitted,
		SubmittedAt:  &now,
	}

	if err := s.repo.CreateSubmission(ctx, sub); err != nil {
		if errors.Is(err, repositories.ErrAlreadyExists) {
			return nil, ErrAlreadySubmitted
		}
		return nil, err
	}

	// Publish event.
	s.publisher.PublishAssignmentSubmitted(ctx, redispkg.AssignmentSubmittedEvent{
		Event:        "ASSIGNMENT_SUBMITTED",
		SubmissionID: sub.ID.Hex(),
		AssignmentID: assignmentID,
		StudentID:    studentID,
		ClassID:      a.ClassID,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	})

	return sub, nil
}

func (s *assignmentService) GetMySubmission(ctx context.Context, assignmentID, studentID string) (*models.Submission, error) {
	// Verify assignment exists.
	if _, err := s.repo.GetAssignmentByID(ctx, assignmentID); err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	sub, err := s.repo.GetSubmissionByAssignmentAndStudent(ctx, assignmentID, studentID)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrNotFound
	}
	return sub, err
}

func (s *assignmentService) ListSubmissions(ctx context.Context, assignmentID, teacherID string) ([]*models.Submission, error) {
	a, err := s.repo.GetAssignmentByID(ctx, assignmentID)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if a.TeacherID != teacherID {
		return nil, ErrForbidden
	}

	return s.repo.ListSubmissionsByAssignment(ctx, assignmentID)
}

func (s *assignmentService) GradeSubmission(ctx context.Context, assignmentID, submissionID string, req *models.GradeSubmissionRequest, teacherID string) (*models.Submission, error) {
	// Verify assignment exists and belongs to teacher.
	a, err := s.repo.GetAssignmentByID(ctx, assignmentID)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if a.TeacherID != teacherID {
		return nil, ErrForbidden
	}

	// Fetch the submission.
	sub, err := s.repo.GetSubmissionByID(ctx, submissionID)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Ensure submission belongs to this assignment.
	if sub.AssignmentID != assignmentID {
		return nil, ErrNotFound
	}

	now := time.Now().UTC()
	update := bson.M{
		"grade":     req.Grade,
		"status":    models.StatusGraded,
		"graded_at": now,
	}

	if err := s.repo.UpdateSubmission(ctx, submissionID, update); err != nil {
		return nil, err
	}

	// Return updated submission.
	updated, err := s.repo.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}

	// Publish event.
	s.publisher.PublishAssignmentGraded(ctx, redispkg.AssignmentGradedEvent{
		Event:        "ASSIGNMENT_GRADED",
		SubmissionID: submissionID,
		AssignmentID: assignmentID,
		StudentID:    sub.StudentID,
		Grade:        req.Grade,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	})

	return updated, nil
}
