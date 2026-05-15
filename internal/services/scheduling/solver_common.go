package scheduling

import (
	"fmt"
	"sort"
	"time"

	"doan/internal/entities"
)

type preparedSchedulingProblem struct {
	variables         []Variable
	presetConflicts   []PreviewConflict
	domains           map[string][]DomainValue
	noDomainConflicts map[string]PreviewConflict
	roomsByID         map[string]entities.Room
	travelMap         map[string]int
	targetLessons     map[string]entities.Lesson
}

func prepareSchedulingProblem(input SolverInput) preparedSchedulingProblem {
	variables, presetConflicts := buildVariables(input)
	domains, noDomainConflicts := buildDomains(variables, input.Rooms, input, indexClassSchedules(input.Classes), input.Shifts)

	targetLessonsByClass := make(map[string][]entities.Lesson)
	for _, lesson := range input.TargetLessons {
		targetLessonsByClass[lesson.ClassID] = append(targetLessonsByClass[lesson.ClassID], lesson)
	}
	for classID := range targetLessonsByClass {
		sort.Slice(targetLessonsByClass[classID], func(i, j int) bool {
			return targetLessonsByClass[classID][i].DateStart.Before(targetLessonsByClass[classID][j].DateStart)
		})
	}

	targetLessons := make(map[string]entities.Lesson)
	for index, variable := range variables {
		lessons := targetLessonsByClass[variable.ClassID]
		if variable.SessionIndex > 0 && variable.SessionIndex <= len(lessons) {
			target := lessons[variable.SessionIndex-1]
			targetLessons[variable.ID] = target
			variables[index].ReplaceLessonID = target.ID
		}
	}

	return preparedSchedulingProblem{
		variables:         variables,
		presetConflicts:   presetConflicts,
		domains:           domains,
		noDomainConflicts: noDomainConflicts,
		roomsByID:         input.RoomsByID,
		travelMap:         input.TravelMap,
		targetLessons:     targetLessons,
	}
}

func buildSolverOutput(
	input SolverInput,
	variables []Variable,
	assignments map[string]PreviewAssignment,
	presetConflicts []PreviewConflict,
	noDomainConflicts map[string]PreviewConflict,
	extraConflicts []PreviewConflict,
	targetLessons map[string]entities.Lesson,
) *SolverOutput {
	assignmentSlice := assignmentsToSlice(assignments)
	conflicts := append([]PreviewConflict(nil), presetConflicts...)
	conflicts = append(conflicts, collectUnassignedConflicts(variables, assignments, noDomainConflicts)...)
	conflicts = append(conflicts, extraConflicts...)

	scheduleChanges := 0
	teacherChanges := 0
	roomChanges := 0

	for _, assignment := range assignmentSlice {
		if target, ok := targetLessons[assignment.VariableID]; ok {
			targetTeacherID := ""
			if target.TeacherID != nil {
				targetTeacherID = *target.TeacherID
			}
			targetRoomID := ""
			if target.RoomID != nil {
				targetRoomID = *target.RoomID
			}

			if !assignment.StartTime.Equal(target.DateStart) {
				scheduleChanges++
			}
			if assignment.TeacherID != targetTeacherID {
				teacherChanges++
			}
			if assignment.RoomID != targetRoomID {
				roomChanges++
			}
		}
	}

	summary := SolverSummary{
		RequestedClasses:    len(input.Classes),
		RequestedSessions:   len(variables),
		ScheduledLessons:    len(assignmentSlice),
		UnscheduledLessons:  maxInt(len(variables)-len(assignmentSlice), 0),
		ConflictCount:       len(conflicts),
		SoftScore:           scoreAssignments(assignmentSlice, targetLessons),
		ScheduleChangeCount: scheduleChanges,
		TeacherChangeCount:  teacherChanges,
		RoomChangeCount:     roomChanges,
	}

	status := "FAILED"
	switch {
	case len(assignmentSlice) == 0:
		status = "FAILED"
	case len(conflicts) > 0:
		status = "PARTIAL"
	default:
		status = "COMPLETED"
	}

	return &SolverOutput{
		Status:      status,
		Assignments: assignmentSlice,
		Conflicts:   conflicts,
		Summary:     summary,
	}
}

func collectUnassignedConflicts(
	variables []Variable,
	assignments map[string]PreviewAssignment,
	noDomainConflicts map[string]PreviewConflict,
) []PreviewConflict {
	conflicts := make([]PreviewConflict, 0)
	for _, variable := range variables {
		if _, ok := assignments[variable.ID]; ok {
			continue
		}

		if conflict, ok := noDomainConflicts[variable.ID]; ok {
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

	return conflicts
}

func newPreviewAssignment(variable Variable, value DomainValue, constraintFit string) PreviewAssignment {
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

func (p *preparedSchedulingProblem) hasConflict(variable Variable, value DomainValue, assignments map[string]PreviewAssignment) bool {
	return len(p.findBlockingAssignments(variable, value, assignments)) > 0
}

func (p *preparedSchedulingProblem) findBlockingAssignments(variable Variable, value DomainValue, assignments map[string]PreviewAssignment) []PreviewAssignment {
	blocking := make([]PreviewAssignment, 0)
	if value.RoomCapacity < variable.ExpectedCapcity {
		return blocking
	}

	for _, assignment := range assignments {
		if overlaps(value.TimeSlot.Start, value.TimeSlot.End, assignment.StartTime, assignment.EndTime) {
			if assignment.ClassID == variable.ClassID || assignment.TeacherID == variable.TeacherID || assignment.RoomID == value.RoomID {
				blocking = append(blocking, assignment)
			}
			continue
		}

		// Travel time check for the same teacher on the same day
		if assignment.TeacherID != "" && assignment.TeacherID == variable.TeacherID && sameDay(assignment.StartTime, value.TimeSlot.Start) {
			var fromRoom, toRoom entities.Room
			var previousEnd, nextStart time.Time

			if assignment.EndTime.Before(value.TimeSlot.Start) || assignment.EndTime.Equal(value.TimeSlot.Start) {
				fromRoom = p.roomsByID[assignment.RoomID]
				toRoom = p.roomsByID[value.RoomID]
				previousEnd = assignment.EndTime
				nextStart = value.TimeSlot.Start
			} else {
				fromRoom = p.roomsByID[value.RoomID]
				toRoom = p.roomsByID[assignment.RoomID]
				previousEnd = value.TimeSlot.End
				nextStart = assignment.StartTime
			}

			if !HasSufficientTravelGap(previousEnd, nextStart, &fromRoom, &toRoom, p.travelMap) {
				blocking = append(blocking, assignment)
			}
		}
	}

	return blocking
}

func buildVariableMap(variables []Variable) map[string]Variable {
	indexed := make(map[string]Variable, len(variables))
	for _, variable := range variables {
		indexed[variable.ID] = variable
	}
	return indexed
}

func slotKeyFromDomain(value DomainValue) string {
	return fmt.Sprintf("%s|%s|%s", value.TimeSlot.Start.Format("2006-01-02T15:04"), value.TimeSlot.End.Format("2006-01-02T15:04"), value.TimeSlot.ShiftID)
}

func slotKeyFromAssignment(assignment PreviewAssignment) string {
	return fmt.Sprintf("%s|%s|%s", assignment.StartTime.Format("2006-01-02T15:04"), assignment.EndTime.Format("2006-01-02T15:04"), assignment.ShiftID)
}

func groupDomainsBySlot(values []DomainValue) map[string][]DomainValue {
	grouped := make(map[string][]DomainValue)
	for _, value := range values {
		key := slotKeyFromDomain(value)
		grouped[key] = append(grouped[key], value)
	}

	for key := range grouped {
		sort.Slice(grouped[key], func(i, j int) bool {
			if grouped[key][i].RoomCapacity == grouped[key][j].RoomCapacity {
				return grouped[key][i].RoomID < grouped[key][j].RoomID
			}
			return grouped[key][i].RoomCapacity < grouped[key][j].RoomCapacity
		})
	}

	return grouped
}

func buildConflictAdjacency(variables []Variable) map[string]map[string]struct{} {
	adjacency := make(map[string]map[string]struct{}, len(variables))
	for _, variable := range variables {
		adjacency[variable.ID] = make(map[string]struct{})
	}

	for i := 0; i < len(variables); i++ {
		for j := i + 1; j < len(variables); j++ {
			if variables[i].ClassID == variables[j].ClassID || variables[i].TeacherID == variables[j].TeacherID {
				adjacency[variables[i].ID][variables[j].ID] = struct{}{}
				adjacency[variables[j].ID][variables[i].ID] = struct{}{}
			}
		}
	}

	return adjacency
}

func cloneAssignments(assignments map[string]PreviewAssignment) map[string]PreviewAssignment {
	cloned := make(map[string]PreviewAssignment, len(assignments))
	for key, value := range assignments {
		cloned[key] = value
	}
	return cloned
}
