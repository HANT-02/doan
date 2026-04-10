package scheduling

type staticSolverCatalog struct {
	solvers map[string]SchedulingSolver
}

func NewSchedulingSolverCatalog(
	graphSolver *graphColoringSolver,
	cpSatSolver *cpSatSolver,
	tabuSolver *tabuSearchSolver,
) SolverCatalog {
	return &staticSolverCatalog{
		solvers: map[string]SchedulingSolver{
			SolverKeyGraphColoring: graphSolver,
			SolverKeyCPSAT:         cpSatSolver,
			SolverKeyTabuSearch:    tabuSolver,
		},
	}
}

func (c *staticSolverCatalog) BenchmarkSolvers() []SolverDescriptor {
	return []SolverDescriptor{
		{
			Key:         SolverKeyGraphColoring,
			Label:       "Graph Coloring + Heuristic",
			Description: "Baseline heuristic benchmark cho bài toán xếp lịch.",
			Readiness:   "READY",
		},
		{
			Key:         SolverKeyCPSAT,
			Label:       "CP-SAT",
			Description: "Constraint optimization solver dùng cho benchmark exact/exact-like.",
			Readiness:   "READY",
		},
		{
			Key:         SolverKeyTabuSearch,
			Label:       "Tabu Search",
			Description: "Local search/metaheuristic solver dùng cho benchmark timetabling.",
			Readiness:   "READY",
		},
	}
}

func (c *staticSolverCatalog) GetSolver(key string) (SchedulingSolver, bool) {
	solver, ok := c.solvers[key]
	return solver, ok
}

func NewDefaultSchedulingSolver() SchedulingSolver {
	return NewLegacyPreviewSolver()
}
