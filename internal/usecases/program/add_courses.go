package program

import (
	"context"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type AddCoursesInput struct {
	ProgramID string   `json:"program_id"`
	CourseIDs []string `json:"course_ids"`
}

type AddCoursesOutput struct {
	Message string `json:"message"`
}

type AddCoursesUseCase interface {
	Execute(ctx context.Context, input AddCoursesInput) (*AddCoursesOutput, error)
}

type addCoursesUseCaseImpl struct {
	repo repointerface.ProgramRepository
}

func NewAddCoursesUseCase(repo repointerface.ProgramRepository) AddCoursesUseCase {
	return &addCoursesUseCaseImpl{repo: repo}
}

func (uc *addCoursesUseCaseImpl) Execute(ctx context.Context, input AddCoursesInput) (*AddCoursesOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	err := uc.repo.AddCourses(ctx, input.ProgramID, input.CourseIDs)
	if err != nil {
		ctxLogger.Errorf("Error adding courses to program: %v", err)
		return nil, err
	}
	return &AddCoursesOutput{Message: "Courses added to program successfully"}, nil
}
