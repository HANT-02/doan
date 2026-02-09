package teacher

import (
	"context"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
)

// DeleteTeacherInput represents the input for deleting a teacher
type DeleteTeacherInput struct {
	ID string `json:"id"`
}

// DeleteTeacherOutput represents the output after deleting a teacher
type DeleteTeacherOutput struct {
	Message string `json:"message"`
}

// DeleteTeacherUseCase defines the interface for deleting a teacher
type DeleteTeacherUseCase interface {
	Execute(ctx context.Context, input DeleteTeacherInput) (*DeleteTeacherOutput, error)
}

type deleteTeacherUseCase struct {
	teacherRepo repointerface.TeacherRepository
}

// NewDeleteTeacherUseCase creates a new instance of DeleteTeacherUseCase
func NewDeleteTeacherUseCase(teacherRepo repointerface.TeacherRepository) DeleteTeacherUseCase {
	return &deleteTeacherUseCase{
		teacherRepo: teacherRepo,
	}
}

func (uc *deleteTeacherUseCase) Execute(ctx context.Context, input DeleteTeacherInput) (*DeleteTeacherOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if input.ID == "" {
		return nil, errors.New("teacher ID is required")
	}

	// Check if teacher exists
	_, err := uc.teacherRepo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to get teacher: %v", err)
		return nil, err
	}

	// Soft delete teacher
	err = uc.teacherRepo.SoftDelete(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to delete teacher: %v", err)
		return nil, err
	}

	return &DeleteTeacherOutput{Message: "Teacher deleted successfully"}, nil
}
