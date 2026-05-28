package service

import (
	"context"
	"progress-service/model"
)

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
