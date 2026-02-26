package program

import (
	"context"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type RemoveCoursesInput struct {
	ProgramID string   `json:"program_id"`
	CourseIDs []string `json:"course_ids"`
}

type RemoveCoursesOutput struct {
	Message string `json:"message"`
}

type RemoveCoursesUseCase interface {
	Execute(ctx context.Context, input RemoveCoursesInput) (*RemoveCoursesOutput, error)
}

type removeCoursesUseCaseImpl struct {
	repo repointerface.ProgramRepository
}

func NewRemoveCoursesUseCase(repo repointerface.ProgramRepository) RemoveCoursesUseCase {
	return &removeCoursesUseCaseImpl{repo: repo}
}

func (uc *removeCoursesUseCaseImpl) Execute(ctx context.Context, input RemoveCoursesInput) (*RemoveCoursesOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	err := uc.repo.RemoveCourses(ctx, input.ProgramID, input.CourseIDs)
	if err != nil {
		ctxLogger.Errorf("Error removing courses from program: %v", err)
		return nil, err
	}
	return &RemoveCoursesOutput{Message: "Courses removed from program successfully"}, nil
}
