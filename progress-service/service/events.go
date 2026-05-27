package service

import (
	"context"
	"progress-service/model"
	"time"
)

func (s *progressService) HandleStudentJoined(ctx context.Context, event *model.StudentJoinedClassEvent) error {
	_, err := s.repo.GetStudentProgress(ctx, event.StudentID, event.ClassID)
	if err == nil {
		// already exists
		return nil
	}

	sp := &model.StudentProgress{
		StudentID:            event.StudentID,
		ClassID:              event.ClassID,
		TotalMaterials:       0,
		CompletedMaterials:   0,
		TotalAssignments:     0,
		SubmittedAssignments: 0,
		GradedAssignments:    0,
		AverageGrade:         0,
	}
	return s.repo.CreateStudentProgress(ctx, sp)
}

func (s *progressService) HandleMaterialCreated(ctx context.Context, event *model.MaterialCreatedEvent) error {
	allProgress, err := s.repo.GetStudentProgressByClassID(ctx, event.ClassID)
	if err != nil {
		return err
	}

	for _, sp := range allProgress {
		if err := s.repo.IncrementStudentProgressField(ctx, sp.StudentID, event.ClassID, "total_materials", 1); err != nil {
			return err
		}

		mp := &model.MaterialProgress{
			StudentID:   sp.StudentID,
			MaterialID:  event.MaterialID,
			ClassID:     event.ClassID,
			IsCompleted: false,
		}
		// ignore duplicate key errors (idempotency)
		_ = s.repo.CreateMaterialProgress(ctx, mp)
	}
	return nil
}

func (s *progressService) HandleMaterialCompleted(ctx context.Context, event *model.MaterialCompletedEvent) error {
	now := time.Now().UTC()

	if err := s.repo.SetMaterialCompleted(ctx, event.StudentID, event.MaterialID, now); err != nil {
		return err
	}

	if err := s.repo.IncrementStudentProgressField(ctx, event.StudentID, event.ClassID, "completed_materials", 1); err != nil {
		return err
	}

	return s.repo.UpdateStudentLastActivity(ctx, event.StudentID, event.ClassID)
}

func (s *progressService) HandleAssignmentCreated(ctx context.Context, event *model.AssignmentCreatedEvent) error {
	allProgress, err := s.repo.GetStudentProgressByClassID(ctx, event.ClassID)
	if err != nil {
		return err
	}

	for _, sp := range allProgress {
		if err := s.repo.IncrementStudentProgressField(ctx, sp.StudentID, event.ClassID, "total_assignments", 1); err != nil {
			return err
		}

		ap := &model.AssignmentProgress{
			StudentID:    sp.StudentID,
			AssignmentID: event.AssignmentID,
			ClassID:      event.ClassID,
			Status:       model.AssignmentStatusPending,
		}
		// ignore duplicate key errors (idempotency)
		_ = s.repo.CreateAssignmentProgress(ctx, ap)
	}
	return nil
}

func (s *progressService) HandleAssignmentSubmitted(ctx context.Context, event *model.AssignmentSubmittedEvent) error {
	now := time.Now().UTC()

	if err := s.repo.UpsertAssignmentSubmitted(ctx, event.StudentID, event.AssignmentID, event.ClassID, now); err != nil {
		return err
	}

	return s.repo.IncrementStudentProgressField(ctx, event.StudentID, event.ClassID, "submitted_assignments", 1)
}

func (s *progressService) HandleAssignmentGraded(ctx context.Context, event *model.AssignmentGradedEvent) error {
	now := time.Now().UTC()

	if err := s.repo.UpdateAssignmentGraded(ctx, event.StudentID, event.AssignmentID, event.Grade, now); err != nil {
		return err
	}

	gr := &model.GradeRecord{
		StudentID:    event.StudentID,
		AssignmentID: event.AssignmentID,
		ClassID:      event.ClassID,
		Grade:        event.Grade,
		GradedAt:     now,
	}
	if err := s.repo.InsertGradeRecord(ctx, gr); err != nil {
		return err
	}

	if err := s.repo.IncrementStudentProgressField(ctx, event.StudentID, event.ClassID, "graded_assignments", 1); err != nil {
		return err
	}

	// Recalculate average grade across all grade records for this student
	records, err := s.repo.GetGradeRecordsByStudent(ctx, event.StudentID)
	if err != nil {
		return err
	}

	var sum float64
	for _, rec := range records {
		sum += rec.Grade
	}
	avg := sum / float64(len(records))

	if err := s.repo.UpdateStudentAverageGrade(ctx, event.StudentID, event.ClassID, avg); err != nil {
		return err
	}

	return s.repo.UpdateStudentLastActivity(ctx, event.StudentID, event.ClassID)
}
