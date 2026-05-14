package scheduling

import (
	"context"
	"sort"
)

type graphColoringSolver struct{}

func NewGraphColoringSolver() *graphColoringSolver {
	return &graphColoringSolver{}
}

func (s *graphColoringSolver) Key() string {
	return SolverKeyGraphColoring
}

func (s *graphColoringSolver) Label() string {
	return "Graph Coloring + Heuristic"
}

func (s *graphColoringSolver) Solve(_ context.Context, input SolverInput) (*SolverOutput, error) {
	problem := prepareSchedulingProblem(input)
	adjacency := buildConflictAdjacency(problem.variables)
	assignments := s.greedyColor(&problem, problem.variables, problem.domains, adjacency)

	return buildSolverOutput(input, problem.variables, assignments, problem.presetConflicts, problem.noDomainConflicts, nil, problem.targetLessons), nil
}

func (s *graphColoringSolver) greedyColor(
	problem *preparedSchedulingProblem,
	variables []Variable,
	domains map[string][]DomainValue,
	adjacency map[string]map[string]struct{},
) map[string]PreviewAssignment {
	assignments := make(map[string]PreviewAssignment)
	slotUsage := make(map[string]int)
	ordered := append([]Variable(nil), variables...)

	sort.Slice(ordered, func(i, j int) bool {
		leftDomain := len(domains[ordered[i].ID])
		rightDomain := len(domains[ordered[j].ID])
		if leftDomain == rightDomain {
			leftDegree := len(adjacency[ordered[i].ID])
			rightDegree := len(adjacency[ordered[j].ID])
			if leftDegree == rightDegree {
				return ordered[i].ClassCode < ordered[j].ClassCode
			}
			return leftDegree > rightDegree
		}
		return leftDomain < rightDomain
	})

	for _, variable := range ordered {
		grouped := groupDomainsBySlot(domains[variable.ID])
		slotKeys := make([]string, 0, len(grouped))
		for key := range grouped {
			slotKeys = append(slotKeys, key)
		}

		sort.Slice(slotKeys, func(i, j int) bool {
			leftScore := graphSlotPenalty(slotKeys[i], adjacency[variable.ID], assignments, slotUsage)
			rightScore := graphSlotPenalty(slotKeys[j], adjacency[variable.ID], assignments, slotUsage)
			if leftScore == rightScore {
				return slotKeys[i] < slotKeys[j]
			}
			return leftScore < rightScore
		})

		for _, slotKey := range slotKeys {
			for _, candidate := range grouped[slotKey] {
				if problem.hasConflict(variable, candidate, assignments) {
					continue
				}

				assignments[variable.ID] = newPreviewAssignment(variable, candidate, "GRAPH_COLORING_OK")
				slotUsage[slotKey]++
				goto assigned
			}
		}

	assigned:
	}

	return assignments
}

func graphSlotPenalty(
	slotKey string,
	neighbors map[string]struct{},
	assignments map[string]PreviewAssignment,
	slotUsage map[string]int,
) int {
	penalty := slotUsage[slotKey] * 10
	for neighborID := range neighbors {
		assignment, ok := assignments[neighborID]
		if !ok {
			continue
		}
		if slotKeyFromAssignment(assignment) == slotKey {
			penalty += 1000
		}
	}
	return penalty
}
