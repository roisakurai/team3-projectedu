package service

import "time"

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
