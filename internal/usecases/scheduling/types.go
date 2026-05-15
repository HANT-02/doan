package scheduling

import (
	"time"

	"doan/internal/entities"
	schedulingservice "doan/internal/services/scheduling"
)

type TimeSlot = schedulingservice.TimeSlot
type Variable = schedulingservice.Variable
type DomainValue = schedulingservice.DomainValue
type PreviewAssignment = schedulingservice.PreviewAssignment
type PreviewCandidateOption = schedulingservice.PreviewCandidateOption
type ExistingLesson = schedulingservice.ExistingLesson
type PreviewConflict = schedulingservice.PreviewConflict
type PreviewSummary = schedulingservice.SolverSummary

type PreviewResult struct {
	RunID             string                              `json:"run_id"`
	Mode              string                              `json:"mode"`
	Status            string                              `json:"status"`
	GeneratedAt       time.Time                           `json:"generated_at"`
	EffectiveDateFrom time.Time                           `json:"effective_date_from"`
	Filters           PreviewFilters                      `json:"filters"`
	Summary           PreviewSummary                      `json:"summary"`
	Assignments       []PreviewAssignment                 `json:"assignments"`
	Conflicts         []PreviewConflict                   `json:"conflicts"`
	ExistingLessons   []ExistingLesson                    `json:"existing_lessons,omitempty"`
	CandidateOptions  map[string][]PreviewCandidateOption `json:"candidate_options,omitempty"`

	Variables         []Variable                     `json:"-"`
	PresetConflicts   []PreviewConflict              `json:"-"`
	NoDomainConflicts map[string]PreviewConflict     `json:"-"`
	DomainOptions     map[string][]DomainValue       `json:"-"`
	ClassStudentIDs   map[string]map[string]struct{} `json:"-"`
	RoomsByID         map[string]entities.Room       `json:"-"`
	TravelMap         map[string]int                 `json:"-"`
}

type PreviewFilters struct {
	DateFrom          time.Time `json:"date_from"`
	DateTo            time.Time `json:"date_to"`
	EffectiveDateFrom time.Time `json:"effective_date_from"`
	Mode              string    `json:"mode,omitempty"`
	ClassIDs          []string  `json:"class_ids,omitempty"`
	TeacherIDs        []string  `json:"teacher_ids,omitempty"`
	RoomIDs           []string  `json:"room_ids,omitempty"`
}

type ManualAssignmentOverride struct {
	VariableID string
	OptionKey  string
}
