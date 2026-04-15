package leave

type CreateLeaveRequestRequest struct {
	LeaveType    string   `json:"leave_type" binding:"required"`
	ApplyDate    string   `json:"apply_date" binding:"required"`
	LateMinutes  int      `json:"late_minutes"`
	EarlyMinutes int      `json:"early_minutes"`
	Reason       string   `json:"reason" binding:"required"`
	Documents    []string `json:"documents"`
	ClassID      *string  `json:"class_id"`
	LessonID     *string  `json:"lesson_id"`
	Subject      string   `json:"subject"`
}

type RejectLeaveRequestRequest struct {
	RejectionReason string `json:"rejection_reason" binding:"required"`
}
