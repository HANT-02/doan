package scheduling

import "time"

type TimeSlot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type Variable struct {
	ID              string `json:"id"`
	ClassID         string `json:"class_id"`
	ClassCode       string `json:"class_code"`
	ClassName       string `json:"class_name"`
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
	TeacherID     string    `json:"teacher_id"`
	TeacherLabel  string    `json:"teacher_label"`
	RoomID        string    `json:"room_id"`
	RoomName      string    `json:"room_name"`
	RoomCapacity  int       `json:"room_capacity"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	ConstraintFit string    `json:"constraint_fit"`
}

type PreviewConflict struct {
	VariableID string `json:"variable_id"`
	ClassID    string `json:"class_id"`
	ClassCode  string `json:"class_code"`
	ClassName  string `json:"class_name"`
	Type       string `json:"type"`
	Message    string `json:"message"`
}

type PreviewSummary struct {
	RequestedClasses   int `json:"requested_classes"`
	ScheduledLessons   int `json:"scheduled_lessons"`
	UnscheduledLessons int `json:"unscheduled_lessons"`
	ConflictCount      int `json:"conflict_count"`
	SoftScore          int `json:"soft_score"`
}

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
