package scheduling

import (
	"context"
	"fmt"
	"sort"
)

type tabuSearchSolver struct {
	maxIterations int
	tabuTenure    int
}

type tabuMove struct {
	variableID string
	slotKey    string
	roomID     string
}

type tabuState struct {
	assignments map[string]PreviewAssignment
	penalty     int
}

func NewTabuSearchSolver() *tabuSearchSolver {
	return &tabuSearchSolver{
		maxIterations: 120,
		tabuTenure:    9,
	}
}

func (s *tabuSearchSolver) Key() string {
	return SolverKeyTabuSearch
}

func (s *tabuSearchSolver) Label() string {
	return "Tabu Search"
}

func (s *tabuSearchSolver) Solve(_ context.Context, input SolverInput) (*SolverOutput, error) {
	problem := prepareSchedulingProblem(input)
	initial := s.buildInitialState(problem.variables, problem.domains)
	best := tabuState{
		assignments: cloneAssignments(initial),
		penalty:     s.evaluate(initial, problem.variables),
	}
	current := tabuState{
		assignments: cloneAssignments(initial),
		penalty:     best.penalty,
	}

	tabu := make(map[string]int)
	variableMap := buildVariableMap(problem.variables)
	for iteration := 0; iteration < s.maxIterations; iteration++ {
		bestMoveFound := false
		var nextAssignments map[string]PreviewAssignment
		nextPenalty := 0
		nextMoveKey := ""

		for _, variable := range problem.variables {
			for _, candidate := range limitedNeighborhood(problem.domains[variable.ID], 6) {
				candidateAssignments := cloneAssignments(current.assignments)
				for _, blocking := range findBlockingAssignments(variable, candidate, candidateAssignments) {
					delete(candidateAssignments, blocking.VariableID)
				}
				candidateAssignments[variable.ID] = newPreviewAssignment(variable, candidate, "TABU_OK")

				penalty := s.evaluate(candidateAssignments, problem.variables)
				moveKey := formatTabuMove(variable.ID, candidate)
				tabuExpiry, isTabu := tabu[moveKey]
				if isTabu && tabuExpiry > iteration && penalty >= best.penalty {
					continue
				}

				if !bestMoveFound || penalty < nextPenalty {
					bestMoveFound = true
					nextAssignments = candidateAssignments
					nextPenalty = penalty
					nextMoveKey = moveKey
				}
			}
		}

		if !bestMoveFound {
			break
		}

		current.assignments = nextAssignments
		current.penalty = nextPenalty
		tabu[nextMoveKey] = iteration + s.tabuTenure

		if current.penalty < best.penalty {
			best.assignments = cloneAssignments(current.assignments)
			best.penalty = current.penalty
		}
	}

	repaired := repairAssignments(problem.variables, variableMap, problem.domains, best.assignments)
	return buildSolverOutput(input, problem.variables, repaired, problem.presetConflicts, problem.noDomainConflicts, nil), nil
}

func (s *tabuSearchSolver) buildInitialState(
	variables []Variable,
	domains map[string][]DomainValue,
) map[string]PreviewAssignment {
	assignments := make(map[string]PreviewAssignment)
	ordered := append([]Variable(nil), variables...)
	sort.Slice(ordered, func(i, j int) bool {
		return len(domains[ordered[i].ID]) < len(domains[ordered[j].ID])
	})

	for _, variable := range ordered {
		for _, candidate := range domains[variable.ID] {
			if hasConflict(variable, candidate, assignments) {
				continue
			}
			assignments[variable.ID] = newPreviewAssignment(variable, candidate, "TABU_INITIAL")
			break
		}
	}

	return assignments
}

func (s *tabuSearchSolver) evaluate(assignments map[string]PreviewAssignment, variables []Variable) int {
	penalty := (len(variables) - len(assignments)) * 1000
	items := assignmentsToSlice(assignments)
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if !overlaps(items[i].StartTime, items[i].EndTime, items[j].StartTime, items[j].EndTime) {
				continue
			}
			if items[i].ClassID == items[j].ClassID || items[i].TeacherID == items[j].TeacherID || items[i].RoomID == items[j].RoomID {
				penalty += 500
			}
		}
	}

	penalty -= scoreAssignments(items)
	return penalty
}

func limitedNeighborhood(values []DomainValue, limit int) []DomainValue {
	if len(values) <= limit {
		return append([]DomainValue(nil), values...)
	}
	return append([]DomainValue(nil), values[:limit]...)
}

func repairAssignments(
	variables []Variable,
	variableMap map[string]Variable,
	domains map[string][]DomainValue,
	assignments map[string]PreviewAssignment,
) map[string]PreviewAssignment {
	repaired := make(map[string]PreviewAssignment)
	items := assignmentsToSlice(assignments)
	for _, assignment := range items {
		variable, ok := variableMap[assignment.VariableID]
		if !ok {
			continue
		}
		domainValue, ok := findDomainValueForAssignment(domains[assignment.VariableID], assignment)
		if !ok || hasConflict(variable, domainValue, repaired) {
			continue
		}
		repaired[assignment.VariableID] = newPreviewAssignment(variable, domainValue, "TABU_OK")
	}

	for _, variable := range variables {
		if _, ok := repaired[variable.ID]; ok {
			continue
		}
		for _, candidate := range domains[variable.ID] {
			if hasConflict(variable, candidate, repaired) {
				continue
			}
			repaired[variable.ID] = newPreviewAssignment(variable, candidate, "TABU_REPAIRED")
			break
		}
	}

	return repaired
}

func findDomainValueForAssignment(values []DomainValue, assignment PreviewAssignment) (DomainValue, bool) {
	for _, value := range values {
		if value.RoomID == assignment.RoomID && value.TimeSlot.Start.Equal(assignment.StartTime) && value.TimeSlot.End.Equal(assignment.EndTime) {
			return value, true
		}
	}
	return DomainValue{}, false
}

func formatTabuMove(variableID string, candidate DomainValue) string {
	move := tabuMove{
		variableID: variableID,
		slotKey:    slotKeyFromDomain(candidate),
		roomID:     candidate.RoomID,
	}

	return fmt.Sprintf("%s|%s|%s", move.variableID, move.slotKey, move.roomID)
}
