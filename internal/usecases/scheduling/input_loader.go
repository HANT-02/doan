package scheduling

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
)

func loadSchedulingClasses(
	ctx context.Context,
	classRepo repositoryinterface.ClassRepository,
	classIDs []string,
	teacherIDs []string,
) ([]entities.Class, error) {
	condition := repositories.NewCommonCondition()
	condition.SetPaging(500, 1)
	condition.SetPreload([]string{"Teacher", "Course", "Room", "ClassSchedules", "ClassSchedules.Room", "ClassSchedules.Shift"})
	condition.AddCondition("status", "OPEN", repositories.Equal)
	if len(classIDs) > 0 {
		condition.AddCondition("id", classIDs, repositories.In)
	}
	if len(teacherIDs) > 0 {
		condition.AddCondition("teacher_id", teacherIDs, repositories.In)
	}

	output, err := classRepo.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}

	classes := make([]entities.Class, 0)
	if output == nil {
		return classes, nil
	}

	for _, item := range output.Data {
		if item != nil {
			classes = append(classes, *item)
		}
	}

	return classes, nil
}

func loadSchedulingRooms(
	ctx context.Context,
	roomRepo repositoryinterface.RoomRepository,
	roomIDs []string,
) ([]entities.Room, error) {
	condition := repositories.NewCommonCondition()
	condition.SetPaging(500, 1)
	if len(roomIDs) > 0 {
		condition.AddCondition("id", roomIDs, repositories.In)
	}

	output, err := roomRepo.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}

	rooms := make([]entities.Room, 0)
	if output == nil {
		return rooms, nil
	}

	for _, item := range output.Data {
		if item != nil {
			rooms = append(rooms, *item)
		}
	}

	return rooms, nil
}

func loadActiveShifts(
	ctx context.Context,
	shiftRepo repositoryinterface.ShiftRepository,
) ([]entities.Shift, error) {
	condition := repositories.NewCommonCondition()
	condition.SetPaging(500, 1)
	condition.AddCondition("is_active", true, repositories.Equal)

	output, err := shiftRepo.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}

	shifts := make([]entities.Shift, 0)
	if output == nil {
		return shifts, nil
	}

	for _, item := range output.Data {
		if item != nil {
			shifts = append(shifts, *item)
		}
	}

	return shifts, nil
}
