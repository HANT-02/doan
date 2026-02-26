package course

import (
	"context"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type DeleteCourseInput struct {
	ID string
}

type DeleteCourseOutput struct {
	Message string `json:"message"`
}

type DeleteCourseUseCase interface {
	Execute(ctx context.Context, input DeleteCourseInput) (*DeleteCourseOutput, error)
}

type deleteCourseUseCaseImpl struct {
	repo repointerface.CourseRepository
}

func NewDeleteCourseUseCase(repo repointerface.CourseRepository) DeleteCourseUseCase {
	return &deleteCourseUseCaseImpl{repo: repo}
}

func (uc *deleteCourseUseCaseImpl) Execute(ctx context.Context, input DeleteCourseInput) (*DeleteCourseOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	err := uc.repo.SoftDelete(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Error soft deleting course: %v", err)
		return nil, err
	}
	return &DeleteCourseOutput{Message: "Course deleted successfully"}, nil
}
