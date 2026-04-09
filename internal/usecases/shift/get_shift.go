package shift

import (
	"context"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetShiftInput struct {
	ID string
}

type GetShiftOutput struct {
	Shift *entities.Shift
}

type GetShiftUseCase interface {
	Execute(ctx context.Context, input GetShiftInput) (*GetShiftOutput, error)
}

type getShiftUseCase struct {
	shiftRepo repointerface.ShiftRepository
}

func NewGetShiftUseCase(shiftRepo repointerface.ShiftRepository) GetShiftUseCase {
	return &getShiftUseCase{shiftRepo: shiftRepo}
}

func (uc *getShiftUseCase) Execute(ctx context.Context, input GetShiftInput) (*GetShiftOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	shiftEntity, err := uc.shiftRepo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to get shift: %v", err)
		return nil, err
	}

	return &GetShiftOutput{Shift: shiftEntity}, nil
}
