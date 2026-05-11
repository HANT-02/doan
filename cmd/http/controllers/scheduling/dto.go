package scheduling

type PreviewScheduleRequest struct {
	DateFrom   string   `json:"date_from" binding:"required"`
	DateTo     string   `json:"date_to" binding:"required"`
	ClassIDs   []string `json:"class_ids"`
	TeacherIDs []string `json:"teacher_ids"`
	RoomIDs    []string `json:"room_ids"`
}

type CommitScheduleRequest struct {
	RunID             string                   `json:"run_id" binding:"required"`
	ManualAssignments []CommitAssignmentChoice `json:"manual_assignments"`
}

type CommitAssignmentChoice struct {
	VariableID string `json:"variable_id" binding:"required"`
	OptionKey  string `json:"option_key" binding:"required"`
}
