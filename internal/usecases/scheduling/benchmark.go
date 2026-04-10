package scheduling

import (
	"context"
	"time"

	schedulingservice "doan/internal/services/scheduling"
)

type BenchmarkInput struct {
	DateFrom   time.Time
	DateTo     time.Time
	ClassIDs   []string
	TeacherIDs []string
	RoomIDs    []string
}

type BenchmarkUseCase interface {
	Execute(ctx context.Context, input BenchmarkInput) (*BenchmarkOutput, error)
}

type BenchmarkOutput struct {
	GeneratedAt time.Time               `json:"generated_at"`
	Filters     PreviewFilters          `json:"filters"`
	Mode        string                  `json:"mode"`
	Solvers     []BenchmarkSolverResult `json:"solvers"`
}

type BenchmarkSolverResult struct {
	Key                string                  `json:"key"`
	Label              string                  `json:"label"`
	Description        string                  `json:"description"`
	Readiness          string                  `json:"readiness"`
	ExecutionStatus    string                  `json:"execution_status"`
	Message            string                  `json:"message"`
	SelectedForMainAPI bool                    `json:"selected_for_main_api"`
	Metrics            *BenchmarkSolverMetrics `json:"metrics,omitempty"`
}

type BenchmarkSolverMetrics struct {
	FeasibilityRate    *float64        `json:"feasibility_rate,omitempty"`
	HardViolationCount *int            `json:"hard_violation_count,omitempty"`
	SoftScore          *int            `json:"soft_score,omitempty"`
	RuntimeMs          *int64          `json:"runtime_ms,omitempty"`
	Summary            *PreviewSummary `json:"summary,omitempty"`
}

type benchmarkUseCase struct {
	catalog schedulingservice.SolverCatalog
}

func NewBenchmarkUseCase(catalog schedulingservice.SolverCatalog) BenchmarkUseCase {
	return &benchmarkUseCase{catalog: catalog}
}

func (uc *benchmarkUseCase) Execute(_ context.Context, input BenchmarkInput) (*BenchmarkOutput, error) {
	descriptors := uc.catalog.BenchmarkSolvers()
	solvers := make([]BenchmarkSolverResult, 0, len(descriptors))
	for _, descriptor := range descriptors {
		solvers = append(solvers, BenchmarkSolverResult{
			Key:                descriptor.Key,
			Label:              descriptor.Label,
			Description:        descriptor.Description,
			Readiness:          descriptor.Readiness,
			ExecutionStatus:    "READY_FOR_BENCHMARK",
			Message:            "Solver da duoc implement o tang service. Metric benchmark se duoc dien khi buoc E chay benchmark thuc te.",
			SelectedForMainAPI: false,
		})
	}

	return &BenchmarkOutput{
		GeneratedAt: time.Now(),
		Filters: PreviewFilters{
			DateFrom:   input.DateFrom,
			DateTo:     input.DateTo,
			ClassIDs:   input.ClassIDs,
			TeacherIDs: input.TeacherIDs,
			RoomIDs:    input.RoomIDs,
		},
		Mode:    "ADMIN_BENCHMARK_CONTRACT",
		Solvers: solvers,
	}, nil
}
