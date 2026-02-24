package class

import (
	"context"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetClassInput struct {
	ID string
}

type GetClassOutput struct {
	Class *entities.Class
}

type GetClassUseCase interface {
	Execute(ctx context.Context, input GetClassInput) (*GetClassOutput, error)
}

type getClassUseCase struct {
	classRepo repointerface.ClassRepository
}

func NewGetClassUseCase(classRepo repointerface.ClassRepository) GetClassUseCase {
	return &getClassUseCase{
		classRepo: classRepo,
	}
}

func (uc *getClassUseCase) Execute(ctx context.Context, input GetClassInput) (*GetClassOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	classEntity, err := uc.classRepo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to get class: %v", err)
		return nil, err
	}

	return &GetClassOutput{Class: classEntity}, nil
}
