package service

import (
	"context"
	"time"

	"progress-service/internal/model"
	"progress-service/internal/repository"
)

// ─── Request/Response types ───────────────────────────────────────────────────

type MaterialProgressSummary struct {
	Total      int     `json:"total"`
	Completed  int     `json:"completed"`
	Percentage float64 `json:"percentage"`
}

type AssignmentProgressSummary struct {
	Total                int     `json:"total"`
	Submitted            int     `json:"submitted"`
	Graded               int     `json:"graded"`
	AverageGrade         float64 `json:"average_grade"`
	SubmissionPercentage float64 `json:"submission_percentage"`
}

type ClassProgressDetail struct {
	ClassID            string                    `json:"class_id"`
	MaterialProgress   MaterialProgressSummary   `json:"material_progress"`
	AssignmentProgress AssignmentProgressSummary `json:"assignment_progress"`
	OverallPercentage  float64                   `json:"overall_percentage"`
}

type DashboardResponse struct {
	StudentID          string                `json:"student_id"`
	Classes            []ClassProgressDetail `json:"classes"`
	GlobalAverageGrade float64               `json:"global_average_grade"`
	LastActivityAt     time.Time             `json:"last_activity_at"`
}

type PlatformSummary struct {
	TotalStudents        int64   `json:"total_students"`
	TotalClasses         int64   `json:"total_classes"`
	PlatformAverageGrade float64 `json:"platform_average_grade"`
	TotalSubmissions     int64   `json:"total_submissions"`
}

// ─── Interface ────────────────────────────────────────────────────────────────

type ProgressService interface {
	GetMyProgress(ctx context.Context, studentID string) ([]*model.StudentProgress, error)
	GetMyClassProgress(ctx context.Context, studentID, classID string) (*ClassProgressDetail, error)
	GetMyDashboard(ctx context.Context, studentID string) (*DashboardResponse, error)
	GetClassProgress(ctx context.Context, classID string) ([]*model.StudentProgress, error)
	GetStudentClassProgress(ctx context.Context, classID, studentID string) (*ClassProgressDetail, error)
	GetPlatformSummary(ctx context.Context) (*PlatformSummary, error)

	// Event handlers
	HandleStudentJoined(ctx context.Context, event *model.StudentJoinedClassEvent) error
	HandleMaterialCreated(ctx context.Context, event *model.MaterialCreatedEvent) error
	HandleMaterialCompleted(ctx context.Context, event *model.MaterialCompletedEvent) error
	HandleAssignmentCreated(ctx context.Context, event *model.AssignmentCreatedEvent) error
	HandleAssignmentSubmitted(ctx context.Context, event *model.AssignmentSubmittedEvent) error
	HandleAssignmentGraded(ctx context.Context, event *model.AssignmentGradedEvent) error
}

// ─── Implementation ───────────────────────────────────────────────────────────

type progressService struct {
	repo repository.ProgressRepository
}

func NewProgressService(repo repository.ProgressRepository) ProgressService {
	return &progressService{repo: repo}
}

func (s *progressService) GetMyProgress(ctx context.Context, studentID string) ([]*model.StudentProgress, error) {
	return s.repo.GetStudentProgressByStudentID(ctx, studentID)
}

func (s *progressService) GetMyClassProgress(ctx context.Context, studentID, classID string) (*ClassProgressDetail, error) {
	return s.buildClassProgressDetail(ctx, studentID, classID)
}

func (s *progressService) GetMyDashboard(ctx context.Context, studentID string) (*DashboardResponse, error) {
	allProgress, err := s.repo.GetStudentProgressByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	dashboard := &DashboardResponse{
		StudentID: studentID,
		Classes:   make([]ClassProgressDetail, 0, len(allProgress)),
	}

	var totalGrade float64
	gradeCount := 0
	var lastActivity time.Time

	for _, sp := range allProgress {
		detail, err := s.buildClassProgressDetail(ctx, studentID, sp.ClassID)
		if err != nil {
			continue
		}
		dashboard.Classes = append(dashboard.Classes, *detail)

		if sp.GradedAssignments > 0 {
			totalGrade += sp.AverageGrade
			gradeCount++
		}

		if sp.LastActivityAt.After(lastActivity) {
			lastActivity = sp.LastActivityAt
		}
	}

	if gradeCount > 0 {
		dashboard.GlobalAverageGrade = totalGrade / float64(gradeCount)
	}
	dashboard.LastActivityAt = lastActivity

	return dashboard, nil
}

func (s *progressService) GetClassProgress(ctx context.Context, classID string) ([]*model.StudentProgress, error) {
	return s.repo.GetStudentProgressByClassID(ctx, classID)
}

func (s *progressService) GetStudentClassProgress(ctx context.Context, classID, studentID string) (*ClassProgressDetail, error) {
	return s.buildClassProgressDetail(ctx, studentID, classID)
}

func (s *progressService) GetPlatformSummary(ctx context.Context) (*PlatformSummary, error) {
	students, err := s.repo.CountDistinctStudents(ctx)
	if err != nil {
		return nil, err
	}

	classes, err := s.repo.CountDistinctClasses(ctx)
	if err != nil {
		return nil, err
	}

	avgGrade, err := s.repo.GetPlatformAverageGrade(ctx)
	if err != nil {
		return nil, err
	}

	totalSubs, err := s.repo.GetTotalSubmissions(ctx)
	if err != nil {
		return nil, err
	}

	return &PlatformSummary{
		TotalStudents:        students,
		TotalClasses:         classes,
		PlatformAverageGrade: avgGrade,
		TotalSubmissions:     totalSubs,
	}, nil
}

// buildClassProgressDetail assembles a ClassProgressDetail for a given student+class.
func (s *progressService) buildClassProgressDetail(ctx context.Context, studentID, classID string) (*ClassProgressDetail, error) {
	sp, err := s.repo.GetStudentProgress(ctx, studentID, classID)
	if err != nil {
		return nil, err
	}

	matPct := 0.0
	if sp.TotalMaterials > 0 {
		matPct = float64(sp.CompletedMaterials) / float64(sp.TotalMaterials) * 100
	}

	subPct := 0.0
	if sp.TotalAssignments > 0 {
		subPct = float64(sp.SubmittedAssignments) / float64(sp.TotalAssignments) * 100
	}

	overallPct := (matPct + subPct) / 2

	return &ClassProgressDetail{
		ClassID: classID,
		MaterialProgress: MaterialProgressSummary{
			Total:      sp.TotalMaterials,
			Completed:  sp.CompletedMaterials,
			Percentage: matPct,
		},
		AssignmentProgress: AssignmentProgressSummary{
			Total:                sp.TotalAssignments,
			Submitted:            sp.SubmittedAssignments,
			Graded:               sp.GradedAssignments,
			AverageGrade:         sp.AverageGrade,
			SubmissionPercentage: subPct,
		},
		OverallPercentage: overallPct,
	}, nil
}

// ─── Event Handlers ───────────────────────────────────────────────────────────

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
