package service

import (
	"context"
	"time"
)

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
