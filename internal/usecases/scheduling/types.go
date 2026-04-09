package scheduling

import (
	"time"

	schedulingservice "doan/internal/services/scheduling"
)

type TimeSlot = schedulingservice.TimeSlot
type Variable = schedulingservice.Variable
type DomainValue = schedulingservice.DomainValue
type PreviewAssignment = schedulingservice.PreviewAssignment
type PreviewConflict = schedulingservice.PreviewConflict
type PreviewSummary = schedulingservice.SolverSummary

type PreviewResult struct {
	RunID       string              `json:"run_id"`
	Status      string              `json:"status"`
	GeneratedAt time.Time           `json:"generated_at"`
	Filters     PreviewFilters      `json:"filters"`
	Summary     PreviewSummary      `json:"summary"`
	Assignments []PreviewAssignment `json:"assignments"`
	Conflicts   []PreviewConflict   `json:"conflicts"`
}

type PreviewFilters struct {
	DateFrom   time.Time `json:"date_from"`
	DateTo     time.Time `json:"date_to"`
	ClassIDs   []string  `json:"class_ids,omitempty"`
	TeacherIDs []string  `json:"teacher_ids,omitempty"`
	RoomIDs    []string  `json:"room_ids,omitempty"`
}
