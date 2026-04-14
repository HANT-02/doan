package class

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"errors"
)

type GetClassSchedulesInput struct {
	ClassID string
}

type GetClassSchedulesOutput struct {
	Schedules []entities.ClassSchedule `json:"schedules"`
}

type GetClassSchedulesUseCase interface {
	Execute(ctx context.Context, input GetClassSchedulesInput) (*GetClassSchedulesOutput, error)
}

type getClassSchedulesUseCase struct {
	classRepo         repointerface.ClassRepository
	classScheduleRepo repointerface.ClassScheduleRepository
}

func NewGetClassSchedulesUseCase(
	classRepo repointerface.ClassRepository,
	classScheduleRepo repointerface.ClassScheduleRepository,
) GetClassSchedulesUseCase {
	return &getClassSchedulesUseCase{
		classRepo:         classRepo,
		classScheduleRepo: classScheduleRepo,
	}
}

func (u *getClassSchedulesUseCase) Execute(ctx context.Context, input GetClassSchedulesInput) (*GetClassSchedulesOutput, error) {
	classEntity, err := u.classRepo.GetByID(ctx, input.ClassID)
	if err != nil || classEntity == nil {
		return nil, errors.New("class not found")
	}

	schedules, err := u.classScheduleRepo.GetSchedulesByClassID(ctx, input.ClassID)
	if err != nil {
		return nil, err
	}

	return &GetClassSchedulesOutput{
		Schedules: schedules,
	}, nil
}
