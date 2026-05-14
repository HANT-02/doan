package scheduling

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	schedulingstore "doan/internal/services/scheduling"
	"doan/pkg/logger"
)

type CommitPreviewInput struct {
	RunID             string
	ManualAssignments []ManualAssignmentOverride
}

type CommitPreviewOutput struct {
	Message          string `json:"message"`
	ScheduledLessons int    `json:"scheduled_lessons"`
	Status           string `json:"status"`
}

type CommitPreviewConflictError struct {
	Message string
	Preview PreviewResult
}

func (err *CommitPreviewConflictError) Error() string {
	return err.Message
}

type CommitPreviewUseCase interface {
	Execute(ctx context.Context, input CommitPreviewInput) (*CommitPreviewOutput, error)
}

type commitPreviewUseCase struct {
	lessonRepo        repositoryinterface.LessonRepository
	classScheduleRepo repositoryinterface.ClassScheduleRepository
	uow               repositories.UnitOfWork
	log               logger.Logger
	store             schedulingstore.PreviewStore[PreviewResult]
}

func NewCommitPreviewUseCase(
	lessonRepo repositoryinterface.LessonRepository,
	classScheduleRepo repositoryinterface.ClassScheduleRepository,
	uow repositories.UnitOfWork,
	log logger.Logger,
	store schedulingstore.PreviewStore[PreviewResult],
) CommitPreviewUseCase {
	return &commitPreviewUseCase{
		lessonRepo:        lessonRepo,
		classScheduleRepo: classScheduleRepo,
		uow:               uow,
		log:               log,
		store:             store,
	}
}

func (uc *commitPreviewUseCase) Execute(ctx context.Context, input CommitPreviewInput) (*CommitPreviewOutput, error) {
	if input.RunID == "" {
		return nil, errors.New("run_id is required")
	}

	preview, ok := uc.store.Get(input.RunID)
	if !ok {
		return nil, errors.New("preview run not found")
	}

	if len(preview.Assignments) == 0 {
		return nil, errors.New("preview does not contain any assignment to commit")
	}

	assignmentsByID := buildAssignmentMap(preview.Assignments)
	manualAssignmentIDs := make(map[string]struct{}, len(input.ManualAssignments))
	if len(input.ManualAssignments) > 0 {
		variableIndex := buildVariableIndex(preview.Variables)
		for _, override := range input.ManualAssignments {
			variable, ok := variableIndex[override.VariableID]
			if !ok {
				return nil, fmt.Errorf("khong tim thay session %s trong preview hien tai", override.VariableID)
			}

			domainValue, ok := findDomainValueByOptionKey(preview.DomainOptions[override.VariableID], override.OptionKey)
			if !ok {
				return nil, fmt.Errorf("phuong an da chon cho session %s khong hop le hoac da het hieu luc", override.VariableID)
			}

			assignmentsByID[override.VariableID] = buildAssignmentFromDomain(variable, domainValue, "MANUAL_OVERRIDE")
			manualAssignmentIDs[override.VariableID] = struct{}{}
		}

		preview = rebuildPreviewResult(preview, assignmentsByID)
		uc.store.Save(preview.RunID, preview)
	} else {
		preview = rebuildPreviewResult(preview, assignmentsByID)
	}

	manualAdjustmentLimit := allowedManualAdjustmentLimit(preview.Summary.RequestedSessions)
	if len(input.ManualAssignments) > manualAdjustmentLimit {
		preview = appendOperationalConflict(preview, PreviewConflict{
			Type: "EXCESSIVE_MANUAL_ADJUSTMENT",
			Message: fmt.Sprintf(
				"Preview đang cần %d chỉnh tay, vượt ngưỡng cho phép %d chỉnh tay cho %d ca học. Hãy tối ưu lại dữ liệu đầu vào hoặc chạy lại solver thay vì vá tay quá nhiều.",
				len(input.ManualAssignments),
				manualAdjustmentLimit,
				preview.Summary.RequestedSessions,
			),
		})
		uc.store.Save(preview.RunID, preview)
		return nil, &CommitPreviewConflictError{
			Message: "so luong chinh tay vuot nguong cho phep de commit",
			Preview: preview,
		}
	}

	if preview.Status != "COMPLETED" {
		return nil, fmt.Errorf(
			"preview %s chưa thể commit vì còn %d buổi chưa xếp và %d conflict. Hãy chạy lại preview đến khi trạng thái COMPLETED",
			input.RunID,
			preview.Summary.UnscheduledLessons,
			preview.Summary.ConflictCount,
		)
	}

	committedLessons := 0
	_, err := repositories.ExecuteInTransaction(ctx, uc.uow, uc.log, func(txCtx context.Context) (interface{}, error) {
		from, to := previewAssignmentWindow(preview.Assignments)
		existingLessons, err := uc.lessonRepo.FindOverlappingLessons(
			txCtx,
			from,
			to,
			collectAssignmentClassIDs(preview.Assignments),
			collectAssignmentTeacherIDs(preview.Assignments),
			collectAssignmentRoomIDs(preview.Assignments),
			[]string{
				entities.LessonStatusPublished,
				entities.LessonStatusDraft,
				entities.LessonStatusUnplanned,
			},
		)
		if err != nil {
			return nil, err
		}

		preview = refreshPreviewWithExistingLessons(preview, existingLessons)
		uc.store.Save(preview.RunID, preview)

		if conflicts := detectCommitConflicts(preview.Assignments, existingLessons); len(conflicts) > 0 {
			return nil, &CommitPreviewConflictError{
				Message: formatCommitConflicts(conflicts),
				Preview: preview,
			}
		}

		for _, assignment := range preview.Assignments {
			teacherID := assignment.TeacherID
			roomID := assignment.RoomID
			publishedAt := time.Now().UTC()
			sourcePreviewRunID := input.RunID

			_, err := uc.lessonRepo.Create(txCtx, &entities.Lesson{
				ClassID:          assignment.ClassID,
				TeacherID:        &teacherID,
				RoomID:           &roomID,
				DateStart:        assignment.StartTime,
				DateEnd:          assignment.EndTime,
				Status:           entities.LessonStatusPublished,
				PublishedAt:      &publishedAt,
				SourcePreviewRun: &sourcePreviewRunID,
				ChangeReason:     entities.LessonChangeReasonInitialSchedulingCommit,
				Notes: fmt.Sprintf(
					"Generated from scheduling preview %s - session %d/%d",
					input.RunID,
					assignment.SessionIndex,
					assignment.SessionTotal,
				),
			})
			if err != nil {
				return nil, err
			}

			committedLessons++
		}

		return committedLessons, nil
	})
	if err != nil {
		return nil, err
	}

	return &CommitPreviewOutput{
		Message: fmt.Sprintf(
			"Da tao %d lesson tu preview %s. Neu can xem lai, hay vao timetable/quan ly lop hoc de doi chieu ket qua.",
			committedLessons,
			input.RunID,
		),
		ScheduledLessons: committedLessons,
		Status:           "COMMITTED",
	}, nil
}

func weekdayToScheduleDay(weekday time.Weekday) string {
	switch weekday {
	case time.Monday:
		return "MONDAY"
	case time.Tuesday:
		return "TUESDAY"
	case time.Wednesday:
		return "WEDNESDAY"
	case time.Thursday:
		return "THURSDAY"
	case time.Friday:
		return "FRIDAY"
	case time.Saturday:
		return "SATURDAY"
	case time.Sunday:
		return "SUNDAY"
	default:
		return ""
	}
}

type commitConflict struct {
	Assignment PreviewAssignment
	Lesson     entities.Lesson
	Reason     string
}

func previewAssignmentWindow(assignments []PreviewAssignment) (time.Time, time.Time) {
	if len(assignments) == 0 {
		return time.Time{}, time.Time{}
	}

	from := assignments[0].StartTime
	to := assignments[0].EndTime
	for _, assignment := range assignments[1:] {
		if assignment.StartTime.Before(from) {
			from = assignment.StartTime
		}
		if assignment.EndTime.After(to) {
			to = assignment.EndTime
		}
	}

	return from, to
}

func collectAssignmentClassIDs(assignments []PreviewAssignment) []string {
	ids := make([]string, 0)
	seen := make(map[string]struct{})
	for _, assignment := range assignments {
		if _, ok := seen[assignment.ClassID]; ok {
			continue
		}
		seen[assignment.ClassID] = struct{}{}
		ids = append(ids, assignment.ClassID)
	}

	return ids
}

func collectAssignmentTeacherIDs(assignments []PreviewAssignment) []string {
	ids := make([]string, 0)
	seen := make(map[string]struct{})
	for _, assignment := range assignments {
		if assignment.TeacherID == "" {
			continue
		}
		if _, ok := seen[assignment.TeacherID]; ok {
			continue
		}
		seen[assignment.TeacherID] = struct{}{}
		ids = append(ids, assignment.TeacherID)
	}

	return ids
}

func collectAssignmentRoomIDs(assignments []PreviewAssignment) []string {
	ids := make([]string, 0)
	seen := make(map[string]struct{})
	for _, assignment := range assignments {
		if assignment.RoomID == "" {
			continue
		}
		if _, ok := seen[assignment.RoomID]; ok {
			continue
		}
		seen[assignment.RoomID] = struct{}{}
		ids = append(ids, assignment.RoomID)
	}

	return ids
}

func detectCommitConflicts(assignments []PreviewAssignment, lessons []entities.Lesson) []commitConflict {
	conflicts := make([]commitConflict, 0)
	for _, assignment := range assignments {
		for _, lesson := range lessons {
			if !overlaps(assignment.StartTime, assignment.EndTime, lesson.DateStart, lesson.DateEnd) {
				continue
			}

			reason := ""
			switch {
			case lesson.ClassID == assignment.ClassID:
				reason = "lop hoc da co lesson trung gio"
			case lesson.TeacherID != nil && *lesson.TeacherID == assignment.TeacherID:
				reason = "giao vien da co lesson trung gio"
			case lesson.RoomID != nil && *lesson.RoomID == assignment.RoomID:
				reason = "phong hoc da co lesson trung gio"
			default:
				continue
			}

			conflicts = append(conflicts, commitConflict{
				Assignment: assignment,
				Lesson:     lesson,
				Reason:     reason,
			})
		}
	}

	return conflicts
}

func refreshPreviewWithExistingLessons(preview PreviewResult, lessons []entities.Lesson) PreviewResult {
	preview.ExistingLessons = mergeExistingLessons(
		preview.ExistingLessons,
		buildExistingLessonEventsForPreview(lessons, preview.ClassStudentIDs),
	)
	return rebuildPreviewResult(preview, buildAssignmentMap(preview.Assignments))
}

func buildExistingLessonEventsForPreview(
	lessons []entities.Lesson,
	classStudentIDs map[string]map[string]struct{},
) []ExistingLesson {
	if len(lessons) == 0 {
		return []ExistingLesson{}
	}

	events := make([]ExistingLesson, 0, len(lessons))
	for _, lesson := range lessons {
		teacherLabel := ""
		if lesson.TeacherID != nil {
			teacherLabel = *lesson.TeacherID
		}
		if lesson.Teacher.FullName != "" {
			teacherLabel = lesson.Teacher.FullName
		} else if lesson.Teacher.Code != "" {
			teacherLabel = lesson.Teacher.Code
		}

		roomID := ""
		if lesson.RoomID != nil {
			roomID = *lesson.RoomID
		}

		teacherID := ""
		if lesson.TeacherID != nil {
			teacherID = *lesson.TeacherID
		}

		events = append(events, ExistingLesson{
			LessonID:     lesson.ID,
			ClassID:      lesson.ClassID,
			ClassCode:    lesson.Class.Code,
			ClassName:    lesson.Class.Name,
			Status:       lesson.Status,
			TeacherID:    teacherID,
			TeacherLabel: teacherLabel,
			RoomID:       roomID,
			RoomName:     lesson.Room.Name,
			StartTime:    lesson.DateStart,
			EndTime:      lesson.DateEnd,
			Notes:        lesson.Notes,
			StudentIDs:   sortedKeys(classStudentIDs[lesson.ClassID]),
		})
	}

	return events
}

func mergeExistingLessons(current []ExistingLesson, incoming []ExistingLesson) []ExistingLesson {
	if len(current) == 0 {
		return append([]ExistingLesson(nil), incoming...)
	}
	if len(incoming) == 0 {
		return append([]ExistingLesson(nil), current...)
	}

	merged := make(map[string]ExistingLesson, len(current)+len(incoming))
	for _, lesson := range current {
		if lesson.LessonID == "" {
			continue
		}
		merged[lesson.LessonID] = lesson
	}
	for _, lesson := range incoming {
		if lesson.LessonID == "" {
			continue
		}
		merged[lesson.LessonID] = lesson
	}

	items := make([]ExistingLesson, 0, len(merged))
	for _, lesson := range merged {
		items = append(items, lesson)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].StartTime.Equal(items[j].StartTime) {
			return items[i].LessonID < items[j].LessonID
		}
		return items[i].StartTime.Before(items[j].StartTime)
	})

	return items
}

func formatCommitConflicts(conflicts []commitConflict) string {
	if len(conflicts) == 0 {
		return ""
	}

	lines := make([]string, 0, minInt(len(conflicts), 5))
	for index, conflict := range conflicts {
		if index >= 5 {
			break
		}

		lines = append(lines, fmt.Sprintf(
			"- %s (%s buoi %d/%d) xung dot vi %s voi lesson %s dang ton tai [%s - %s]",
			conflict.Assignment.ClassName,
			conflict.Assignment.ClassCode,
			conflict.Assignment.SessionIndex,
			conflict.Assignment.SessionTotal,
			conflict.Reason,
			lessonLifecycleLabel(conflict.Lesson.Status),
			conflict.Lesson.DateStart.Format("02/01/2006 15:04"),
			conflict.Lesson.DateEnd.Format("15:04"),
		))
	}

	suffix := ""
	if len(conflicts) > 5 {
		suffix = fmt.Sprintf("\n... va %d xung dot khac", len(conflicts)-5)
	}

	return "Khong the commit preview vi da ton tai lesson trung lich:\n" + strings.Join(lines, "\n") + suffix
}

func overlaps(startA, endA, startB, endB time.Time) bool {
	return startA.Before(endB) && startB.Before(endA)
}

func minInt(left, right int) int {
	if left < right {
		return left
	}
	return right
}
