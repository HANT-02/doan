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
	Mode       string
	ClassIDs   []string
	TeacherIDs []string
	RoomIDs    []string
}

func normalizePreviewMode(mode string) string {
	switch mode {
	case schedulingservice.PreviewModeColdStart:
		return schedulingservice.PreviewModeColdStart
	case schedulingservice.PreviewModeReplanDraft:
		return schedulingservice.PreviewModeReplanDraft
	case schedulingservice.PreviewModeReplanWithPublishedLock:
		return schedulingservice.PreviewModeReplanWithPublishedLock
	default:
		return schedulingservice.PreviewModeReplanWithPublishedLock
	}
}

type PreviewUseCase interface {
	Execute(ctx context.Context, input PreviewInput) (*PreviewResult, error)
}

type previewUseCase struct {
	classRepo        repositoryinterface.ClassRepository
	roomRepo         repositoryinterface.RoomRepository
	shiftRepo        repositoryinterface.ShiftRepository
	lessonRepo       repositoryinterface.LessonRepository
	enrollmentRepo   repositoryinterface.EnrollmentRepository
	campusTravelRepo repositoryinterface.CampusTravelTimeRepository
	store            schedulingservice.PreviewStore[PreviewResult]
	solver           schedulingservice.SchedulingSolver
}

func NewPreviewUseCase(
	classRepo repositoryinterface.ClassRepository,
	roomRepo repositoryinterface.RoomRepository,
	shiftRepo repositoryinterface.ShiftRepository,
	lessonRepo repositoryinterface.LessonRepository,
	enrollmentRepo repositoryinterface.EnrollmentRepository,
	campusTravelRepo repositoryinterface.CampusTravelTimeRepository,
	store schedulingservice.PreviewStore[PreviewResult],
	solver schedulingservice.SchedulingSolver,
) PreviewUseCase {
	return &previewUseCase{
		classRepo:        classRepo,
		roomRepo:         roomRepo,
		shiftRepo:        shiftRepo,
		lessonRepo:       lessonRepo,
		enrollmentRepo:   enrollmentRepo,
		campusTravelRepo: campusTravelRepo,
		store:            store,
		solver:           solver,
	}
}

func normalizePreviewDateRange(dateFrom, dateTo time.Time, classes []entities.Class) (time.Time, time.Time) {
	today := truncateToDate(time.Now())

	if dateFrom.IsZero() {
		dateFrom = earliestClassStartDate(classes)
	}
	if dateFrom.IsZero() {
		dateFrom = today
	}
	dateFrom = truncateToDate(dateFrom)
	if dateFrom.Before(today) {
		dateFrom = today
	}

	if dateTo.IsZero() {
		dateTo = latestClassEndDate(classes)
	}
	if !dateTo.IsZero() {
		dateTo = truncateToDate(dateTo)
	}
	if dateTo.IsZero() || dateTo.Before(dateFrom) {
		dateTo = dateFrom.AddDate(0, 3, 0)
	}

	return dateFrom, dateTo
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

func collectClassIDs(classes []entities.Class) []string {
	if len(classes) == 0 {
		return nil
	}

	ids := make([]string, 0, len(classes))
	seen := make(map[string]struct{}, len(classes))
	for _, classEntity := range classes {
		if classEntity.ID == "" {
			continue
		}
		if _, ok := seen[classEntity.ID]; ok {
			continue
		}
		seen[classEntity.ID] = struct{}{}
		ids = append(ids, classEntity.ID)
	}
	return ids
}

func truncateToDate(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func (uc *previewUseCase) Execute(ctx context.Context, input PreviewInput) (*PreviewResult, error) {
	ctxLogger := logger.NewLogger(ctx)
	input.Mode = normalizePreviewMode(input.Mode)

	classes, err := loadSchedulingClasses(ctx, uc.classRepo, input.ClassIDs, input.TeacherIDs)
	if err != nil {
		ctxLogger.Errorf("Failed to load classes for scheduling preview: %v", err)
		return nil, err
	}

	input.DateFrom, input.DateTo = normalizePreviewDateRange(input.DateFrom, input.DateTo, classes)
	eligibleClasses, enrollmentConflicts, err := filterClassesByEnrollment(ctx, uc.enrollmentRepo, classes)
	if err != nil {
		ctxLogger.Errorf("Failed to evaluate class enrollment eligibility for scheduling preview: %v", err)
		return nil, err
	}
	classes = eligibleClasses
	if input.DateTo.Before(input.DateFrom) {
		return &PreviewResult{
			RunID:             utils.GenerateUUIDWithPrefix("sched-preview-"),
			Mode:              input.Mode,
			Status:            "FAILED",
			GeneratedAt:       time.Now(),
			EffectiveDateFrom: input.DateFrom,
			Filters: PreviewFilters{
				DateFrom:          input.DateFrom,
				DateTo:            input.DateTo,
				EffectiveDateFrom: input.DateFrom,
				Mode:              input.Mode,
				ClassIDs:          input.ClassIDs,
				TeacherIDs:        input.TeacherIDs,
				RoomIDs:           input.RoomIDs,
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

	shifts, err := loadActiveShifts(ctx, uc.shiftRepo)
	if err != nil {
		ctxLogger.Errorf("Failed to load shifts for scheduling preview: %v", err)
		return nil, err
	}

	classWindows := schedulingservice.BuildClassSchedulingWindows(input.DateFrom, classes, shifts, time.Now())
	previewWindowFrom, previewWindowTo := schedulingservice.AggregateClassSchedulingWindow(input.DateFrom, input.DateTo, classWindows)
	input.DateFrom = previewWindowFrom
	input.DateTo = previewWindowTo

	rooms, err := loadSchedulingRooms(ctx, uc.roomRepo, input.RoomIDs)
	if err != nil {
		ctxLogger.Errorf("Failed to load rooms for scheduling preview: %v", err)
		return nil, err
	}

	travelTimePage, err := uc.campusTravelRepo.GetByCondition(ctx, repositories.NewCommonCondition().WithPaging(1000, 1))
	if err != nil {
		ctxLogger.Errorf("Failed to load campus travel times for scheduling preview: %v", err)
		return nil, err
	}
	travelTimes := make([]entities.CampusTravelTime, 0, len(travelTimePage.Data))
	for _, ptr := range travelTimePage.Data {
		if ptr != nil {
			travelTimes = append(travelTimes, *ptr)
		}
	}
	travelMap := schedulingservice.BuildCampusTravelTimeMap(travelTimes)

	roomsByID := make(map[string]entities.Room, len(rooms))
	for _, room := range rooms {
		roomsByID[room.ID] = room
	}

	runID := utils.GenerateUUIDWithPrefix("sched-preview-")
	result := PreviewResult{
		RunID:             runID,
		Mode:              input.Mode,
		Status:            "FAILED",
		GeneratedAt:       time.Now(),
		EffectiveDateFrom: input.DateFrom,
		Filters: PreviewFilters{
			DateFrom:          input.DateFrom,
			DateTo:            input.DateTo,
			EffectiveDateFrom: input.DateFrom,
			Mode:              input.Mode,
			ClassIDs:          input.ClassIDs,
			TeacherIDs:        input.TeacherIDs,
			RoomIDs:           input.RoomIDs,
		},
		Assignments: []PreviewAssignment{},
		Conflicts:   append([]PreviewConflict{}, enrollmentConflicts...),
	}

	if len(classes) == 0 && len(enrollmentConflicts) == 0 {
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

	targetLessonStatuses := []string{entities.LessonStatusPublished, entities.LessonStatusDraft}
	if input.Mode == schedulingservice.PreviewModeColdStart {
		targetLessonStatuses = nil
	}
	var targetLessonsList []entities.Lesson
	targetClassIDs := collectClassIDs(classes)
	if len(targetLessonStatuses) > 0 && len(targetClassIDs) > 0 {
		targetLessonsList, err = uc.lessonRepo.FindOverlappingLessons(
			ctx,
			input.DateFrom,
			input.DateTo.Add(24*time.Hour),
			targetClassIDs,
			nil,
			nil,
			targetLessonStatuses,
		)
		if err != nil {
			ctxLogger.Errorf("Failed to load target lessons for scheduling preview: %v", err)
			return nil, err
		}
	}

	output, err := uc.solver.Solve(ctx, schedulingservice.SolverInput{
		DateFrom:      input.DateFrom,
		DateTo:        input.DateTo,
		ClassIDs:      input.ClassIDs,
		TeacherIDs:    input.TeacherIDs,
		RoomIDs:       input.RoomIDs,
		Classes:       classes,
		ClassWindows:  classWindows,
		Rooms:         rooms,
		Shifts:        shifts,
		RoomsByID:     roomsByID,
		TravelMap:     travelMap,
		TargetLessons: targetLessonsList,
	})
	if err != nil {
		return nil, err
	}

	previewContext := schedulingservice.BuildPreviewContext(schedulingservice.SolverInput{
		DateFrom:      input.DateFrom,
		DateTo:        input.DateTo,
		ClassIDs:      input.ClassIDs,
		TeacherIDs:    input.TeacherIDs,
		RoomIDs:       input.RoomIDs,
		Classes:       classes,
		ClassWindows:  classWindows,
		Rooms:         rooms,
		Shifts:        shifts,
		RoomsByID:     roomsByID,
		TravelMap:     travelMap,
		TargetLessons: targetLessonsList,
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
	result.RoomsByID = roomsByID
	result.TravelMap = travelMap

	if result.Summary.RequestedClasses == 0 {
		result.Summary.RequestedClasses = len(classes)
	}

	existingLessons, classStudentIDs, existingConflicts, err := uc.collectExistingLessonConflicts(ctx, result, travelMap, roomsByID)
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
	result = maybeAppendConflictDensityConflict(result)

	// B5: Capacity utilization calculation
	classMaxStudents := make(map[string]int)
	for _, classEntity := range classes {
		classMaxStudents[classEntity.ID] = classEntity.MaxStudents
	}

	totalStudentCount := 0
	totalCapacity := 0

	for i := range result.Assignments {
		classID := result.Assignments[i].ClassID
		studentCount := len(result.ClassStudentIDs[classID])
		result.Assignments[i].ExpectedStudentCount = studentCount

		maxStudents := classMaxStudents[classID]
		roomCapacity := result.Assignments[i].RoomCapacity

		limit := CalculateCapacityLimit(roomCapacity, maxStudents)

		totalStudentCount += studentCount
		totalCapacity += limit
	}

	if totalCapacity > 0 {
		result.Summary.AverageCapacityUtilization = float64(totalStudentCount) / float64(totalCapacity)
	}

	uc.store.Save(runID, result)
	return &result, nil
}
