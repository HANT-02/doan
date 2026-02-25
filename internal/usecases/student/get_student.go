package student

import (
	"context"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetStudentOutput struct {
	Student *entities.Student
}

type GetStudentUseCase interface {
	Execute(ctx context.Context, id string) (*GetStudentOutput, error)
}

type getStudentUseCase struct {
	studentRepo repointerface.StudentRepository
}

func NewGetStudentUseCase(studentRepo repointerface.StudentRepository) GetStudentUseCase {
	return &getStudentUseCase{
		studentRepo: studentRepo,
	}
}

func (uc *getStudentUseCase) Execute(ctx context.Context, id string) (*GetStudentOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	student, err := uc.studentRepo.GetByID(ctx, id)
	if err != nil {
		ctxLogger.Errorf("Failed to get student: %v", err)
		return nil, err
	}

	return &GetStudentOutput{Student: student}, nil
}
