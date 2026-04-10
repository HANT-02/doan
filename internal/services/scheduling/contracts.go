package scheduling

import (
	"context"
	"time"

	"doan/internal/entities"
)

const (
	SolverKeyLegacyPreview = "legacy_preview"
	SolverKeyGraphColoring = "graph_coloring"
	SolverKeyCPSAT         = "cp_sat"
	SolverKeyTabuSearch    = "tabu_search"
)

type TimeSlot struct {
	Start           time.Time `json:"start"`
	End             time.Time `json:"end"`
	PreferredRoomID string    `json:"preferred_room_id,omitempty"`
	ShiftID         string    `json:"shift_id,omitempty"`
	ShiftCode       string    `json:"shift_code,omitempty"`
	ShiftName       string    `json:"shift_name,omitempty"`
	ShiftType       string    `json:"shift_type,omitempty"`
}

type Variable struct {
	ID              string `json:"id"`
	ClassID         string `json:"class_id"`
	ClassCode       string `json:"class_code"`
	ClassName       string `json:"class_name"`
	SessionIndex    int    `json:"session_index"`
	SessionTotal    int    `json:"session_total"`
	TeacherID       string `json:"teacher_id"`
	TeacherLabel    string `json:"teacher_label"`
	ExpectedCapcity int    `json:"expected_capacity"`
	DurationMinutes int    `json:"duration_minutes"`
	PreferredRoomID string `json:"preferred_room_id,omitempty"`
}

type DomainValue struct {
	RoomID       string   `json:"room_id"`
	RoomName     string   `json:"room_name"`
	RoomCapacity int      `json:"room_capacity"`
	TimeSlot     TimeSlot `json:"time_slot"`
}

type PreviewAssignment struct {
	VariableID    string    `json:"variable_id"`
	ClassID       string    `json:"class_id"`
	ClassCode     string    `json:"class_code"`
	ClassName     string    `json:"class_name"`
	SessionIndex  int       `json:"session_index"`
	SessionTotal  int       `json:"session_total"`
	TeacherID     string    `json:"teacher_id"`
	TeacherLabel  string    `json:"teacher_label"`
	RoomID        string    `json:"room_id"`
	RoomName      string    `json:"room_name"`
	RoomCapacity  int       `json:"room_capacity"`
	ShiftID       string    `json:"shift_id,omitempty"`
	ShiftCode     string    `json:"shift_code,omitempty"`
	ShiftName     string    `json:"shift_name,omitempty"`
	ShiftType     string    `json:"shift_type,omitempty"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	ConstraintFit string    `json:"constraint_fit"`
}

type PreviewConflict struct {
	VariableID   string `json:"variable_id"`
	ClassID      string `json:"class_id"`
	ClassCode    string `json:"class_code"`
	ClassName    string `json:"class_name"`
	SessionIndex int    `json:"session_index,omitempty"`
	SessionTotal int    `json:"session_total,omitempty"`
	Type         string `json:"type"`
	Message      string `json:"message"`
}

type SolverSummary struct {
	RequestedClasses   int `json:"requested_classes"`
	RequestedSessions  int `json:"requested_sessions"`
	ScheduledLessons   int `json:"scheduled_lessons"`
	UnscheduledLessons int `json:"unscheduled_lessons"`
	ConflictCount      int `json:"conflict_count"`
	SoftScore          int `json:"soft_score"`
}

type SolverInput struct {
	DateFrom   time.Time
	DateTo     time.Time
	ClassIDs   []string
	TeacherIDs []string
	RoomIDs    []string
	Classes    []entities.Class
	Rooms      []entities.Room
	Shifts     []entities.Shift
}

type SolverOutput struct {
	Status      string
	Assignments []PreviewAssignment
	Conflicts   []PreviewConflict
	Summary     SolverSummary
}

type SchedulingSolver interface {
	Key() string
	Label() string
	Solve(ctx context.Context, input SolverInput) (*SolverOutput, error)
}

type SolverDescriptor struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Readiness   string `json:"readiness"`
}

type SolverCatalog interface {
	BenchmarkSolvers() []SolverDescriptor
	GetSolver(key string) (SchedulingSolver, bool)
}
