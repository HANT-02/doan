package scheduling

import (
	"context"
	"testing"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	schedulingservice "doan/internal/services/scheduling"
)

func TestBenchmarkUseCase_ExecuteRunsAllSolvers(t *testing.T) {
	t.Parallel()

	uc := NewBenchmarkUseCase(
		benchmarkClassRepoStub{
			data: []entities.Class{
				{
					ID:          "class-1",
					Code:        "L-001",
					Name:        "Lop 1",
					MaxStudents: 20,
					Status:      "OPEN",
					TeacherID:   stringPtr("teacher-1"),
					Teacher: entities.Teacher{
						ID:       "teacher-1",
						FullName: "GV 1",
						Code:     "GV-1",
					},
					CourseID: stringPtr("course-1"),
					Course: entities.Course{
						ID:                     "course-1",
						SessionCount:           1,
						SessionDurationMinutes: 120,
					},
				},
				{
					ID:          "class-2",
					Code:        "L-002",
					Name:        "Lop 2",
					MaxStudents: 18,
					Status:      "OPEN",
					TeacherID:   stringPtr("teacher-2"),
					Teacher: entities.Teacher{
						ID:       "teacher-2",
						FullName: "GV 2",
						Code:     "GV-2",
					},
					CourseID: stringPtr("course-2"),
					Course: entities.Course{
						ID:                     "course-2",
						SessionCount:           1,
						SessionDurationMinutes: 120,
					},
				},
			},
		},
		benchmarkRoomRepoStub{
			data: []entities.Room{
				{ID: "room-1", Name: "Phong 1", Capacity: 30},
			},
		},
		benchmarkShiftRepoStub{
			data: []entities.Shift{
				{ID: "shift-1", Code: "S1", Name: "Ca sang", StartTime: "08:00", EndTime: "10:00", SessionType: "MORNING", IsActive: true},
				{ID: "shift-2", Code: "S2", Name: "Ca chieu", StartTime: "13:30", EndTime: "15:30", SessionType: "AFTERNOON", IsActive: true},
			},
		},
		schedulingservice.NewSchedulingSolverCatalog(
			schedulingservice.NewGraphColoringSolver(),
			schedulingservice.NewCPSATSolver(),
			schedulingservice.NewTabuSearchSolver(),
		),
	)

	output, err := uc.Execute(context.Background(), BenchmarkInput{
		DateFrom: time.Date(2026, 4, 14, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 14, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output == nil {
		t.Fatalf("expected output")
	}

	if output.Mode != "ADMIN_BENCHMARK_EXECUTION" {
		t.Fatalf("expected execution mode, got %s", output.Mode)
	}

	if len(output.Solvers) != 3 {
		t.Fatalf("expected 3 solver results, got %d", len(output.Solvers))
	}

	for _, solver := range output.Solvers {
		if solver.ExecutionStatus != "COMPLETED" {
			t.Fatalf("expected solver %s to complete, got %s", solver.Key, solver.ExecutionStatus)
		}

		if solver.Metrics == nil {
			t.Fatalf("expected metrics for solver %s", solver.Key)
		}

		if solver.Metrics.FeasibilityRate == nil || *solver.Metrics.FeasibilityRate != 1 {
			t.Fatalf("expected feasibility rate 1 for solver %s", solver.Key)
		}

		if solver.Metrics.HardViolationCount == nil || *solver.Metrics.HardViolationCount != 0 {
			t.Fatalf("expected no hard violations for solver %s", solver.Key)
		}

		if solver.Metrics.Summary == nil || solver.Metrics.Summary.ScheduledLessons != 2 {
			t.Fatalf("expected scheduled lessons summary for solver %s", solver.Key)
		}
	}
}

var _ repositoryinterface.ClassRepository = benchmarkClassRepoStub{}
var _ repositoryinterface.RoomRepository = benchmarkRoomRepoStub{}
var _ repositoryinterface.ShiftRepository = benchmarkShiftRepoStub{}

type benchmarkClassRepoStub struct {
	data []entities.Class
}

func (s benchmarkClassRepoStub) GetTable() string { return "classes" }

func (s benchmarkClassRepoStub) GetByCondition(_ context.Context, _ *repositories.CommonCondition) (*repositories.Pagination[entities.Class], error) {
	return benchmarkPaginationFromSlice(s.data), nil
}

func (s benchmarkClassRepoStub) GetTotal(_ context.Context, _ *repositories.CommonCondition) (uint64, error) {
	return uint64(len(s.data)), nil
}

func (s benchmarkClassRepoStub) Create(_ context.Context, entity *entities.Class) (*entities.Class, error) {
	return entity, nil
}

func (s benchmarkClassRepoStub) Update(_ context.Context, _ interface{}, _ map[string]interface{}) error {
	return nil
}

func (s benchmarkClassRepoStub) UpdateWithIDs(_ context.Context, _ []string, _ map[string]interface{}) error {
	return nil
}

func (s benchmarkClassRepoStub) SoftDelete(_ context.Context, _ interface{}) error {
	return nil
}

func (s benchmarkClassRepoStub) HardDelete(_ context.Context, _ interface{}) error {
	return nil
}

func (s benchmarkClassRepoStub) GetByID(_ context.Context, _ interface{}) (*entities.Class, error) {
	return nil, nil
}

type benchmarkRoomRepoStub struct {
	data []entities.Room
}

func (s benchmarkRoomRepoStub) GetTable() string { return "rooms" }

func (s benchmarkRoomRepoStub) GetByCondition(_ context.Context, _ *repositories.CommonCondition) (*repositories.Pagination[entities.Room], error) {
	return benchmarkPaginationFromSlice(s.data), nil
}

func (s benchmarkRoomRepoStub) GetTotal(_ context.Context, _ *repositories.CommonCondition) (uint64, error) {
	return uint64(len(s.data)), nil
}

func (s benchmarkRoomRepoStub) Create(_ context.Context, entity *entities.Room) (*entities.Room, error) {
	return entity, nil
}

func (s benchmarkRoomRepoStub) Update(_ context.Context, _ interface{}, _ map[string]interface{}) error {
	return nil
}

func (s benchmarkRoomRepoStub) UpdateWithIDs(_ context.Context, _ []string, _ map[string]interface{}) error {
	return nil
}

func (s benchmarkRoomRepoStub) SoftDelete(_ context.Context, _ interface{}) error {
	return nil
}

func (s benchmarkRoomRepoStub) HardDelete(_ context.Context, _ interface{}) error {
	return nil
}

func (s benchmarkRoomRepoStub) GetByID(_ context.Context, _ interface{}) (*entities.Room, error) {
	return nil, nil
}

type benchmarkShiftRepoStub struct {
	data []entities.Shift
}

func (s benchmarkShiftRepoStub) GetTable() string { return "shifts" }

func (s benchmarkShiftRepoStub) GetByCondition(_ context.Context, _ *repositories.CommonCondition) (*repositories.Pagination[entities.Shift], error) {
	return benchmarkPaginationFromSlice(s.data), nil
}

func (s benchmarkShiftRepoStub) GetTotal(_ context.Context, _ *repositories.CommonCondition) (uint64, error) {
	return uint64(len(s.data)), nil
}

func (s benchmarkShiftRepoStub) Create(_ context.Context, entity *entities.Shift) (*entities.Shift, error) {
	return entity, nil
}

func (s benchmarkShiftRepoStub) Update(_ context.Context, _ interface{}, _ map[string]interface{}) error {
	return nil
}

func (s benchmarkShiftRepoStub) UpdateWithIDs(_ context.Context, _ []string, _ map[string]interface{}) error {
	return nil
}

func (s benchmarkShiftRepoStub) SoftDelete(_ context.Context, _ interface{}) error {
	return nil
}

func (s benchmarkShiftRepoStub) HardDelete(_ context.Context, _ interface{}) error {
	return nil
}

func (s benchmarkShiftRepoStub) GetByID(_ context.Context, _ interface{}) (*entities.Shift, error) {
	return nil, nil
}

func benchmarkPaginationFromSlice[T any](items []T) *repositories.Pagination[T] {
	data := make([]*T, 0, len(items))
	for i := range items {
		item := items[i]
		data = append(data, &item)
	}

	return &repositories.Pagination[T]{
		Data: data,
	}
}

func stringPtr(value string) *string {
	return &value
}
