package lesson

import "time"

type UpsertLessonAttendanceRequest struct {
	Records []AttendanceRecordPayload `json:"records" binding:"required"`
}

type AttendanceRecordPayload struct {
	StudentID string `json:"student_id" binding:"required"`
	Status    int    `json:"status" binding:"required"`
	Note      string `json:"note"`
}

type UpsertLessonSummaryRequest struct {
	Topic            string     `json:"topic"`
	LessonContent    string     `json:"lesson_content"`
	ClassFeedback    string     `json:"class_feedback"`
	Homework         string     `json:"homework"`
	HomeworkDeadline *time.Time `json:"homework_deadline"`
	TeacherNotes     string     `json:"teacher_notes"`
}

type UpsertLessonAcademicRecordsRequest struct {
	Records []AcademicRecordPayload `json:"records" binding:"required"`
}

type AcademicRecordPayload struct {
	StudentID          string  `json:"student_id" binding:"required"`
	HomeworkCompleted  bool    `json:"homework_completed"`
	HomeworkScore      float64 `json:"homework_score"`
	AttitudeRating     int     `json:"attitude_rating"`
	ParticipationScore float64 `json:"participation_score"`
	PersonalComment    string  `json:"personal_comment"`
}
