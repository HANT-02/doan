package shift

import (
	"context"

	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type DeleteShiftInput struct {
	ID string
}

type DeleteShiftOutput struct {
	Message string
}

type DeleteShiftUseCase interface {
	Execute(ctx context.Context, input DeleteShiftInput) (*DeleteShiftOutput, error)
}

type deleteShiftUseCase struct {
	shiftRepo repointerface.ShiftRepository
}

func NewDeleteShiftUseCase(shiftRepo repointerface.ShiftRepository) DeleteShiftUseCase {
	return &deleteShiftUseCase{shiftRepo: shiftRepo}
}

func (uc *deleteShiftUseCase) Execute(ctx context.Context, input DeleteShiftInput) (*DeleteShiftOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if err := uc.shiftRepo.HardDelete(ctx, input.ID); err != nil {
		ctxLogger.Errorf("Failed to delete shift: %v", err)
		return nil, err
	}

	return &DeleteShiftOutput{Message: "Xóa ca học thành công"}, nil
}
