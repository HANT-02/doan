package scheduling

import (
	"context"
	"time"

	"doan/internal/entities"
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
	classRepo      repositoryinterface.ClassRepository
	roomRepo       repositoryinterface.RoomRepository
	shiftRepo      repositoryinterface.ShiftRepository
	lessonRepo     repositoryinterface.LessonRepository
	enrollmentRepo repositoryinterface.EnrollmentRepository
	store          schedulingservice.PreviewStore[PreviewResult]
	solver         schedulingservice.SchedulingSolver
}

func NewPreviewUseCase(
	classRepo repositoryinterface.ClassRepository,
	roomRepo repositoryinterface.RoomRepository,
	shiftRepo repositoryinterface.ShiftRepository,
	lessonRepo repositoryinterface.LessonRepository,
	enrollmentRepo repositoryinterface.EnrollmentRepository,
	store schedulingservice.PreviewStore[PreviewResult],
	solver schedulingservice.SchedulingSolver,
) PreviewUseCase {
	return &previewUseCase{
		classRepo:      classRepo,
		roomRepo:       roomRepo,
		shiftRepo:      shiftRepo,
		lessonRepo:     lessonRepo,
		enrollmentRepo: enrollmentRepo,
		store:          store,
		solver:         solver,
	}
}

func normalizePreviewDateRange(dateFrom, dateTo time.Time, classes []entities.Class) (time.Time, time.Time) {
	if dateFrom.IsZero() {
		dateFrom = earliestClassStartDate(classes)
	}
	if dateFrom.IsZero() {
		dateFrom = truncateToDate(time.Now())
	}

	if dateTo.IsZero() {
		dateTo = latestClassEndDate(classes)
	}
	if dateTo.IsZero() || dateTo.Before(dateFrom) {
		dateTo = dateFrom.AddDate(0, 3, 0)
	}

	return truncateToDate(dateFrom), truncateToDate(dateTo)
}

func earliestClassStartDate(classes []entities.Class) time.Time {
	var earliest time.Time
	for _, classEntity := range classes {
		if classEntity.StartDate.IsZero() {
			continue
		}
		if earliest.IsZero() || classEntity.StartDate.Before(earliest) {
			earliest = classEntity.StartDate
		}
	}
	return earliest
}

func latestClassEndDate(classes []entities.Class) time.Time {
	var latest time.Time
	for _, classEntity := range classes {
		if classEntity.EndDate == nil || classEntity.EndDate.IsZero() {
			continue
		}
		if latest.IsZero() || classEntity.EndDate.After(latest) {
			latest = *classEntity.EndDate
		}
	}
	return latest
}

func truncateToDate(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func (uc *previewUseCase) Execute(ctx context.Context, input PreviewInput) (*PreviewResult, error) {
	ctxLogger := logger.NewLogger(ctx)

	classes, err := loadSchedulingClasses(ctx, uc.classRepo, input.ClassIDs, input.TeacherIDs)
	if err != nil {
		ctxLogger.Errorf("Failed to load classes for scheduling preview: %v", err)
		return nil, err
	}

	input.DateFrom, input.DateTo = normalizePreviewDateRange(input.DateFrom, input.DateTo, classes)
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

	rooms, err := loadSchedulingRooms(ctx, uc.roomRepo, input.RoomIDs)
	if err != nil {
		ctxLogger.Errorf("Failed to load rooms for scheduling preview: %v", err)
		return nil, err
	}

	shifts, err := loadActiveShifts(ctx, uc.shiftRepo)
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

	previewContext := schedulingservice.BuildPreviewContext(schedulingservice.SolverInput{
		DateFrom:   input.DateFrom,
		DateTo:     input.DateTo,
		ClassIDs:   input.ClassIDs,
		TeacherIDs: input.TeacherIDs,
		RoomIDs:    input.RoomIDs,
		Classes:    classes,
		Rooms:      rooms,
		Shifts:     shifts,
	})

	result.Status = output.Status
	result.Assignments = output.Assignments
	result.Conflicts = append(result.Conflicts, output.Conflicts...)
	result.Summary = output.Summary
	result.CandidateOptions = previewContext.CandidateOptions
	result.Variables = previewContext.Variables
	result.PresetConflicts = previewContext.PresetConflicts
	result.NoDomainConflicts = previewContext.NoDomainConflicts
	result.DomainOptions = previewContext.Domains

	if result.Summary.RequestedClasses == 0 {
		result.Summary.RequestedClasses = len(classes)
	}

	existingLessons, classStudentIDs, existingConflicts, err := uc.collectExistingLessonConflicts(ctx, result)
	if err != nil {
		ctxLogger.Errorf("Failed to collect existing lesson conflicts for scheduling preview: %v", err)
		return nil, err
	}

	result.ExistingLessons = existingLessons
	result.ClassStudentIDs = classStudentIDs
	result.Conflicts = append(result.Conflicts, existingConflicts...)
	result.Summary.ConflictCount = len(result.Conflicts)
	if len(result.Assignments) == 0 {
		result.Status = "FAILED"
	} else if len(result.Conflicts) > 0 || result.Summary.UnscheduledLessons > 0 {
		result.Status = "PARTIAL"
	} else {
		result.Status = "COMPLETED"
	}

	uc.store.Save(runID, result)
	return &result, nil
}
