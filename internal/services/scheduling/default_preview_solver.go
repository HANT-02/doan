package scheduling

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"doan/internal/entities"
	skillservice "doan/internal/services/skills"
)

type legacyPreviewSolver struct{}

func NewLegacyPreviewSolver() SchedulingSolver {
	return &legacyPreviewSolver{}
}

func (s *legacyPreviewSolver) Key() string {
	return SolverKeyLegacyPreview
}

func (s *legacyPreviewSolver) Label() string {
	return "Bộ giải xem trước cũ"
}

func (s *legacyPreviewSolver) Solve(_ context.Context, input SolverInput) (*SolverOutput, error) {
	startedAt := time.Now()
	problem := prepareSchedulingProblem(input)
	telemetry := newSolverTelemetry(s.Key(), s.Label(), input, problem)
	solver := newBacktrackingSolver(&problem)
	assignments, solverConflicts := solver.Solve()

	assignmentsByID := make(map[string]PreviewAssignment, len(assignments))
	for _, assignment := range assignments {
		assignmentsByID[assignment.VariableID] = assignment
	}

	finishedAt := time.Now()
	return buildSolverOutput(
		input,
		problem.variables,
		assignmentsByID,
		problem.presetConflicts,
		problem.noDomainConflicts,
		solverConflicts,
		problem.targetLessons,
		finalizeSolverTelemetry(telemetry, startedAt, finishedAt),
	), nil
}

func buildVariables(input SolverInput) ([]Variable, []PreviewConflict) {
	teacherFilter := make(map[string]struct{}, len(input.TeacherIDs))
	for _, teacherID := range input.TeacherIDs {
		teacherFilter[teacherID] = struct{}{}
	}

	var (
		variables []Variable
		conflicts []PreviewConflict
	)

	for _, classEntity := range input.Classes {
		if classEntity.TeacherID == nil || *classEntity.TeacherID == "" {
			conflicts = append(conflicts, PreviewConflict{
				VariableID: classEntity.ID,
				ClassID:    classEntity.ID,
				ClassCode:  classEntity.Code,
				ClassName:  classEntity.Name,
				Type:       "MISSING_TEACHER",
				Message:    "Lớp chưa có giáo viên phụ trách nên không thể đưa vào preview.",
			})
			continue
		}

		if classEntity.CourseID == nil || *classEntity.CourseID == "" {
			conflicts = append(conflicts, PreviewConflict{
				VariableID: classEntity.ID,
				ClassID:    classEntity.ID,
				ClassCode:  classEntity.Code,
				ClassName:  classEntity.Name,
				Type:       "MISSING_COURSE",
				Message:    "Lớp chưa gắn khóa học nên chưa thể sinh số buổi học thực tế cho preview.",
			})
			continue
		}

		if classEntity.Course.SessionDurationMinutes <= 0 {
			conflicts = append(conflicts, PreviewConflict{
				VariableID: classEntity.ID,
				ClassID:    classEntity.ID,
				ClassCode:  classEntity.Code,
				ClassName:  classEntity.Name,
				Type:       "INVALID_COURSE_DURATION",
				Message:    "Khóa học của lớp chưa có `session_duration_minutes` hợp lệ, nên chưa thể sinh thời lượng buổi học chính xác.",
			})
			continue
		}

		if len(classEntity.ClassSchedules) == 0 {
			conflicts = append(conflicts, PreviewConflict{
				VariableID: classEntity.ID,
				ClassID:    classEntity.ID,
				ClassCode:  classEntity.Code,
				ClassName:  classEntity.Name,
				Type:       "MISSING_CLASS_SCHEDULE",
				Message:    "Lớp chưa có lịch tuần (`class_schedule`), nên chưa thể tính số buổi cần xếp từ thời gian học của lớp.",
			})
			continue
		}

		if classEntity.Course.SessionCount <= 0 {
			conflicts = append(conflicts, PreviewConflict{
				VariableID: classEntity.ID,
				ClassID:    classEntity.ID,
				ClassCode:  classEntity.Code,
				ClassName:  classEntity.Name,
				Type:       "INVALID_COURSE_SESSION_COUNT",
				Message:    "Khóa học của lớp chưa có tổng số buổi hợp lệ, nên chưa thể xác định đủ số buổi cần xếp.",
			})
			continue
		}

		missingRequiredSkills := skillservice.MissingRequiredCodes(
			[]string(classEntity.Teacher.Skills),
			[]string(classEntity.Course.RequiredSkills),
		)
		if len(missingRequiredSkills) > 0 {
			conflicts = append(conflicts, PreviewConflict{
				VariableID: classEntity.ID,
				ClassID:    classEntity.ID,
				ClassCode:  classEntity.Code,
				ClassName:  classEntity.Name,
				Type:       "SKILL_MISMATCH",
				Message: fmt.Sprintf(
					"Giáo viên phụ trách chưa đáp ứng kỹ năng/chứng chỉ bắt buộc của khóa học. Thiếu: %s.",
					strings.Join(missingRequiredSkills, ", "),
				),
			})
			continue
		}

		if len(teacherFilter) > 0 {
			if _, ok := teacherFilter[*classEntity.TeacherID]; !ok {
				continue
			}
		}

		teacherLabel := *classEntity.TeacherID
		if classEntity.Teacher.FullName != "" {
			teacherLabel = classEntity.Teacher.FullName
		} else if classEntity.Teacher.Code != "" {
			teacherLabel = classEntity.Teacher.Code
		}

		var preferredRoomID string
		if classEntity.RoomID != nil {
			preferredRoomID = *classEntity.RoomID
		}

		window := resolveClassSchedulingWindow(input, classEntity)
		if window.DateFrom.IsZero() || window.DateTo.IsZero() {
			conflicts = append(conflicts, PreviewConflict{
				VariableID: classEntity.ID,
				ClassID:    classEntity.ID,
				ClassCode:  classEntity.Code,
				ClassName:  classEntity.Name,
				Type:       "NO_SLOT_IN_RANGE",
				Message: fmt.Sprintf(
					"Không thể suy ra khoảng ngày xếp lịch hợp lệ từ ngày bắt đầu dự kiến và lịch tuần của lớp (%s). Hãy kiểm tra ngày bắt đầu lớp hoặc lịch tuần lớp.",
					describeScheduleDays(classEntity.ClassSchedules),
				),
			})
			continue
		}

		sessionTotal := window.SessionTotal
		slotCapacity := countUniqueTimeSlots(generateTimeSlotsForVariable(
			window.DateFrom,
			window.DateTo,
			classEntity.Course.SessionDurationMinutes,
			classEntity.ClassSchedules,
			input.Shifts,
		))
		if slotCapacity < sessionTotal {
			conflicts = append(conflicts, PreviewConflict{
				VariableID: classEntity.ID,
				ClassID:    classEntity.ID,
				ClassCode:  classEntity.Code,
				ClassName:  classEntity.Name,
				Type:       "INSUFFICIENT_SCHEDULE_SLOTS",
				Message: fmt.Sprintf(
					"Khoảng ngày dự kiến của lớp chỉ tạo được %d slot hợp lệ theo lịch tuần (%s), chưa đủ để phủ %d buổi của khóa học.",
					slotCapacity,
					describeScheduleDays(classEntity.ClassSchedules),
					sessionTotal,
				),
			})
			continue
		}

		for sessionIndex := 1; sessionIndex <= sessionTotal; sessionIndex++ {
			variables = append(variables, Variable{
				ID:              fmt.Sprintf("%s-session-%02d", classEntity.ID, sessionIndex),
				ClassID:         classEntity.ID,
				ClassCode:       classEntity.Code,
				ClassName:       classEntity.Name,
				SessionIndex:    sessionIndex,
				SessionTotal:    sessionTotal,
				TeacherID:       *classEntity.TeacherID,
				TeacherLabel:    teacherLabel,
				ExpectedCapcity: classEntity.MaxStudents,
				DurationMinutes: classEntity.Course.SessionDurationMinutes,
				PreferredRoomID: preferredRoomID,
			})
		}
	}

	return variables, conflicts
}

func buildDomains(
	variables []Variable,
	rooms []entities.Room,
	input SolverInput,
	classSchedulesByClass map[string][]entities.ClassSchedule,
	defaultShifts []entities.Shift,
) (map[string][]DomainValue, map[string]PreviewConflict) {
	domains := make(map[string][]DomainValue, len(variables))
	noDomainConflicts := make(map[string]PreviewConflict, len(variables))
	classUniqueSlotCapacity := make(map[string]int, len(classSchedulesByClass))

	for _, variable := range variables {
		classSchedules := classSchedulesByClass[variable.ClassID]
		if len(classSchedules) == 0 {
			domains[variable.ID] = []DomainValue{}
			noDomainConflicts[variable.ID] = PreviewConflict{
				VariableID:   variable.ID,
				ClassID:      variable.ClassID,
				ClassCode:    variable.ClassCode,
				ClassName:    variable.ClassName,
				SessionIndex: variable.SessionIndex,
				SessionTotal: variable.SessionTotal,
				Type:         "MISSING_CLASS_SCHEDULE",
				Message:      fmt.Sprintf("Lớp chưa có lịch tuần (`class_schedule`), nên chưa thể sinh buổi %d/%d cho preview. Hãy cấu hình ít nhất một ngày học và ca học cho lớp.", variable.SessionIndex, variable.SessionTotal),
			}
			continue
		}

		classEntity, ok := findClassByID(input.Classes, variable.ClassID)
		if !ok {
			domains[variable.ID] = []DomainValue{}
			noDomainConflicts[variable.ID] = PreviewConflict{
				VariableID:   variable.ID,
				ClassID:      variable.ClassID,
				ClassCode:    variable.ClassCode,
				ClassName:    variable.ClassName,
				SessionIndex: variable.SessionIndex,
				SessionTotal: variable.SessionTotal,
				Type:         "MISSING_CLASS",
				Message:      fmt.Sprintf("Không tìm thấy dữ liệu lớp để xếp buổi %d/%d.", variable.SessionIndex, variable.SessionTotal),
			}
			continue
		}

		window := resolveClassSchedulingWindow(input, classEntity)
		if window.DateFrom.IsZero() || window.DateTo.IsZero() {
			domains[variable.ID] = []DomainValue{}
			noDomainConflicts[variable.ID] = PreviewConflict{
				VariableID:   variable.ID,
				ClassID:      variable.ClassID,
				ClassCode:    variable.ClassCode,
				ClassName:    variable.ClassName,
				SessionIndex: variable.SessionIndex,
				SessionTotal: variable.SessionTotal,
				Type:         "NO_SLOT_IN_RANGE",
				Message:      fmt.Sprintf("Chưa suy ra được khoảng ngày hợp lệ để xếp buổi %d/%d của lớp.", variable.SessionIndex, variable.SessionTotal),
			}
			continue
		}

		slots := generateTimeSlotsForVariable(window.DateFrom, window.DateTo, variable.DurationMinutes, classSchedules, defaultShifts)
		if _, ok := classUniqueSlotCapacity[variable.ClassID]; !ok {
			classUniqueSlotCapacity[variable.ClassID] = countUniqueTimeSlots(slots)
		}
		if classUniqueSlotCapacity[variable.ClassID] > 0 && variable.SessionIndex > classUniqueSlotCapacity[variable.ClassID] {
			domains[variable.ID] = []DomainValue{}
			noDomainConflicts[variable.ID] = PreviewConflict{
				VariableID:   variable.ID,
				ClassID:      variable.ClassID,
				ClassCode:    variable.ClassCode,
				ClassName:    variable.ClassName,
				SessionIndex: variable.SessionIndex,
				SessionTotal: variable.SessionTotal,
				Type:         "INSUFFICIENT_SCHEDULE_SLOTS",
				Message: fmt.Sprintf(
					"Trong khoảng ngày đã chọn chỉ tạo được %d slot hợp lệ theo lịch tuần của lớp (%s), nên chưa đủ để xếp buổi %d/%d. Hãy nới khoảng ngày preview hoặc bổ sung thêm ngày/ca trong lịch tuần lớp.",
					classUniqueSlotCapacity[variable.ClassID],
					describeScheduleDays(classSchedules),
					variable.SessionIndex,
					variable.SessionTotal,
				),
			}
			continue
		}

		values := make([]DomainValue, 0)
		for _, room := range rooms {
			if variable.PreferredRoomID != "" && variable.PreferredRoomID != room.ID {
				continue
			}
			if room.Capacity < variable.ExpectedCapcity {
				continue
			}
			for _, slot := range slots {
				if slot.PreferredRoomID != "" && slot.PreferredRoomID != room.ID {
					continue
				}
				if slot.End.Hour() > 22 || (slot.End.Hour() == 22 && slot.End.Minute() > 0) {
					continue
				}
				values = append(values, DomainValue{
					RoomID:       room.ID,
					RoomName:     room.Name,
					RoomCapacity: room.Capacity,
					TimeSlot:     slot,
				})
			}
		}

		sort.Slice(values, func(i, j int) bool {
			if values[i].RoomID == variable.PreferredRoomID && values[j].RoomID != variable.PreferredRoomID {
				return true
			}
			return values[i].TimeSlot.Start.Before(values[j].TimeSlot.Start)
		})
		domains[variable.ID] = values
		if len(values) == 0 {
			noDomainConflicts[variable.ID] = explainNoDomain(variable, rooms, slots, len(classSchedules) > 0)
		}
	}

	return domains, noDomainConflicts
}

func generateTimeSlots(dateFrom, dateTo time.Time, durationMinutes int, shifts []entities.Shift) []TimeSlot {
	if dateTo.Before(dateFrom) {
		dateTo = dateFrom
	}

	slots := make([]TimeSlot, 0)
	for day := startOfDay(dateFrom); !day.After(startOfDay(dateTo)); day = day.AddDate(0, 0, 1) {
		for _, shift := range shifts {
			if !shift.IsActive {
				continue
			}

			startHour, startMinute, startOk := parseClockValue(shift.StartTime)
			endHour, endMinute, endOk := parseClockValue(shift.EndTime)
			if !startOk || !endOk {
				continue
			}

			start := time.Date(day.Year(), day.Month(), day.Day(), startHour, startMinute, 0, 0, day.Location())
			limitEnd := time.Date(day.Year(), day.Month(), day.Day(), endHour, endMinute, 0, 0, day.Location())
			end := start.Add(time.Duration(durationMinutes) * time.Minute)
			if !end.After(start) || end.After(limitEnd) {
				continue
			}

			slots = append(slots, TimeSlot{
				Start:     start,
				End:       end,
				ShiftID:   shift.ID,
				ShiftCode: shift.Code,
				ShiftName: shift.Name,
				ShiftType: shift.SessionType,
			})
		}
	}

	return slots
}

func generateTimeSlotsForVariable(
	dateFrom, dateTo time.Time,
	durationMinutes int,
	classSchedules []entities.ClassSchedule,
	defaultShifts []entities.Shift,
) []TimeSlot {
	if len(classSchedules) == 0 {
		return generateTimeSlots(dateFrom, dateTo, durationMinutes, defaultShifts)
	}

	slots := make([]TimeSlot, 0)
	for day := startOfDay(dateFrom); !day.After(startOfDay(dateTo)); day = day.AddDate(0, 0, 1) {
		for _, schedule := range classSchedules {
			if !matchesScheduleDay(day, schedule.DayOfWeek) {
				continue
			}

			if schedule.ShiftID == "" || schedule.Shift.ID == "" || !schedule.Shift.IsActive {
				continue
			}

			startHour, startMinute, startOk := parseClockValue(schedule.Shift.StartTime)
			endHour, endMinute, endOk := parseClockValue(schedule.Shift.EndTime)
			if !startOk || !endOk {
				continue
			}

			start := time.Date(day.Year(), day.Month(), day.Day(), startHour, startMinute, 0, 0, day.Location())
			limitEnd := time.Date(day.Year(), day.Month(), day.Day(), endHour, endMinute, 0, 0, day.Location())
			end := start.Add(time.Duration(durationMinutes) * time.Minute)

			if !end.After(start) || end.After(limitEnd) {
				continue
			}

			slot := TimeSlot{
				Start:     start,
				End:       end,
				ShiftID:   schedule.Shift.ID,
				ShiftCode: schedule.Shift.Code,
				ShiftName: schedule.Shift.Name,
				ShiftType: schedule.Shift.SessionType,
			}
			if schedule.RoomID != nil {
				slot.PreferredRoomID = *schedule.RoomID
			}

			slots = append(slots, slot)
		}
	}

	sort.Slice(slots, func(i, j int) bool {
		return slots[i].Start.Before(slots[j].Start)
	})

	return slots
}

func matchesScheduleDay(day time.Time, scheduleDay string) bool {
	switch strings.ToLower(strings.TrimSpace(scheduleDay)) {
	case "monday", "mon", "thu_hai", "thứ hai", "thu 2", "thứ 2", "2":
		return day.Weekday() == time.Monday
	case "tuesday", "tue", "tues", "thu_ba", "thứ ba", "thu 3", "thứ 3", "3":
		return day.Weekday() == time.Tuesday
	case "wednesday", "wed", "thu_tu", "thứ tư", "thu 4", "thứ 4", "4":
		return day.Weekday() == time.Wednesday
	case "thursday", "thu", "thur", "thurs", "thu_nam", "thứ năm", "thu 5", "thứ 5", "5":
		return day.Weekday() == time.Thursday
	case "friday", "fri", "thu_sau", "thứ sáu", "thu 6", "thứ 6", "6":
		return day.Weekday() == time.Friday
	case "saturday", "sat", "thu_bay", "thứ bảy", "thu 7", "thứ 7", "7":
		return day.Weekday() == time.Saturday
	case "sunday", "sun", "chu_nhat", "chủ nhật", "cn":
		return day.Weekday() == time.Sunday
	default:
		return false
	}
}

func countUniqueTimeSlots(slots []TimeSlot) int {
	if len(slots) == 0 {
		return 0
	}
	unique := make(map[string]struct{}, len(slots))
	for _, slot := range slots {
		key := fmt.Sprintf("%s|%s|%s", slot.Start.Format(time.RFC3339), slot.End.Format(time.RFC3339), slot.ShiftID)
		unique[key] = struct{}{}
	}
	return len(unique)
}

func describeScheduleDays(classSchedules []entities.ClassSchedule) string {
	if len(classSchedules) == 0 {
		return "chưa cấu hình ngày học"
	}
	labels := make([]string, 0, len(classSchedules))
	seen := make(map[string]struct{}, len(classSchedules))
	for _, schedule := range classSchedules {
		label := strings.TrimSpace(schedule.DayOfWeek)
		if label == "" {
			continue
		}
		normalized := strings.ToUpper(label)
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		labels = append(labels, normalized)
	}
	if len(labels) == 0 {
		return "chưa cấu hình ngày học"
	}
	sort.Strings(labels)
	return strings.Join(labels, ", ")
}

func parseClockValue(raw string) (hour int, minute int, ok bool) {
	for _, layout := range []string{"15:04", "15:04:05"} {
		parsed, err := time.Parse(layout, raw)
		if err == nil {
			return parsed.Hour(), parsed.Minute(), true
		}
	}

	return 0, 0, false
}

func startOfDay(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

func scoreAssignments(assignments []PreviewAssignment, targetLessons map[string]entities.Lesson) int {
	score := 0
	for _, assignment := range assignments {
		if target, ok := targetLessons[assignment.VariableID]; ok {
			targetTeacherID := ""
			if target.TeacherID != nil {
				targetTeacherID = *target.TeacherID
			}
			targetRoomID := ""
			if target.RoomID != nil {
				targetRoomID = *target.RoomID
			}

			if target.Status == entities.LessonStatusPublished {
				if !assignment.StartTime.Equal(target.DateStart) {
					score -= 1000
				}
				if assignment.TeacherID != targetTeacherID {
					score -= 500
				}
				if assignment.RoomID != targetRoomID {
					score -= 500
				}
			} else {
				if !assignment.StartTime.Equal(target.DateStart) {
					score -= 100
				}
				if assignment.TeacherID != targetTeacherID {
					score -= 50
				}
				if assignment.RoomID != targetRoomID {
					score -= 50
				}
			}
		}
	}

	for i := 0; i < len(assignments)-1; i++ {
		current := assignments[i]
		next := assignments[i+1]
		if current.TeacherID == next.TeacherID && sameDay(current.StartTime, next.StartTime) {
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

func sameDay(left, right time.Time) bool {
	return left.Year() == right.Year() && left.Month() == right.Month() && left.Day() == right.Day()
}

type backtrackingSolver struct {
	problem           *preparedSchedulingProblem
	variables         []Variable
	domains           map[string][]DomainValue
	noDomainConflicts map[string]PreviewConflict
}

func newBacktrackingSolver(problem *preparedSchedulingProblem) *backtrackingSolver {
	clonedDomains := make(map[string][]DomainValue, len(problem.domains))
	for key, values := range problem.domains {
		clonedDomains[key] = append([]DomainValue(nil), values...)
	}

	return &backtrackingSolver{
		problem:           problem,
		variables:         append([]Variable(nil), problem.variables...),
		domains:           clonedDomains,
		noDomainConflicts: problem.noDomainConflicts,
	}
}

func (s *backtrackingSolver) Solve() ([]PreviewAssignment, []PreviewConflict) {
	assignments := make(map[string]PreviewAssignment)
	orderedVariables := s.sortedVariables()
	conflicts := make([]PreviewConflict, 0)

	if s.backtrack(orderedVariables, assignments) {
		return assignmentsToSlice(assignments), conflicts
	}

	assignments = s.greedyAssign(orderedVariables)

	assignedKeys := make(map[string]struct{}, len(assignments))
	for key := range assignments {
		assignedKeys[key] = struct{}{}
	}

	for _, variable := range orderedVariables {
		if _, ok := assignedKeys[variable.ID]; ok {
			continue
		}

		conflictType := "NO_DOMAIN"
		message := "Không tìm thấy phương án hợp lệ với hard constraints hiện tại."
		if precomputedConflict, ok := s.noDomainConflicts[variable.ID]; ok {
			conflictType = precomputedConflict.Type
			message = precomputedConflict.Message
		} else if len(s.domains[variable.ID]) == 0 {
			message = "Không còn room/timeslot khả dụng cho lớp này trong khoảng ngày đã chọn."
		}

		conflicts = append(conflicts, PreviewConflict{
			VariableID:   variable.ID,
			ClassID:      variable.ClassID,
			ClassCode:    variable.ClassCode,
			ClassName:    variable.ClassName,
			SessionIndex: variable.SessionIndex,
			SessionTotal: variable.SessionTotal,
			Type:         conflictType,
			Message:      message,
		})
	}

	return assignmentsToSlice(assignments), conflicts
}

func (s *backtrackingSolver) greedyAssign(variables []Variable) map[string]PreviewAssignment {
	assignments := make(map[string]PreviewAssignment)
	for _, variable := range variables {
		for _, value := range s.domains[variable.ID] {
			if !s.isConsistent(variable, value, assignments) {
				continue
			}

			assignments[variable.ID] = newPreviewAssignment(variable, value, "HARD_OK_PARTIAL")
			break
		}
	}
	return assignments
}

func (s *backtrackingSolver) sortedVariables() []Variable {
	ordered := append([]Variable(nil), s.variables...)
	sort.Slice(ordered, func(i, j int) bool {
		return len(s.domains[ordered[i].ID]) < len(s.domains[ordered[j].ID])
	})
	return ordered
}

func (s *backtrackingSolver) backtrack(variables []Variable, assignments map[string]PreviewAssignment) bool {
	if len(assignments) == len(variables) {
		return true
	}

	variable, ok := s.selectUnassignedVariable(variables, assignments)
	if !ok {
		return true
	}

	for _, value := range s.domains[variable.ID] {
		if !s.isConsistent(variable, value, assignments) {
			continue
		}

		assignments[variable.ID] = newPreviewAssignment(variable, value, "HARD_OK")

		originalDomains, forwardOK := s.forwardCheck(variable, assignments)
		if forwardOK && s.backtrack(variables, assignments) {
			return true
		}

		s.restoreDomains(originalDomains)
		delete(assignments, variable.ID)
	}

	return false
}

func (s *backtrackingSolver) selectUnassignedVariable(variables []Variable, assignments map[string]PreviewAssignment) (Variable, bool) {
	var (
		selected Variable
		found    bool
	)
	bestDomainSize := int(^uint(0) >> 1)

	for _, variable := range variables {
		if _, ok := assignments[variable.ID]; ok {
			continue
		}

		domainSize := len(s.domains[variable.ID])
		if !found || domainSize < bestDomainSize {
			selected = variable
			bestDomainSize = domainSize
			found = true
		}
	}

	return selected, found
}

func (s *backtrackingSolver) isConsistent(variable Variable, value DomainValue, assignments map[string]PreviewAssignment) bool {
	if value.TimeSlot.End.Hour() > 22 || (value.TimeSlot.End.Hour() == 22 && value.TimeSlot.End.Minute() > 0) {
		return false
	}

	return !s.problem.hasConflict(variable, value, assignments)
}

func (s *backtrackingSolver) forwardCheck(variable Variable, assignments map[string]PreviewAssignment) (map[string][]DomainValue, bool) {
	original := make(map[string][]DomainValue)
	for _, other := range s.variables {
		if other.ID == variable.ID {
			continue
		}
		if _, assigned := assignments[other.ID]; assigned {
			continue
		}

		original[other.ID] = append([]DomainValue(nil), s.domains[other.ID]...)
		filtered := make([]DomainValue, 0, len(s.domains[other.ID]))
		for _, candidate := range s.domains[other.ID] {
			if s.isConsistent(other, candidate, assignments) {
				filtered = append(filtered, candidate)
			}
		}
		s.domains[other.ID] = filtered
		if len(filtered) == 0 {
			return original, false
		}
	}
	return original, true
}

func (s *backtrackingSolver) restoreDomains(original map[string][]DomainValue) {
	for key, values := range original {
		s.domains[key] = values
	}
}

func assignmentsToSlice(assignments map[string]PreviewAssignment) []PreviewAssignment {
	items := make([]PreviewAssignment, 0, len(assignments))
	for _, assignment := range assignments {
		items = append(items, assignment)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].StartTime.Equal(items[j].StartTime) {
			return fmt.Sprintf("%s-%s", items[i].ClassCode, items[i].RoomID) < fmt.Sprintf("%s-%s", items[j].ClassCode, items[j].RoomID)
		}
		return items[i].StartTime.Before(items[j].StartTime)
	})
	return items
}

func overlaps(startA, endA, startB, endB time.Time) bool {
	return startA.Before(endB) && startB.Before(endA)
}

func explainNoDomain(variable Variable, rooms []entities.Room, slots []TimeSlot, hasClassSchedule bool) PreviewConflict {
	if len(slots) == 0 {
		if hasClassSchedule {
			return PreviewConflict{
				VariableID:   variable.ID,
				ClassID:      variable.ClassID,
				ClassCode:    variable.ClassCode,
				ClassName:    variable.ClassName,
				SessionIndex: variable.SessionIndex,
				SessionTotal: variable.SessionTotal,
				Type:         "CLASS_SCHEDULE_NO_SLOT",
				Message:      fmt.Sprintf("Lịch mẫu của lớp không tạo được slot hợp lệ cho buổi %d/%d trong khoảng ngày đã chọn. Hãy kiểm tra `day_of_week`, `shift_id` hoặc nới khoảng ngày preview.", variable.SessionIndex, variable.SessionTotal),
			}
		}

		return PreviewConflict{
			VariableID:   variable.ID,
			ClassID:      variable.ClassID,
			ClassCode:    variable.ClassCode,
			ClassName:    variable.ClassName,
			SessionIndex: variable.SessionIndex,
			SessionTotal: variable.SessionTotal,
			Type:         "NO_ACTIVE_SHIFT",
			Message:      fmt.Sprintf("Không có ca học khả dụng để sinh slot cho buổi %d/%d của lớp này. Hãy tạo hoặc kích hoạt `Shift` trước khi chạy preview.", variable.SessionIndex, variable.SessionTotal),
		}
	}

	if len(rooms) == 0 {
		return PreviewConflict{
			VariableID:   variable.ID,
			ClassID:      variable.ClassID,
			ClassCode:    variable.ClassCode,
			ClassName:    variable.ClassName,
			SessionIndex: variable.SessionIndex,
			SessionTotal: variable.SessionTotal,
			Type:         "NO_ACTIVE_ROOM",
			Message:      fmt.Sprintf("Không có phòng khả dụng để xếp buổi %d/%d của lớp này. Hãy kiểm tra bộ lọc phòng hoặc dữ liệu phòng đang hoạt động.", variable.SessionIndex, variable.SessionTotal),
		}
	}

	requiredSlotRooms := collectRequiredSlotRooms(slots)
	if len(requiredSlotRooms) > 0 && !containsAnyRoom(rooms, requiredSlotRooms) {
		return PreviewConflict{
			VariableID:   variable.ID,
			ClassID:      variable.ClassID,
			ClassCode:    variable.ClassCode,
			ClassName:    variable.ClassName,
			SessionIndex: variable.SessionIndex,
			SessionTotal: variable.SessionTotal,
			Type:         "CLASS_SCHEDULE_ROOM_UNAVAILABLE",
			Message:      fmt.Sprintf("Lịch mẫu của lớp đang khóa buổi %d/%d vào phòng cụ thể theo ca học đã chọn, nhưng phòng đó không nằm trong tập phòng khả dụng hiện tại. Hãy kiểm tra `class_schedule.room_id` hoặc bỏ bớt bộ lọc phòng.", variable.SessionIndex, variable.SessionTotal),
		}
	}

	if variable.PreferredRoomID != "" {
		for _, room := range rooms {
			if room.ID != variable.PreferredRoomID {
				continue
			}

			if room.Capacity < variable.ExpectedCapcity {
				return PreviewConflict{
					VariableID:   variable.ID,
					ClassID:      variable.ClassID,
					ClassCode:    variable.ClassCode,
					ClassName:    variable.ClassName,
					SessionIndex: variable.SessionIndex,
					SessionTotal: variable.SessionTotal,
					Type:         "ROOM_CAPACITY_BLOCK",
					Message:      fmt.Sprintf("Phòng đang gán sẵn cho buổi %d/%d chỉ chứa %d chỗ, nhỏ hơn sĩ số tối đa %d. Hãy đổi phòng hoặc giảm sĩ số tối đa.", variable.SessionIndex, variable.SessionTotal, room.Capacity, variable.ExpectedCapcity),
				}
			}
		}

		return PreviewConflict{
			VariableID:   variable.ID,
			ClassID:      variable.ClassID,
			ClassCode:    variable.ClassCode,
			ClassName:    variable.ClassName,
			SessionIndex: variable.SessionIndex,
			SessionTotal: variable.SessionTotal,
			Type:         "PREFERRED_ROOM_UNAVAILABLE",
			Message:      fmt.Sprintf("Buổi %d/%d đang gán vào một phòng không nằm trong tập phòng khả dụng hiện tại. Hãy bỏ lọc phòng hoặc gán lại phòng học.", variable.SessionIndex, variable.SessionTotal),
		}
	}

	for _, room := range rooms {
		if room.Capacity >= variable.ExpectedCapcity {
			return PreviewConflict{
				VariableID:   variable.ID,
				ClassID:      variable.ClassID,
				ClassCode:    variable.ClassCode,
				ClassName:    variable.ClassName,
				SessionIndex: variable.SessionIndex,
				SessionTotal: variable.SessionTotal,
				Type:         "NO_DOMAIN",
				Message:      fmt.Sprintf("Có phòng phù hợp nhưng không tìm được tổ hợp phòng/khung giờ cho buổi %d/%d thỏa hard constraints hiện tại. Hãy nới bộ lọc hoặc đổi khoảng ngày preview.", variable.SessionIndex, variable.SessionTotal),
			}
		}
	}

	return PreviewConflict{
		VariableID:   variable.ID,
		ClassID:      variable.ClassID,
		ClassCode:    variable.ClassCode,
		ClassName:    variable.ClassName,
		SessionIndex: variable.SessionIndex,
		SessionTotal: variable.SessionTotal,
		Type:         "ROOM_CAPACITY_BLOCK",
		Message:      fmt.Sprintf("Không có phòng khả dụng nào đủ sức chứa cho buổi %d/%d của lớp này. Sĩ số tối đa hiện tại là %d học viên.", variable.SessionIndex, variable.SessionTotal, variable.ExpectedCapcity),
	}
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}

func indexClassSchedules(classes []entities.Class) map[string][]entities.ClassSchedule {
	indexed := make(map[string][]entities.ClassSchedule, len(classes))
	for _, classEntity := range classes {
		if len(classEntity.ClassSchedules) == 0 {
			continue
		}

		indexed[classEntity.ID] = append([]entities.ClassSchedule(nil), classEntity.ClassSchedules...)
	}

	return indexed
}

func findClassByID(classes []entities.Class, classID string) (entities.Class, bool) {
	for _, classEntity := range classes {
		if classEntity.ID == classID {
			return classEntity, true
		}
	}
	return entities.Class{}, false
}

func collectRequiredSlotRooms(slots []TimeSlot) map[string]struct{} {
	requiredRooms := make(map[string]struct{})
	for _, slot := range slots {
		if slot.PreferredRoomID == "" {
			continue
		}

		requiredRooms[slot.PreferredRoomID] = struct{}{}
	}

	return requiredRooms
}

func containsAnyRoom(rooms []entities.Room, candidates map[string]struct{}) bool {
	for _, room := range rooms {
		if _, ok := candidates[room.ID]; ok {
			return true
		}
	}

	return false
}
