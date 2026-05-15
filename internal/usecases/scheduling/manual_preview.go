package scheduling

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"doan/internal/entities"
	schedulingservice "doan/internal/services/scheduling"
)

func buildAssignmentMap(assignments []PreviewAssignment) map[string]PreviewAssignment {
	indexed := make(map[string]PreviewAssignment, len(assignments))
	for _, assignment := range assignments {
		indexed[assignment.VariableID] = assignment
	}
	return indexed
}

func buildVariableIndex(variables []Variable) map[string]Variable {
	indexed := make(map[string]Variable, len(variables))
	for _, variable := range variables {
		indexed[variable.ID] = variable
	}
	return indexed
}

func findDomainValueByOptionKey(values []DomainValue, optionKey string) (DomainValue, bool) {
	for _, value := range values {
		if previewCandidateOptionKey(value) == optionKey || legacyPreviewCandidateOptionKey(value) == optionKey {
			return value, true
		}
	}

	return DomainValue{}, false
}

func previewCandidateOptionKey(value DomainValue) string {
	return fmt.Sprintf(
		"%s|%s|%s|%s",
		value.TimeSlot.Start.Format("2006-01-02T15:04"),
		value.TimeSlot.End.Format("2006-01-02T15:04"),
		value.TimeSlot.ShiftID,
		value.RoomID,
	)
}

func legacyPreviewCandidateOptionKey(value DomainValue) string {
	return fmt.Sprintf(
		"%s|%s|%s|%s",
		value.TimeSlot.Start.Format(time.RFC3339),
		value.TimeSlot.End.Format(time.RFC3339),
		value.TimeSlot.ShiftID,
		value.RoomID,
	)
}

func buildAssignmentFromDomain(variable Variable, value DomainValue, constraintFit string) PreviewAssignment {
	return PreviewAssignment{
		VariableID:      variable.ID,
		ClassID:         variable.ClassID,
		ClassCode:       variable.ClassCode,
		ClassName:       variable.ClassName,
		SessionIndex:    variable.SessionIndex,
		SessionTotal:    variable.SessionTotal,
		TeacherID:       variable.TeacherID,
		TeacherLabel:    variable.TeacherLabel,
		RoomID:          value.RoomID,
		RoomName:        value.RoomName,
		RoomCapacity:    value.RoomCapacity,
		ReplaceLessonID: variable.ReplaceLessonID,
		ShiftID:         value.TimeSlot.ShiftID,
		ShiftCode:       value.TimeSlot.ShiftCode,
		ShiftName:       value.TimeSlot.ShiftName,
		ShiftType:       value.TimeSlot.ShiftType,
		StartTime:       value.TimeSlot.Start,
		EndTime:         value.TimeSlot.End,
		ConstraintFit:   constraintFit,
	}
}

func rebuildPreviewResult(base PreviewResult, assignments map[string]PreviewAssignment) PreviewResult {
	result := base
	result.GeneratedAt = time.Now()
	result.Assignments = sortAssignments(assignments)
	result.Conflicts = buildPreviewConflicts(base, assignments)
	result.Summary = PreviewSummary{
		RequestedClasses:   base.Summary.RequestedClasses,
		RequestedSessions:  len(base.Variables),
		ScheduledLessons:   len(result.Assignments),
		UnscheduledLessons: maxInt(len(base.Variables)-len(result.Assignments), 0),
		ConflictCount:      len(result.Conflicts),
		SoftScore:          scorePreviewAssignments(result.Assignments),
	}
	result.Status = derivePreviewStatus(result.Assignments, result.Conflicts)
	return result
}

func buildPreviewConflicts(base PreviewResult, assignments map[string]PreviewAssignment) []PreviewConflict {
	conflicts := make([]PreviewConflict, 0, len(base.PresetConflicts)+len(base.NoDomainConflicts))
	conflicts = append(conflicts, base.PresetConflicts...)

	for _, variable := range base.Variables {
		if _, ok := assignments[variable.ID]; ok {
			continue
		}

		if conflict, ok := base.NoDomainConflicts[variable.ID]; ok {
			conflicts = append(conflicts, conflict)
			continue
		}

		conflicts = append(conflicts, PreviewConflict{
			VariableID:   variable.ID,
			ClassID:      variable.ClassID,
			ClassCode:    variable.ClassCode,
			ClassName:    variable.ClassName,
			SessionIndex: variable.SessionIndex,
			SessionTotal: variable.SessionTotal,
			Type:         "NO_DOMAIN",
			Message:      "Không tìm thấy phương án hợp lệ với hard constraints hiện tại.",
		})
	}

	for _, conflict := range buildAssignmentConflicts(sortAssignments(assignments), base.TravelMap, base.RoomsByID) {
		conflicts = append(conflicts, conflict)
	}

	for _, conflict := range buildExistingLessonPreviewConflicts(sortAssignments(assignments), base.ExistingLessons, base.ClassStudentIDs) {
		conflicts = append(conflicts, conflict)
	}

	return conflicts
}

func buildAssignmentConflicts(assignments []PreviewAssignment, travelMap map[string]int, roomsByID map[string]entities.Room) []PreviewConflict {
	reasonMap := make(map[string][]string)

	for i := 0; i < len(assignments); i++ {
		for j := i + 1; j < len(assignments); j++ {
			left := assignments[i]
			right := assignments[j]
			if !assignmentsOverlap(left, right) {
				if left.TeacherID != "" && left.TeacherID == right.TeacherID && sameCalendarDay(left.StartTime, right.StartTime) {
					var previousEnd, nextStart time.Time
					var fromRoom, toRoom entities.Room

					if left.EndTime.Before(right.StartTime) || left.EndTime.Equal(right.StartTime) {
						previousEnd = left.EndTime
						nextStart = right.StartTime
						fromRoom = roomsByID[left.RoomID]
						toRoom = roomsByID[right.RoomID]
					} else {
						previousEnd = right.EndTime
						nextStart = left.StartTime
						fromRoom = roomsByID[right.RoomID]
						toRoom = roomsByID[left.RoomID]
					}

					if !schedulingservice.HasSufficientTravelGap(previousEnd, nextStart, &fromRoom, &toRoom, travelMap) {
						reasonMap[left.VariableID] = append(reasonMap[left.VariableID], buildConflictReason("di chuyển không kịp", right))
						reasonMap[right.VariableID] = append(reasonMap[right.VariableID], buildConflictReason("di chuyển không kịp", left))
					}
				}
				continue
			}

			if left.ClassID == right.ClassID {
				reasonMap[left.VariableID] = append(reasonMap[left.VariableID], buildConflictReason("lớp học", right))
				reasonMap[right.VariableID] = append(reasonMap[right.VariableID], buildConflictReason("lớp học", left))
			}
			if left.TeacherID != "" && left.TeacherID == right.TeacherID {
				reasonMap[left.VariableID] = append(reasonMap[left.VariableID], buildConflictReason("giáo viên", right))
				reasonMap[right.VariableID] = append(reasonMap[right.VariableID], buildConflictReason("giáo viên", left))
			}
			if left.RoomID != "" && left.RoomID == right.RoomID {
				reasonMap[left.VariableID] = append(reasonMap[left.VariableID], buildConflictReason("phòng học", right))
				reasonMap[right.VariableID] = append(reasonMap[right.VariableID], buildConflictReason("phòng học", left))
			}
		}
	}

	conflicts := make([]PreviewConflict, 0, len(reasonMap))
	for _, assignment := range assignments {
		reasons := uniqueStrings(reasonMap[assignment.VariableID])
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
			Type:         "ASSIGNMENT_CONFLICT",
			Message:      "Buổi học đang bị trùng " + strings.Join(reasons, "; "),
		})
	}

	return conflicts
}

func buildConflictReason(scope string, other PreviewAssignment) string {
	return fmt.Sprintf(
		"%s với %s (%s buổi %d/%d lúc %s)",
		scope,
		other.ClassName,
		other.ClassCode,
		other.SessionIndex,
		other.SessionTotal,
		other.StartTime.Format("02/01 15:04"),
	)
}

func uniqueStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(values))
	items := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		items = append(items, value)
	}

	return items
}

func assignmentsOverlap(left, right PreviewAssignment) bool {
	return left.StartTime.Before(right.EndTime) && right.StartTime.Before(left.EndTime)
}

func sortAssignments(assignments map[string]PreviewAssignment) []PreviewAssignment {
	items := make([]PreviewAssignment, 0, len(assignments))
	for _, assignment := range assignments {
		items = append(items, assignment)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].StartTime.Equal(items[j].StartTime) {
			return items[i].VariableID < items[j].VariableID
		}
		return items[i].StartTime.Before(items[j].StartTime)
	})

	return items
}

func scorePreviewAssignments(assignments []PreviewAssignment) int {
	score := 0
	for i := 0; i < len(assignments)-1; i++ {
		current := assignments[i]
		next := assignments[i+1]
		if current.TeacherID == next.TeacherID && sameCalendarDay(current.StartTime, next.StartTime) {
			gap := next.StartTime.Sub(current.EndTime)
			if gap >= 0 && gap <= 30*time.Minute {
				score += 10
			} else if gap > 2*time.Hour {
				score -= int(gap.Hours()) * 3
			}
		}
	}

	return score
}

func sameCalendarDay(left, right time.Time) bool {
	return left.Year() == right.Year() && left.Month() == right.Month() && left.Day() == right.Day()
}

func derivePreviewStatus(assignments []PreviewAssignment, conflicts []PreviewConflict) string {
	switch {
	case len(assignments) == 0:
		return "FAILED"
	case len(conflicts) > 0:
		return "PARTIAL"
	default:
		return "COMPLETED"
	}
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}

func buildExistingLessonPreviewConflicts(
	assignments []PreviewAssignment,
	existingLessons []ExistingLesson,
	classStudentIDs map[string]map[string]struct{},
) []PreviewConflict {
	conflicts := make([]PreviewConflict, 0)
	for _, assignment := range assignments {
		for _, lesson := range existingLessons {
			if !(assignment.StartTime.Before(lesson.EndTime) && lesson.StartTime.Before(assignment.EndTime)) {
				continue
			}

			reasons := make([]string, 0, 4)
			if lesson.ClassID == assignment.ClassID {
				reasons = append(reasons, "trùng lớp")
			}
			if lesson.TeacherID != "" && lesson.TeacherID == assignment.TeacherID {
				reasons = append(reasons, "trùng giáo viên")
			}
			if lesson.RoomID != "" && lesson.RoomID == assignment.RoomID {
				reasons = append(reasons, "trùng phòng")
			}
			if hasStudentIntersection(classStudentIDs[assignment.ClassID], sliceToSet(lesson.StudentIDs)) {
				reasons = append(reasons, "trùng học sinh")
			}
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
					lesson.StartTime.Format("02/01/2006 15:04"),
					lesson.EndTime.Format("15:04"),
					strings.Join(reasons, ", "),
				),
			})
		}
	}
	return conflicts
}

func sliceToSet(items []string) map[string]struct{} {
	if len(items) == 0 {
		return nil
	}

	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		if item == "" {
			continue
		}
		set[item] = struct{}{}
	}
	return set
}
