package scheduling

import (
	"context"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	schedulingservice "doan/internal/services/scheduling"
	"doan/pkg/logger"
	"doan/pkg/utils"
)

type PreviewInput struct {
	DateFrom   time.Time
	DateTo     time.Time
	ClassIDs   []string
	TeacherIDs []string
	RoomIDs    []string
}

type PreviewUseCase interface {
	Execute(ctx context.Context, input PreviewInput) (*PreviewResult, error)
}

type previewUseCase struct {
	classRepo repositoryinterface.ClassRepository
	roomRepo  repositoryinterface.RoomRepository
	shiftRepo repositoryinterface.ShiftRepository
	store     schedulingservice.PreviewStore[PreviewResult]
	solver    schedulingservice.SchedulingSolver
}

func NewPreviewUseCase(
	classRepo repositoryinterface.ClassRepository,
	roomRepo repositoryinterface.RoomRepository,
	shiftRepo repositoryinterface.ShiftRepository,
	store schedulingservice.PreviewStore[PreviewResult],
	solver schedulingservice.SchedulingSolver,
) PreviewUseCase {
	return &previewUseCase{
		classRepo: classRepo,
		roomRepo:  roomRepo,
		shiftRepo: shiftRepo,
		store:     store,
		solver:    solver,
	}
}

func (uc *previewUseCase) Execute(ctx context.Context, input PreviewInput) (*PreviewResult, error) {
	ctxLogger := logger.NewLogger(ctx)

	if input.DateTo.Before(input.DateFrom) {
		return &PreviewResult{
			RunID:       utils.GenerateUUIDWithPrefix("sched-preview-"),
			Status:      "FAILED",
			GeneratedAt: time.Now(),
			Filters: PreviewFilters{
				DateFrom:   input.DateFrom,
				DateTo:     input.DateTo,
				ClassIDs:   input.ClassIDs,
				TeacherIDs: input.TeacherIDs,
				RoomIDs:    input.RoomIDs,
			},
			Summary:     PreviewSummary{},
			Assignments: []PreviewAssignment{},
			Conflicts: []PreviewConflict{
				{
					Type:    "NO_VALID_DATE_RANGE",
					Message: "Khoảng ngày preview không hợp lệ. Hãy chọn ngày kết thúc lớn hơn hoặc bằng ngày bắt đầu.",
				},
			},
		}, nil
	}

	classes, err := uc.loadClasses(ctx, input)
	if err != nil {
		ctxLogger.Errorf("Failed to load classes for scheduling preview: %v", err)
		return nil, err
	}

	rooms, err := uc.loadRooms(ctx, input)
	if err != nil {
		ctxLogger.Errorf("Failed to load rooms for scheduling preview: %v", err)
		return nil, err
	}

	shifts, err := uc.loadShifts(ctx)
	if err != nil {
		ctxLogger.Errorf("Failed to load shifts for scheduling preview: %v", err)
		return nil, err
	}

	runID := utils.GenerateUUIDWithPrefix("sched-preview-")
	result := PreviewResult{
		RunID:       runID,
		Status:      "FAILED",
		GeneratedAt: time.Now(),
		Filters: PreviewFilters{
			DateFrom:   input.DateFrom,
			DateTo:     input.DateTo,
			ClassIDs:   input.ClassIDs,
			TeacherIDs: input.TeacherIDs,
			RoomIDs:    input.RoomIDs,
		},
		Assignments: []PreviewAssignment{},
		Conflicts:   []PreviewConflict{},
	}

	if len(classes) == 0 {
		result.Conflicts = append(result.Conflicts, PreviewConflict{
			Type:    "NO_CLASS_INPUT",
			Message: "Không có lớp OPEN nào phù hợp bộ lọc hiện tại. Kiểm tra lại khoảng ngày, lớp đã chọn hoặc giáo viên đã lọc.",
		})
	}

	if len(rooms) == 0 {
		result.Conflicts = append(result.Conflicts, PreviewConflict{
			Type:    "NO_ACTIVE_ROOM",
			Message: "Không có phòng khả dụng để xếp lịch. Kiểm tra bộ lọc phòng hoặc dữ liệu phòng học đang hoạt động.",
		})
	}

	if len(shifts) == 0 {
		result.Conflicts = append(result.Conflicts, PreviewConflict{
			Type:    "NO_ACTIVE_SHIFT",
			Message: "Chưa có ca học nào đang hoạt động để sinh slot xếp lịch. Hãy tạo hoặc bật `Shift` trước khi chạy preview.",
		})
	}

	output, err := uc.solver.Solve(ctx, schedulingservice.SolverInput{
		DateFrom:   input.DateFrom,
		DateTo:     input.DateTo,
		ClassIDs:   input.ClassIDs,
		TeacherIDs: input.TeacherIDs,
		RoomIDs:    input.RoomIDs,
		Classes:    classes,
		Rooms:      rooms,
		Shifts:     shifts,
	})
	if err != nil {
		return nil, err
	}

	result.Status = output.Status
	result.Assignments = output.Assignments
	result.Conflicts = append(result.Conflicts, output.Conflicts...)
	result.Summary = output.Summary

	if result.Summary.RequestedClasses == 0 {
		result.Summary.RequestedClasses = len(classes)
	}

	uc.store.Save(runID, result)
	return &result, nil
}

func (uc *previewUseCase) loadClasses(ctx context.Context, input PreviewInput) ([]entities.Class, error) {
	condition := repositories.NewCommonCondition()
	condition.SetPaging(500, 1)
	condition.SetPreload([]string{"Teacher", "Course", "Room", "ClassSchedules", "ClassSchedules.Room", "ClassSchedules.Shift"})
	condition.AddCondition("status", "OPEN", repositories.Equal)
	if len(input.ClassIDs) > 0 {
		condition.AddCondition("id", input.ClassIDs, repositories.In)
	}
	if len(input.TeacherIDs) > 0 {
		condition.AddCondition("teacher_id", input.TeacherIDs, repositories.In)
	}

	output, err := uc.classRepo.GetByCondition(ctx, condition)
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

func (uc *previewUseCase) loadShifts(ctx context.Context) ([]entities.Shift, error) {
	condition := repositories.NewCommonCondition()
	condition.SetPaging(500, 1)
	condition.AddCondition("is_active", true, repositories.Equal)

	output, err := uc.shiftRepo.GetByCondition(ctx, condition)
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

func (uc *previewUseCase) loadRooms(ctx context.Context, input PreviewInput) ([]entities.Room, error) {
	condition := repositories.NewCommonCondition()
	condition.SetPaging(500, 1)
	if len(input.RoomIDs) > 0 {
		condition.AddCondition("id", input.RoomIDs, repositories.In)
	}

	output, err := uc.roomRepo.GetByCondition(ctx, condition)
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
