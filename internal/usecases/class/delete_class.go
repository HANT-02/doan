package class

import (
	"context"

	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type DeleteClassInput struct {
	ID string
}

type DeleteClassOutput struct {
	Message string
}

type DeleteClassUseCase interface {
	Execute(ctx context.Context, input DeleteClassInput) (*DeleteClassOutput, error)
}

type deleteClassUseCase struct {
	classRepo repointerface.ClassRepository
}

func NewDeleteClassUseCase(classRepo repointerface.ClassRepository) DeleteClassUseCase {
	return &deleteClassUseCase{
		classRepo: classRepo,
	}
}

func (uc *deleteClassUseCase) Execute(ctx context.Context, input DeleteClassInput) (*DeleteClassOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	err := uc.classRepo.SoftDelete(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to delete class: %v", err)
		return nil, err
	}

	return &DeleteClassOutput{Message: "Class deleted successfully"}, nil
}
