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
			Label:       "Tô màu đồ thị",
			Description: "Thuật toán tham lam làm mốc so sánh ban đầu cho bài toán xếp lịch.",
			Readiness:   "READY",
		},
		{
			Key:         SolverKeyCPSAT,
			Label:       "CP-SAT",
			Description: "Bộ giải tối ưu hóa ràng buộc để so sánh chất lượng nghiệm và độ ổn định.",
			Readiness:   "READY",
		},
		{
			Key:         SolverKeyTabuSearch,
			Label:       "Tìm kiếm Tabu",
			Description: "Thuật toán tìm kiếm cục bộ có danh sách cấm để tối ưu lịch học.",
			Readiness:   "READY",
		},
	}
}

func (c *staticSolverCatalog) GetSolver(key string) (SchedulingSolver, bool) {
	solver, ok := c.solvers[key]
	return solver, ok
}

func NewDefaultSchedulingSolver(cpSatSolver *cpSatSolver) SchedulingSolver {
	return cpSatSolver
}
