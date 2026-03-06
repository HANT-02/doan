package scheduling

import (
	"context"
	"fmt"
	"sort"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	schedulingstore "doan/internal/services/scheduling"
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
	store     schedulingstore.PreviewStore[PreviewResult]
}

func NewPreviewUseCase(
	classRepo repositoryinterface.ClassRepository,
	roomRepo repositoryinterface.RoomRepository,
	store schedulingstore.PreviewStore[PreviewResult],
) PreviewUseCase {
	return &previewUseCase{
		classRepo: classRepo,
		roomRepo:  roomRepo,
		store:     store,
	}
}

func (uc *previewUseCase) Execute(ctx context.Context, input PreviewInput) (*PreviewResult, error) {
	ctxLogger := logger.NewLogger(ctx)

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

	variables, presetConflicts := buildVariables(classes, input.TeacherIDs)
	result.Conflicts = append(result.Conflicts, presetConflicts...)

	domains := buildDomains(variables, rooms, input)
	solver := newSolver(variables, domains)
	assignments, solverConflicts := solver.Solve()
	result.Conflicts = append(result.Conflicts, solverConflicts...)
	result.Assignments = assignments
	result.Summary = PreviewSummary{
		RequestedClasses:   len(classes),
		ScheduledLessons:   len(assignments),
		UnscheduledLessons: len(classes) - len(assignments),
		ConflictCount:      len(result.Conflicts),
		SoftScore:          scoreAssignments(assignments),
	}

	switch {
	case len(assignments) == 0:
		result.Status = "FAILED"
	case len(result.Conflicts) > 0:
		result.Status = "PARTIAL"
	default:
		result.Status = "COMPLETED"
	}

	uc.store.Save(runID, result)
	return &result, nil
}

func (uc *previewUseCase) loadClasses(ctx context.Context, input PreviewInput) ([]entities.Class, error) {
	condition := repositories.NewCommonCondition()
	condition.SetPaging(500, 1)
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

func (uc *previewUseCase) loadRooms(ctx context.Context, input PreviewInput) ([]entities.Room, error) {
	condition := repositories.NewCommonCondition()
	condition.SetPaging(500, 1)
	condition.AddCondition("status", "ACTIVE", repositories.Equal)
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

		if len(teacherFilter) > 0 {
			if _, ok := teacherFilter[*classEntity.TeacherID]; !ok {
				continue
			}
		}

		teacherLabel := *classEntity.TeacherID
		var preferredRoomID string
		if classEntity.RoomID != nil {
			preferredRoomID = *classEntity.RoomID
		}

		variables = append(variables, Variable{
			ID:              classEntity.ID,
			ClassID:         classEntity.ID,
			ClassCode:       classEntity.Code,
			ClassName:       classEntity.Name,
			TeacherID:       *classEntity.TeacherID,
			TeacherLabel:    teacherLabel,
			ExpectedCapcity: classEntity.MaxStudents,
			DurationMinutes: 120,
			PreferredRoomID: preferredRoomID,
		})
	}

	return variables, conflicts
}

func buildDomains(variables []Variable, rooms []entities.Room, input PreviewInput) map[string][]DomainValue {
	slots := generateTimeSlots(input.DateFrom, input.DateTo, 120)
	domains := make(map[string][]DomainValue, len(variables))

	for _, variable := range variables {
		values := make([]DomainValue, 0)
		for _, room := range rooms {
			if variable.PreferredRoomID != "" && variable.PreferredRoomID != room.ID {
				continue
			}
			if room.Capacity < variable.ExpectedCapcity {
				continue
			}
			for _, slot := range slots {
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
	}

	return domains
}

func generateTimeSlots(dateFrom, dateTo time.Time, durationMinutes int) []TimeSlot {
	if dateTo.Before(dateFrom) {
		dateTo = dateFrom
	}

	baseSlots := []struct {
		hour   int
		minute int
	}{
		{8, 0},
		{10, 15},
		{13, 30},
		{15, 45},
		{18, 0},
		{20, 0},
	}

	slots := make([]TimeSlot, 0)
	for day := startOfDay(dateFrom); !day.After(startOfDay(dateTo)); day = day.AddDate(0, 0, 1) {
		for _, baseSlot := range baseSlots {
			start := time.Date(day.Year(), day.Month(), day.Day(), baseSlot.hour, baseSlot.minute, 0, 0, day.Location())
			end := start.Add(time.Duration(durationMinutes) * time.Minute)
			slots = append(slots, TimeSlot{Start: start, End: end})
		}
	}

	return slots
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

type solver struct {
	variables []Variable
	domains   map[string][]DomainValue
}

func newSolver(variables []Variable, domains map[string][]DomainValue) *solver {
	clonedDomains := make(map[string][]DomainValue, len(domains))
	for key, values := range domains {
		clonedDomains[key] = append([]DomainValue(nil), values...)
	}

	return &solver{
		variables: append([]Variable(nil), variables...),
		domains:   clonedDomains,
	}
}

func (s *solver) Solve() ([]PreviewAssignment, []PreviewConflict) {
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
		if len(s.domains[variable.ID]) == 0 {
			message = "Không còn room/timeslot khả dụng cho lớp này trong khoảng ngày đã chọn."
		}

		conflicts = append(conflicts, PreviewConflict{
			VariableID: variable.ID,
			ClassID:    variable.ClassID,
			ClassCode:  variable.ClassCode,
			ClassName:  variable.ClassName,
			Type:       conflictType,
			Message:    message,
		})
	}

	return assignmentsToSlice(assignments), conflicts
}

func (s *solver) greedyAssign(variables []Variable) map[string]PreviewAssignment {
	assignments := make(map[string]PreviewAssignment)
	for _, variable := range variables {
		for _, value := range s.domains[variable.ID] {
			if !s.isConsistent(variable, value, assignments) {
				continue
			}

			assignments[variable.ID] = PreviewAssignment{
				VariableID:    variable.ID,
				ClassID:       variable.ClassID,
				ClassCode:     variable.ClassCode,
				ClassName:     variable.ClassName,
				TeacherID:     variable.TeacherID,
				TeacherLabel:  variable.TeacherLabel,
				RoomID:        value.RoomID,
				RoomName:      value.RoomName,
				RoomCapacity:  value.RoomCapacity,
				StartTime:     value.TimeSlot.Start,
				EndTime:       value.TimeSlot.End,
				ConstraintFit: "HARD_OK_PARTIAL",
			}
			break
		}
	}
	return assignments
}

func (s *solver) sortedVariables() []Variable {
	ordered := append([]Variable(nil), s.variables...)
	sort.Slice(ordered, func(i, j int) bool {
		return len(s.domains[ordered[i].ID]) < len(s.domains[ordered[j].ID])
	})
	return ordered
}

func (s *solver) backtrack(variables []Variable, assignments map[string]PreviewAssignment) bool {
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

		assignments[variable.ID] = PreviewAssignment{
			VariableID:    variable.ID,
			ClassID:       variable.ClassID,
			ClassCode:     variable.ClassCode,
			ClassName:     variable.ClassName,
			TeacherID:     variable.TeacherID,
			TeacherLabel:  variable.TeacherLabel,
			RoomID:        value.RoomID,
			RoomName:      value.RoomName,
			RoomCapacity:  value.RoomCapacity,
			StartTime:     value.TimeSlot.Start,
			EndTime:       value.TimeSlot.End,
			ConstraintFit: "HARD_OK",
		}

		originalDomains, forwardOK := s.forwardCheck(variable, assignments)
		if forwardOK && s.backtrack(variables, assignments) {
			return true
		}

		s.restoreDomains(originalDomains)
		delete(assignments, variable.ID)
	}

	return false
}

func (s *solver) selectUnassignedVariable(variables []Variable, assignments map[string]PreviewAssignment) (Variable, bool) {
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

func (s *solver) isConsistent(variable Variable, value DomainValue, assignments map[string]PreviewAssignment) bool {
	if value.RoomCapacity < variable.ExpectedCapcity {
		return false
	}
	if value.TimeSlot.End.Hour() > 22 || (value.TimeSlot.End.Hour() == 22 && value.TimeSlot.End.Minute() > 0) {
		return false
	}

	for _, assignment := range assignments {
		if overlaps(value.TimeSlot.Start, value.TimeSlot.End, assignment.StartTime, assignment.EndTime) {
			if assignment.TeacherID == variable.TeacherID {
				return false
			}
			if assignment.RoomID == value.RoomID {
				return false
			}
		}
	}

	return true
}

func (s *solver) forwardCheck(variable Variable, assignments map[string]PreviewAssignment) (map[string][]DomainValue, bool) {
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

func (s *solver) restoreDomains(original map[string][]DomainValue) {
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
