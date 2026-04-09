package scheduling

type staticSolverCatalog struct{}

func NewSchedulingSolverCatalog() SolverCatalog {
	return &staticSolverCatalog{}
}

func (c *staticSolverCatalog) BenchmarkSolvers() []SolverDescriptor {
	return []SolverDescriptor{
		{
			Key:         SolverKeyGraphColoring,
			Label:       "Graph Coloring + Heuristic",
			Description: "Baseline heuristic benchmark cho bài toán xếp lịch.",
			Readiness:   "PLANNED",
		},
		{
			Key:         SolverKeyCPSAT,
			Label:       "CP-SAT",
			Description: "Constraint optimization solver dùng cho benchmark exact/exact-like.",
			Readiness:   "PLANNED",
		},
		{
			Key:         SolverKeyTabuSearch,
			Label:       "Tabu Search",
			Description: "Local search/metaheuristic solver dùng cho benchmark timetabling.",
			Readiness:   "PLANNED",
		},
	}
}

func NewDefaultSchedulingSolver() SchedulingSolver {
	return NewLegacyPreviewSolver()
}
