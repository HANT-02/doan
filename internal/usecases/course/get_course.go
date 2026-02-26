package course

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetCourseInput struct {
	ID string
}

type GetCourseOutput struct {
	Course *entities.Course `json:"course"`
}

type GetCourseUseCase interface {
	Execute(ctx context.Context, input GetCourseInput) (*GetCourseOutput, error)
}

type getCourseUseCaseImpl struct {
	repo repointerface.CourseRepository
}

func NewGetCourseUseCase(repo repointerface.CourseRepository) GetCourseUseCase {
	return &getCourseUseCaseImpl{repo: repo}
}

func (uc *getCourseUseCaseImpl) Execute(ctx context.Context, input GetCourseInput) (*GetCourseOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	course, err := uc.repo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Error getting course by ID: %v", err)
		return nil, err
	}
	return &GetCourseOutput{Course: course}, nil
}
