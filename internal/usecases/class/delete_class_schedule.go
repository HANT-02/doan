package class

import (
	"context"
	repointerface "doan/internal/repositories/interface"
	"errors"
)

type DeleteClassScheduleInput struct {
	ClassID    string
	ScheduleID string
}

type DeleteClassScheduleOutput struct {
	Message string `json:"message"`
}

type DeleteClassScheduleUseCase interface {
	Execute(ctx context.Context, input DeleteClassScheduleInput) (*DeleteClassScheduleOutput, error)
}

type deleteClassScheduleUseCase struct {
	classScheduleRepo repointerface.ClassScheduleRepository
}

func NewDeleteClassScheduleUseCase(
	classScheduleRepo repointerface.ClassScheduleRepository,
) DeleteClassScheduleUseCase {
	return &deleteClassScheduleUseCase{
		classScheduleRepo: classScheduleRepo,
	}
}

func (u *deleteClassScheduleUseCase) Execute(ctx context.Context, input DeleteClassScheduleInput) (*DeleteClassScheduleOutput, error) {
	schedule, err := u.classScheduleRepo.GetByID(ctx, input.ScheduleID)
	if err != nil || schedule == nil {
		return nil, errors.New("schedule not found")
	}

	if schedule.ClassID != input.ClassID {
		return nil, errors.New("schedule does not belong to this class")
	}

	err = u.classScheduleRepo.HardDelete(ctx, input.ScheduleID)
	if err != nil {
		return nil, err
	}

	return &DeleteClassScheduleOutput{
		Message: "Schedule deleted successfully",
	}, nil
}
