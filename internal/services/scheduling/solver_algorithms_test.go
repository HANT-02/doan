package scheduling

import (
	"context"
	"testing"
	"time"

	"doan/internal/entities"
)

func TestSchedulingSolvers_BasicFeasibleInput(t *testing.T) {
	t.Parallel()

	input := SolverInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		Classes: []entities.Class{
			{
				ID:          "class-1",
				Code:        "L-001",
				Name:        "Lop 1",
				MaxStudents: 20,
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
				ClassSchedules: []entities.ClassSchedule{
					{
						ID:        "schedule-1",
						ClassID:   "class-1",
						DayOfWeek: "MONDAY",
						ShiftID:   "shift-1",
						Shift: entities.Shift{
							ID:          "shift-1",
							Code:        "S1",
							Name:        "Ca sang",
							StartTime:   "08:00",
							EndTime:     "10:00",
							SessionType: "MORNING",
							IsActive:    true,
						},
					},
				},
			},
			{
				ID:          "class-2",
				Code:        "L-002",
				Name:        "Lop 2",
				MaxStudents: 18,
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
				ClassSchedules: []entities.ClassSchedule{
					{
						ID:        "schedule-2",
						ClassID:   "class-2",
						DayOfWeek: "MONDAY",
						ShiftID:   "shift-2",
						Shift: entities.Shift{
							ID:          "shift-2",
							Code:        "S2",
							Name:        "Ca chieu",
							StartTime:   "13:30",
							EndTime:     "15:30",
							SessionType: "AFTERNOON",
							IsActive:    true,
						},
					},
				},
			},
		},
		Rooms: []entities.Room{
			{ID: "room-1", Name: "Phong 1", Capacity: 30},
		},
		Shifts: []entities.Shift{
			{ID: "shift-1", Code: "S1", Name: "Ca sang", StartTime: "08:00", EndTime: "10:00", SessionType: "MORNING", IsActive: true},
			{ID: "shift-2", Code: "S2", Name: "Ca chieu", StartTime: "13:30", EndTime: "15:30", SessionType: "AFTERNOON", IsActive: true},
		},
	}

	solvers := []SchedulingSolver{
		NewGraphColoringSolver(),
		NewCPSATSolver(),
		NewTabuSearchSolver(),
	}

	for _, solver := range solvers {
		solver := solver
		t.Run(solver.Key(), func(t *testing.T) {
			t.Parallel()

			output, err := solver.Solve(context.Background(), input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if output == nil {
				t.Fatalf("expected output")
			}

			if len(output.Assignments) != 2 {
				t.Fatalf("expected 2 assignments, got %d", len(output.Assignments))
			}

			if output.Summary.UnscheduledLessons != 0 {
				t.Fatalf("expected no unscheduled lessons, got %d", output.Summary.UnscheduledLessons)
			}

			if output.Telemetry == nil {
				t.Fatalf("expected telemetry for solver %s", solver.Key())
			}

			if output.Telemetry.SolverKey != solver.Key() {
				t.Fatalf("expected telemetry solver key %s, got %s", solver.Key(), output.Telemetry.SolverKey)
			}

			if output.Telemetry.RequestedSessionCount != 2 {
				t.Fatalf("expected telemetry requested sessions 2, got %d", output.Telemetry.RequestedSessionCount)
			}
		})
	}
}

func stringPtr(value string) *string {
	return &value
}
