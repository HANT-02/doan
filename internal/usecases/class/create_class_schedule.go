package class

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"errors"
)

type CreateClassScheduleInput struct {
	ClassID   string
	ShiftID   string
	DayOfWeek string
	RoomID    *string
}

type CreateClassScheduleOutput struct {
	Schedule *entities.ClassSchedule `json:"schedule"`
}

type CreateClassScheduleUseCase interface {
	Execute(ctx context.Context, input CreateClassScheduleInput) (*CreateClassScheduleOutput, error)
}

type createClassScheduleUseCase struct {
	classRepo         repointerface.ClassRepository
	shiftRepo         repointerface.ShiftRepository
	roomRepo          repointerface.RoomRepository
	classScheduleRepo repointerface.ClassScheduleRepository
}

func NewCreateClassScheduleUseCase(
	classRepo repointerface.ClassRepository,
	shiftRepo repointerface.ShiftRepository,
	roomRepo repointerface.RoomRepository,
	classScheduleRepo repointerface.ClassScheduleRepository,
) CreateClassScheduleUseCase {
	return &createClassScheduleUseCase{
		classRepo:         classRepo,
		shiftRepo:         shiftRepo,
		roomRepo:          roomRepo,
		classScheduleRepo: classScheduleRepo,
	}
}

func (u *createClassScheduleUseCase) Execute(ctx context.Context, input CreateClassScheduleInput) (*CreateClassScheduleOutput, error) {
	// 1. Verify Class
	classEntity, err := u.classRepo.GetByID(ctx, input.ClassID)
	if err != nil || classEntity == nil {
		return nil, errors.New("class not found")
	}

	// 2. Verify Shift
	shiftEntity, err := u.shiftRepo.GetByID(ctx, input.ShiftID)
	if err != nil || shiftEntity == nil {
		return nil, errors.New("shift not found")
	}

	// 3. Verify Room if provided
	if input.RoomID != nil && *input.RoomID != "" {
		roomEntity, err := u.roomRepo.GetByID(ctx, *input.RoomID)
		if err != nil || roomEntity == nil {
			return nil, errors.New("room not found")
		}
	} else {
		input.RoomID = nil
	}

	// 4. Create Schedule
	schedule := &entities.ClassSchedule{
		ClassID:   input.ClassID,
		ShiftID:   input.ShiftID,
		DayOfWeek: input.DayOfWeek,
		RoomID:    input.RoomID,
	}

	created, err := u.classScheduleRepo.Create(ctx, schedule)
	if err != nil {
		return nil, err
	}

	// Preload the relations for the response
	schedules, _ := u.classScheduleRepo.GetSchedulesByClassID(ctx, input.ClassID)
	for _, s := range schedules {
		if s.ID == created.ID {
			created = &s
			break
		}
	}

	return &CreateClassScheduleOutput{
		Schedule: created,
	}, nil
}
