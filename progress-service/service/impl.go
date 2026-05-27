package service

import (
	"context"
	"progress-service/model"
	"progress-service/repository"
)

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

func (s *progressService) GetClassProgress(ctx context.Context, classID string) ([]*model.StudentProgress, error) {
	return s.repo.GetStudentProgressByClassID(ctx, classID)
}

func (s *progressService) GetStudentClassProgress(ctx context.Context, classID, studentID string) (*ClassProgressDetail, error) {
	return s.buildClassProgressDetail(ctx, studentID, classID)
}
