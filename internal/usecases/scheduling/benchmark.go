package scheduling

import (
	"context"
	"fmt"
	"time"

	repositoryinterface "doan/internal/repositories/interface"
	schedulingservice "doan/internal/services/scheduling"
	"doan/pkg/logger"
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
	FeasibilityRate    *float64                           `json:"feasibility_rate,omitempty"`
	HardViolationCount *int                               `json:"hard_violation_count,omitempty"`
	SoftScore          *int                               `json:"soft_score,omitempty"`
	RuntimeMs          *int64                             `json:"runtime_ms,omitempty"`
	StartedAt          *time.Time                         `json:"started_at,omitempty"`
	FinishedAt         *time.Time                         `json:"finished_at,omitempty"`
	Summary            *PreviewSummary                    `json:"summary,omitempty"`
	Telemetry          *schedulingservice.SolverTelemetry `json:"telemetry,omitempty"`
}

type benchmarkUseCase struct {
	classRepo repositoryinterface.ClassRepository
	roomRepo  repositoryinterface.RoomRepository
	shiftRepo repositoryinterface.ShiftRepository
	catalog   schedulingservice.SolverCatalog
}

func NewBenchmarkUseCase(
	classRepo repositoryinterface.ClassRepository,
	roomRepo repositoryinterface.RoomRepository,
	shiftRepo repositoryinterface.ShiftRepository,
	catalog schedulingservice.SolverCatalog,
) BenchmarkUseCase {
	return &benchmarkUseCase{
		classRepo: classRepo,
		roomRepo:  roomRepo,
		shiftRepo: shiftRepo,
		catalog:   catalog,
	}
}

func (uc *benchmarkUseCase) Execute(ctx context.Context, input BenchmarkInput) (*BenchmarkOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	classes, err := loadSchedulingClasses(ctx, uc.classRepo, input.ClassIDs, input.TeacherIDs)
	if err != nil {
		ctxLogger.Errorf("Failed to load classes for scheduling benchmark: %v", err)
		return nil, err
	}

	rooms, err := loadSchedulingRooms(ctx, uc.roomRepo, input.RoomIDs)
	if err != nil {
		ctxLogger.Errorf("Failed to load rooms for scheduling benchmark: %v", err)
		return nil, err
	}

	shifts, err := loadActiveShifts(ctx, uc.shiftRepo)
	if err != nil {
		ctxLogger.Errorf("Failed to load shifts for scheduling benchmark: %v", err)
		return nil, err
	}

	solverInput := schedulingservice.SolverInput{
		DateFrom:   input.DateFrom,
		DateTo:     input.DateTo,
		ClassIDs:   input.ClassIDs,
		TeacherIDs: input.TeacherIDs,
		RoomIDs:    input.RoomIDs,
		Classes:    classes,
		Rooms:      rooms,
		Shifts:     shifts,
	}

	descriptors := uc.catalog.BenchmarkSolvers()
	solvers := make([]BenchmarkSolverResult, 0, len(descriptors))
	for _, descriptor := range descriptors {
		solver, ok := uc.catalog.GetSolver(descriptor.Key)
		if !ok {
			solvers = append(solvers, BenchmarkSolverResult{
				Key:                descriptor.Key,
				Label:              descriptor.Label,
				Description:        descriptor.Description,
				Readiness:          descriptor.Readiness,
				ExecutionStatus:    "UNAVAILABLE",
				Message:            "Solver chua duoc dang ky trong catalog benchmark.",
				SelectedForMainAPI: false,
			})
			continue
		}

		startedAt := time.Now()
		output, solveErr := solver.Solve(ctx, solverInput)
		finishedAt := time.Now()
		runtimeMs := finishedAt.Sub(startedAt).Milliseconds()
		if solveErr != nil {
			ctxLogger.Errorf("Failed to execute scheduling benchmark solver %s: %v", descriptor.Key, solveErr)
			solvers = append(solvers, BenchmarkSolverResult{
				Key:                descriptor.Key,
				Label:              descriptor.Label,
				Description:        descriptor.Description,
				Readiness:          descriptor.Readiness,
				ExecutionStatus:    "ERROR",
				Message:            solveErr.Error(),
				SelectedForMainAPI: false,
			})
			continue
		}

		feasibilityRate := calculateFeasibilityRate(output.Summary)
		hardViolationCount := output.Summary.ConflictCount
		softScore := output.Summary.SoftScore

		solvers = append(solvers, BenchmarkSolverResult{
			Key:                descriptor.Key,
			Label:              descriptor.Label,
			Description:        descriptor.Description,
			Readiness:          descriptor.Readiness,
			ExecutionStatus:    output.Status,
			Message:            buildBenchmarkMessage(output.Summary, runtimeMs),
			SelectedForMainAPI: false,
			Metrics: &BenchmarkSolverMetrics{
				FeasibilityRate:    &feasibilityRate,
				HardViolationCount: &hardViolationCount,
				SoftScore:          &softScore,
				RuntimeMs:          &runtimeMs,
				StartedAt:          &startedAt,
				FinishedAt:         &finishedAt,
				Summary: &PreviewSummary{
					RequestedClasses:   output.Summary.RequestedClasses,
					RequestedSessions:  output.Summary.RequestedSessions,
					ScheduledLessons:   output.Summary.ScheduledLessons,
					UnscheduledLessons: output.Summary.UnscheduledLessons,
					ConflictCount:      output.Summary.ConflictCount,
					SoftScore:          output.Summary.SoftScore,
				},
				Telemetry: output.Telemetry,
			},
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
		Mode:    "ADMIN_BENCHMARK_EXECUTION",
		Solvers: solvers,
	}, nil
}

func calculateFeasibilityRate(summary schedulingservice.SolverSummary) float64 {
	if summary.RequestedSessions == 0 {
		return 0
	}

	return float64(summary.ScheduledLessons) / float64(summary.RequestedSessions)
}

func buildBenchmarkMessage(summary schedulingservice.SolverSummary, runtimeMs int64) string {
	return fmt.Sprintf(
		"Da xep %d/%d buoi, %d trung, diem so mem %d, runtime %d ms.",
		summary.ScheduledLessons,
		summary.RequestedSessions,
		summary.ConflictCount,
		summary.SoftScore,
		runtimeMs,
	)
}
