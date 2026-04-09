package shift

import (
	"context"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type CreateShiftInput struct {
	Code            string
	Name            string
	StartTime       string
	EndTime         string
	DurationMinutes int
	SessionType     string
	IsActive        bool
	Notes           string
}

type CreateShiftOutput struct {
	Shift *entities.Shift
}

type CreateShiftUseCase interface {
	Execute(ctx context.Context, input CreateShiftInput) (*CreateShiftOutput, error)
}

type createShiftUseCase struct {
	shiftRepo repointerface.ShiftRepository
}

func NewCreateShiftUseCase(shiftRepo repointerface.ShiftRepository) CreateShiftUseCase {
	return &createShiftUseCase{shiftRepo: shiftRepo}
}

func (uc *createShiftUseCase) Execute(ctx context.Context, input CreateShiftInput) (*CreateShiftOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	createdShift, err := uc.shiftRepo.Create(ctx, &entities.Shift{
		Code:            input.Code,
		Name:            input.Name,
		StartTime:       input.StartTime,
		EndTime:         input.EndTime,
		DurationMinutes: input.DurationMinutes,
		SessionType:     input.SessionType,
		IsActive:        input.IsActive,
		Notes:           input.Notes,
	})
	if err != nil {
		ctxLogger.Errorf("Failed to create shift: %v", err)
		return nil, err
	}

	return &CreateShiftOutput{Shift: createdShift}, nil
}
