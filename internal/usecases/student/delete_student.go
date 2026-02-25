package student

import (
	"context"

	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type DeleteStudentUseCase interface {
	Execute(ctx context.Context, id string) error
}

type deleteStudentUseCase struct {
	studentRepo repointerface.StudentRepository
}

func NewDeleteStudentUseCase(studentRepo repointerface.StudentRepository) DeleteStudentUseCase {
	return &deleteStudentUseCase{
		studentRepo: studentRepo,
	}
}

func (uc *deleteStudentUseCase) Execute(ctx context.Context, id string) error {
	ctxLogger := logger.NewLogger(ctx)

	err := uc.studentRepo.SoftDelete(ctx, id)
	if err != nil {
		ctxLogger.Errorf("Failed to delete student: %v", err)
		return err
	}

	return nil
}
