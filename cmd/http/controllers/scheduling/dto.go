package scheduling

import "time"

type PreviewScheduleRequest struct {
	DateFrom   time.Time `json:"date_from" binding:"required"`
	DateTo     time.Time `json:"date_to" binding:"required"`
	ClassIDs   []string  `json:"class_ids"`
	TeacherIDs []string  `json:"teacher_ids"`
	RoomIDs    []string  `json:"room_ids"`
}

type CommitScheduleRequest struct {
	RunID string `json:"run_id" binding:"required"`
}
