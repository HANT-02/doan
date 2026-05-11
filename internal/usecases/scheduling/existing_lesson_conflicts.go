package scheduling

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"doan/internal/entities"
)

func (uc *previewUseCase) collectExistingLessonConflicts(
	ctx context.Context,
	preview PreviewResult,
) ([]ExistingLesson, map[string]map[string]struct{}, []PreviewConflict, error) {
	if preview.Filters.DateTo.Before(preview.Filters.DateFrom) {
		return []ExistingLesson{}, map[string]map[string]struct{}{}, []PreviewConflict{}, nil
	}

	from, to := preview.Filters.DateFrom, preview.Filters.DateTo.Add(24*time.Hour)
	if len(preview.Assignments) > 0 {
		from, to = previewAssignmentWindow(preview.Assignments)
	}

	lessons, err := uc.lessonRepo.FindOverlappingLessons(
		ctx,
		from,
		to,
		collectAssignmentClassIDs(preview.Assignments),
		collectAssignmentTeacherIDs(preview.Assignments),
		collectAssignmentRoomIDs(preview.Assignments),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	if len(lessons) == 0 {
		return []ExistingLesson{}, map[string]map[string]struct{}{}, []PreviewConflict{}, nil
	}

	classIDs := make(map[string]struct{})
	for _, variable := range preview.Variables {
		if variable.ClassID != "" {
			classIDs[variable.ClassID] = struct{}{}
		}
	}
	for _, lesson := range lessons {
		if lesson.ClassID != "" {
			classIDs[lesson.ClassID] = struct{}{}
		}
	}

	studentIDsByClass, err := uc.loadStudentIDsByClass(ctx, classIDs)
	if err != nil {
		return nil, nil, nil, err
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
			TeacherID:    teacherID,
			TeacherLabel: teacherLabel,
			RoomID:       roomID,
			RoomName:     lesson.Room.Name,
			StartTime:    lesson.DateStart,
			EndTime:      lesson.DateEnd,
			Notes:        lesson.Notes,
			StudentIDs:   sortedKeys(studentIDsByClass[lesson.ClassID]),
		})
	}

	conflicts := make([]PreviewConflict, 0)
	for _, assignment := range preview.Assignments {
		studentSet := studentIDsByClass[assignment.ClassID]
		for _, lesson := range lessons {
			if !overlaps(assignment.StartTime, assignment.EndTime, lesson.DateStart, lesson.DateEnd) {
				continue
			}

			reasons := buildExistingLessonConflictReasons(assignment, lesson, studentSet, studentIDsByClass[lesson.ClassID])
			if len(reasons) == 0 {
				continue
			}

			conflicts = append(conflicts, PreviewConflict{
				VariableID:   assignment.VariableID,
				ClassID:      assignment.ClassID,
				ClassCode:    assignment.ClassCode,
				ClassName:    assignment.ClassName,
				SessionIndex: assignment.SessionIndex,
				SessionTotal: assignment.SessionTotal,
				Type:         "SYSTEM_LESSON_CONFLICT",
				Message: fmt.Sprintf(
					"Trùng với lesson đã lưu [%s - %s] vì %s.",
					lesson.DateStart.Format("02/01/2006 15:04"),
					lesson.DateEnd.Format("15:04"),
					strings.Join(reasons, ", "),
				),
			})
		}
	}

	return events, studentIDsByClass, conflicts, nil
}

func (uc *previewUseCase) loadStudentIDsByClass(
	ctx context.Context,
	classIDs map[string]struct{},
) (map[string]map[string]struct{}, error) {
	output := make(map[string]map[string]struct{}, len(classIDs))
	for classID := range classIDs {
		if classID == "" {
			continue
		}

		enrollments, err := uc.enrollmentRepo.ListByClassID(ctx, classID)
		if err != nil {
			return nil, err
		}

		students := make(map[string]struct{})
		for _, enrollment := range enrollments {
			if enrollment.Status != "" && enrollment.Status != "ENROLLED" {
				continue
			}
			if enrollment.StudentID == "" {
				continue
			}
			students[enrollment.StudentID] = struct{}{}
		}
		output[classID] = students
	}

	return output, nil
}

func buildExistingLessonConflictReasons(
	assignment PreviewAssignment,
	lesson entities.Lesson,
	previewStudentIDs map[string]struct{},
	existingStudentIDs map[string]struct{},
) []string {
	reasons := make([]string, 0, 4)
	switch {
	case lesson.ClassID == assignment.ClassID:
		reasons = append(reasons, "trùng lớp")
	}

	if lesson.TeacherID != nil && *lesson.TeacherID != "" && *lesson.TeacherID == assignment.TeacherID {
		reasons = append(reasons, "trùng giáo viên")
	}

	if lesson.RoomID != nil && *lesson.RoomID != "" && *lesson.RoomID == assignment.RoomID {
		reasons = append(reasons, "trùng phòng")
	}

	if hasStudentIntersection(previewStudentIDs, existingStudentIDs) {
		reasons = append(reasons, "trùng học sinh")
	}

	return reasons
}

func hasStudentIntersection(left map[string]struct{}, right map[string]struct{}) bool {
	if len(left) == 0 || len(right) == 0 {
		return false
	}

	if len(left) > len(right) {
		left, right = right, left
	}

	for studentID := range left {
		if _, ok := right[studentID]; ok {
			return true
		}
	}

	return false
}

func sortedKeys(items map[string]struct{}) []string {
	if len(items) == 0 {
		return nil
	}

	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
