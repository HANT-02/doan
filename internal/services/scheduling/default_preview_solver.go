package scheduling

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"doan/internal/entities"
)

type legacyPreviewSolver struct{}

func NewLegacyPreviewSolver() SchedulingSolver {
	return &legacyPreviewSolver{}
}

func (s *legacyPreviewSolver) Key() string {
	return SolverKeyLegacyPreview
}

func (s *legacyPreviewSolver) Label() string {
	return "Legacy Preview Solver"
}

func (s *legacyPreviewSolver) Solve(_ context.Context, input SolverInput) (*SolverOutput, error) {
	problem := prepareSchedulingProblem(input)
	solver := newBacktrackingSolver(problem.variables, problem.domains, problem.noDomainConflicts)
	assignments, solverConflicts := solver.Solve()

	assignmentsByID := make(map[string]PreviewAssignment, len(assignments))
	for _, assignment := range assignments {
		assignmentsByID[assignment.VariableID] = assignment
	}

	return buildSolverOutput(input, problem.variables, assignmentsByID, problem.presetConflicts, problem.noDomainConflicts, solverConflicts), nil
}

func buildVariables(classes []entities.Class, teacherIDs []string) ([]Variable, []PreviewConflict) {
	teacherFilter := make(map[string]struct{}, len(teacherIDs))
	for _, teacherID := range teacherIDs {
		teacherFilter[teacherID] = struct{}{}
	}

	var (
		variables []Variable
		conflicts []PreviewConflict
	)

	for _, classEntity := range classes {
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

		if classEntity.Course.SessionCount <= 0 {
			conflicts = append(conflicts, PreviewConflict{
				VariableID: classEntity.ID,
				ClassID:    classEntity.ID,
				ClassCode:  classEntity.Code,
				ClassName:  classEntity.Name,
				Type:       "INVALID_COURSE_SESSION_COUNT",
				Message:    "Khóa học của lớp chưa có `session_count` hợp lệ, nên chưa thể sinh đủ số buổi cho preview.",
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

		for sessionIndex := 1; sessionIndex <= classEntity.Course.SessionCount; sessionIndex++ {
			variables = append(variables, Variable{
				ID:              fmt.Sprintf("%s-session-%02d", classEntity.ID, sessionIndex),
				ClassID:         classEntity.ID,
				ClassCode:       classEntity.Code,
				ClassName:       classEntity.Name,
				SessionIndex:    sessionIndex,
				SessionTotal:    classEntity.Course.SessionCount,
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

	for _, variable := range variables {
		classSchedules := classSchedulesByClass[variable.ClassID]
		slots := generateTimeSlotsForVariable(input.DateFrom, input.DateTo, variable.DurationMinutes, classSchedules, defaultShifts)
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

func scoreAssignments(assignments []PreviewAssignment) int {
	score := 0
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
	variables         []Variable
	domains           map[string][]DomainValue
	noDomainConflicts map[string]PreviewConflict
}

func newBacktrackingSolver(
	variables []Variable,
	domains map[string][]DomainValue,
	noDomainConflicts map[string]PreviewConflict,
) *backtrackingSolver {
	clonedDomains := make(map[string][]DomainValue, len(domains))
	for key, values := range domains {
		clonedDomains[key] = append([]DomainValue(nil), values...)
	}

	return &backtrackingSolver{
		variables:         append([]Variable(nil), variables...),
		domains:           clonedDomains,
		noDomainConflicts: noDomainConflicts,
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

	return !hasConflict(variable, value, assignments)
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
