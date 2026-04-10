package scheduling

import (
	"context"
	"sort"
)

type cpSatSolver struct {
	maxNodes int
}

type cpSearchState struct {
	bestAssignments map[string]PreviewAssignment
	bestAssigned    int
	bestSoftScore   int
	nodesVisited    int
}

func NewCPSATSolver() *cpSatSolver {
	return &cpSatSolver{maxNodes: 50000}
}

func (s *cpSatSolver) Key() string {
	return SolverKeyCPSAT
}

func (s *cpSatSolver) Label() string {
	return "CP-SAT"
}

func (s *cpSatSolver) Solve(_ context.Context, input SolverInput) (*SolverOutput, error) {
	problem := prepareSchedulingProblem(input)
	state := &cpSearchState{
		bestAssignments: make(map[string]PreviewAssignment),
		bestAssigned:    -1,
		bestSoftScore:   -1,
	}

	ordered := append([]Variable(nil), problem.variables...)
	sort.Slice(ordered, func(i, j int) bool {
		leftDomain := len(problem.domains[ordered[i].ID])
		rightDomain := len(problem.domains[ordered[j].ID])
		if leftDomain == rightDomain {
			return ordered[i].SessionTotal > ordered[j].SessionTotal
		}
		return leftDomain < rightDomain
	})

	s.search(ordered, problem.domains, 0, make(map[string]PreviewAssignment), state)

	return buildSolverOutput(input, problem.variables, state.bestAssignments, problem.presetConflicts, problem.noDomainConflicts, nil), nil
}

func (s *cpSatSolver) search(
	variables []Variable,
	domains map[string][]DomainValue,
	index int,
	assignments map[string]PreviewAssignment,
	state *cpSearchState,
) {
	if state.nodesVisited >= s.maxNodes {
		return
	}
	state.nodesVisited++

	remaining := len(variables) - index
	if len(assignments)+remaining < state.bestAssigned {
		return
	}

	if index >= len(variables) {
		s.maybeUpdateBest(assignments, state)
		return
	}

	variable := variables[index]
	candidates := append([]DomainValue(nil), domains[variable.ID]...)
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].TimeSlot.Start.Equal(candidates[j].TimeSlot.Start) {
			return candidates[i].RoomCapacity < candidates[j].RoomCapacity
		}
		return candidates[i].TimeSlot.Start.Before(candidates[j].TimeSlot.Start)
	})

	for _, candidate := range candidates {
		if hasConflict(variable, candidate, assignments) {
			continue
		}

		assignments[variable.ID] = newPreviewAssignment(variable, candidate, "CP_SAT_OK")
		s.search(variables, domains, index+1, assignments, state)
		delete(assignments, variable.ID)
	}

	s.search(variables, domains, index+1, assignments, state)
}

func (s *cpSatSolver) maybeUpdateBest(assignments map[string]PreviewAssignment, state *cpSearchState) {
	assignmentCount := len(assignments)
	softScore := scoreAssignments(assignmentsToSlice(assignments))

	if assignmentCount > state.bestAssigned || (assignmentCount == state.bestAssigned && softScore > state.bestSoftScore) {
		state.bestAssigned = assignmentCount
		state.bestSoftScore = softScore
		state.bestAssignments = cloneAssignments(assignments)
	}
}
