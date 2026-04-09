package shift

import (
	"context"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type UpdateShiftInput struct {
	ID              string
	Code            string
	Name            string
	StartTime       string
	EndTime         string
	DurationMinutes int
	SessionType     string
	IsActive        bool
	Notes           string
}

type UpdateShiftOutput struct {
	Shift *entities.Shift
}

type UpdateShiftUseCase interface {
	Execute(ctx context.Context, input UpdateShiftInput) (*UpdateShiftOutput, error)
}

type updateShiftUseCase struct {
	shiftRepo repointerface.ShiftRepository
}

func NewUpdateShiftUseCase(shiftRepo repointerface.ShiftRepository) UpdateShiftUseCase {
	return &updateShiftUseCase{shiftRepo: shiftRepo}
}

func (uc *updateShiftUseCase) Execute(ctx context.Context, input UpdateShiftInput) (*UpdateShiftOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	shiftEntity, err := uc.shiftRepo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Shift not found: %v", err)
		return nil, err
	}

	updateData := map[string]interface{}{
		"code":             input.Code,
		"name":             input.Name,
		"start_time":       input.StartTime,
		"end_time":         input.EndTime,
		"duration_minutes": input.DurationMinutes,
		"session_type":     input.SessionType,
		"is_active":        input.IsActive,
		"notes":            input.Notes,
	}

	if err := uc.shiftRepo.Update(ctx, input.ID, updateData); err != nil {
		ctxLogger.Errorf("Failed to update shift: %v", err)
		return nil, err
	}

	shiftEntity.Code = input.Code
	shiftEntity.Name = input.Name
	shiftEntity.StartTime = input.StartTime
	shiftEntity.EndTime = input.EndTime
	shiftEntity.DurationMinutes = input.DurationMinutes
	shiftEntity.SessionType = input.SessionType
	shiftEntity.IsActive = input.IsActive
	shiftEntity.Notes = input.Notes

	return &UpdateShiftOutput{Shift: shiftEntity}, nil
}
