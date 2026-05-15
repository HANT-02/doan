package scheduling

import (
	"context"
	"sort"
	"time"

	"doan/internal/entities"
)

type cpSatSolver struct {
	maxNodes int
}

type cpSearchState struct {
	bestAssignments map[string]PreviewAssignment
	bestAssigned    int
	bestSoftScore   int
	nodesVisited    int
	targetLessons   map[string]entities.Lesson
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
	startedAt := time.Now()
	problem := prepareSchedulingProblem(input)
	telemetry := newSolverTelemetry(s.Key(), s.Label(), input, problem)
	state := &cpSearchState{
		bestAssignments: make(map[string]PreviewAssignment),
		bestAssigned:    -1,
		bestSoftScore:   -1,
		targetLessons:   problem.targetLessons,
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

	s.search(&problem, ordered, problem.domains, 0, make(map[string]PreviewAssignment), state, telemetry)

	telemetry.NodesVisited = state.nodesVisited
	finishedAt := time.Now()
	return buildSolverOutput(
		input,
		problem.variables,
		state.bestAssignments,
		problem.presetConflicts,
		problem.noDomainConflicts,
		nil,
		problem.targetLessons,
		finalizeSolverTelemetry(telemetry, startedAt, finishedAt),
	), nil
}

func (s *cpSatSolver) search(
	problem *preparedSchedulingProblem,
	variables []Variable,
	domains map[string][]DomainValue,
	index int,
	assignments map[string]PreviewAssignment,
	state *cpSearchState,
	telemetry *SolverTelemetry,
) {
	if state.nodesVisited >= s.maxNodes {
		if telemetry != nil {
			telemetry.HitMaxNodeLimit = true
		}
		return
	}
	state.nodesVisited++

	remaining := len(variables) - index
	if len(assignments)+remaining < state.bestAssigned {
		if telemetry != nil {
			telemetry.PrunedBranchCount++
		}
		return
	}

	if index >= len(variables) {
		if telemetry != nil {
			telemetry.LeafSolutionCount++
		}
		if s.maybeUpdateBest(assignments, state) && telemetry != nil {
			telemetry.BestSolutionUpdateCount++
		}
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
		if telemetry != nil {
			telemetry.CandidateEvaluatedCount++
		}
		if problem.hasConflict(variable, candidate, assignments) {
			if telemetry != nil {
				telemetry.CandidateRejectedConflictCount++
			}
			continue
		}

		assignments[variable.ID] = newPreviewAssignment(variable, candidate, "CP_SAT_OK")
		s.search(problem, variables, domains, index+1, assignments, state, telemetry)
		delete(assignments, variable.ID)
	}

	s.search(problem, variables, domains, index+1, assignments, state, telemetry)
}

func (s *cpSatSolver) maybeUpdateBest(assignments map[string]PreviewAssignment, state *cpSearchState) bool {
	assignmentCount := len(assignments)
	softScore := scoreAssignments(assignmentsToSlice(assignments), state.targetLessons)

	if assignmentCount > state.bestAssigned || (assignmentCount == state.bestAssigned && softScore > state.bestSoftScore) {
		state.bestAssigned = assignmentCount
		state.bestSoftScore = softScore
		state.bestAssignments = cloneAssignments(assignments)
		return true
	}
	return false
}
