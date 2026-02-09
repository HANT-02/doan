package teacher

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
)

// GetTeacherInput represents the input for getting a teacher
type GetTeacherInput struct {
	ID string `json:"id"`
}

// GetTeacherOutput represents the output after getting a teacher
type GetTeacherOutput struct {
	Teacher *entities.Teacher `json:"teacher"`
}

// GetTeacherUseCase defines the interface for getting a teacher by ID
type GetTeacherUseCase interface {
	Execute(ctx context.Context, input GetTeacherInput) (*GetTeacherOutput, error)
}

type getTeacherUseCase struct {
	teacherRepo repointerface.TeacherRepository
}

// NewGetTeacherUseCase creates a new instance of GetTeacherUseCase
func NewGetTeacherUseCase(teacherRepo repointerface.TeacherRepository) GetTeacherUseCase {
	return &getTeacherUseCase{
		teacherRepo: teacherRepo,
	}
}

func (uc *getTeacherUseCase) Execute(ctx context.Context, input GetTeacherInput) (*GetTeacherOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if input.ID == "" {
		return nil, errors.New("teacher ID is required")
	}

	teacher, err := uc.teacherRepo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to get teacher: %v", err)
		return nil, err
	}

	return &GetTeacherOutput{Teacher: teacher}, nil
}
