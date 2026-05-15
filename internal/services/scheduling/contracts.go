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

	PreviewModeColdStart               = "cold_start"
	PreviewModeReplanDraft             = "replan_draft"
	PreviewModeReplanWithPublishedLock = "replan_with_published_lock"
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
	ReplaceLessonID string `json:"replace_lesson_id,omitempty"`
}

type DomainValue struct {
	RoomID       string   `json:"room_id"`
	RoomName     string   `json:"room_name"`
	RoomCapacity int      `json:"room_capacity"`
	TimeSlot     TimeSlot `json:"time_slot"`
}

type PreviewAssignment struct {
	VariableID           string    `json:"variable_id"`
	ClassID              string    `json:"class_id"`
	ClassCode            string    `json:"class_code"`
	ClassName            string    `json:"class_name"`
	SessionIndex         int       `json:"session_index"`
	SessionTotal         int       `json:"session_total"`
	TeacherID            string    `json:"teacher_id"`
	TeacherLabel         string    `json:"teacher_label"`
	RoomID               string    `json:"room_id"`
	RoomName             string    `json:"room_name"`
	RoomCapacity         int       `json:"room_capacity"`
	ExpectedStudentCount int       `json:"expected_student_count"`
	ReplaceLessonID      string    `json:"replace_lesson_id,omitempty"`
	ShiftID              string    `json:"shift_id,omitempty"`
	ShiftCode            string    `json:"shift_code,omitempty"`
	ShiftName            string    `json:"shift_name,omitempty"`
	ShiftType            string    `json:"shift_type,omitempty"`
	StartTime            time.Time `json:"start_time"`
	EndTime              time.Time `json:"end_time"`
	ConstraintFit        string    `json:"constraint_fit"`
}

type PreviewCandidateOption struct {
	Key          string    `json:"key"`
	RoomID       string    `json:"room_id"`
	RoomName     string    `json:"room_name"`
	RoomCapacity int       `json:"room_capacity"`
	ShiftID      string    `json:"shift_id,omitempty"`
	ShiftCode    string    `json:"shift_code,omitempty"`
	ShiftName    string    `json:"shift_name,omitempty"`
	ShiftType    string    `json:"shift_type,omitempty"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}

type ExistingLesson struct {
	LessonID     string    `json:"lesson_id"`
	ClassID      string    `json:"class_id"`
	ClassCode    string    `json:"class_code"`
	ClassName    string    `json:"class_name"`
	Status       string    `json:"status,omitempty"`
	TeacherID    string    `json:"teacher_id,omitempty"`
	TeacherLabel string    `json:"teacher_label,omitempty"`
	RoomID       string    `json:"room_id,omitempty"`
	RoomName     string    `json:"room_name,omitempty"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Notes        string    `json:"notes,omitempty"`
	StudentIDs   []string  `json:"student_ids,omitempty"`
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
	RequestedClasses           int     `json:"requested_classes"`
	RequestedSessions          int     `json:"requested_sessions"`
	ScheduledLessons           int     `json:"scheduled_lessons"`
	UnscheduledLessons         int     `json:"unscheduled_lessons"`
	ConflictCount              int     `json:"conflict_count"`
	SoftScore                  int     `json:"soft_score"`
	ScheduleChangeCount        int     `json:"schedule_change_count"`
	TeacherChangeCount         int     `json:"teacher_change_count"`
	RoomChangeCount            int     `json:"room_change_count"`
	AverageCapacityUtilization float64 `json:"average_capacity_utilization"`
}

type SolverInput struct {
	DateFrom         time.Time
	DateTo           time.Time
	ClassIDs         []string
	TeacherIDs       []string
	RoomIDs          []string
	Classes          []entities.Class
	ClassWindows     map[string]ClassSchedulingWindow
	Rooms            []entities.Room
	Shifts           []entities.Shift
	RoomsByID        map[string]entities.Room
	TravelMap        map[string]int
	TargetLessons    []entities.Lesson
	BenchmarkOptions *BenchmarkOptions
}

type SolverOutput struct {
	Status      string
	Assignments []PreviewAssignment
	Conflicts   []PreviewConflict
	Summary     SolverSummary
	Telemetry   *SolverTelemetry
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
