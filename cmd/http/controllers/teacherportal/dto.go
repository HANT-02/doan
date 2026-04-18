package teacherportal

import "time"

type TeacherLessonShiftResponse struct {
	ID              string `json:"id"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	DurationMinutes int    `json:"duration_minutes"`
	SessionType     string `json:"session_type"`
}

type TeacherLessonResponse struct {
	ID        string                      `json:"id"`
	ClassID   string                      `json:"class_id"`
	ClassName string                      `json:"class_name"`
	ClassCode string                      `json:"class_code"`
	RoomID    *string                     `json:"room_id,omitempty"`
	RoomName  *string                     `json:"room_name,omitempty"`
	DateStart time.Time                   `json:"date_start"`
	DateEnd   time.Time                   `json:"date_end"`
	Notes     string                      `json:"notes"`
	Shift     *TeacherLessonShiftResponse `json:"shift,omitempty"`
}

type GetTeacherLessonsResponse struct {
	TeacherID string                  `json:"teacher_id"`
	Lessons   []TeacherLessonResponse `json:"lessons"`
}

type TeacherAttendanceStudentResponse struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"full_name"`
}

type TeacherAttendanceRecordResponse struct {
	AttendanceID *string                          `json:"attendance_id,omitempty"`
	Student      TeacherAttendanceStudentResponse `json:"student"`
	Status       int                              `json:"status"`
	Note         string                           `json:"note"`
	MarkedAt     *time.Time                       `json:"marked_at,omitempty"`
}

type GetTeacherLessonAttendanceResponse struct {
	Lesson  TeacherLessonResponse             `json:"lesson"`
	Records []TeacherAttendanceRecordResponse `json:"records"`
}

type SubmitTeacherLessonAttendanceRequest struct {
	Records []SubmitTeacherLessonAttendanceRecordRequest `json:"records" binding:"required"`
}

type SubmitTeacherLessonAttendanceRecordRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Status    int    `json:"status" binding:"required"`
	Note      string `json:"note"`
}

type UpdateTeacherLessonAttendanceRequest struct {
	Status int    `json:"status" binding:"required"`
	Note   string `json:"note"`
}

type TeacherAttendanceSaveResponse struct {
	SavedCount int `json:"saved_count"`
}

type TeacherAttendanceStudentSummaryResponse struct {
	Student        TeacherAttendanceStudentResponse `json:"student"`
	TotalLessons   int                              `json:"total_lessons"`
	MarkedCount    int                              `json:"marked_count"`
	PresentCount   int                              `json:"present_count"`
	AbsentCount    int                              `json:"absent_count"`
	LateCount      int                              `json:"late_count"`
	ExcusedCount   int                              `json:"excused_count"`
	UnmarkedCount  int                              `json:"unmarked_count"`
	AttendanceRate float64                          `json:"attendance_rate"`
}

type GetTeacherAttendanceSummaryResponse struct {
	TeacherID    string                                    `json:"teacher_id"`
	ClassID      string                                    `json:"class_id"`
	TotalLessons int                                       `json:"total_lessons"`
	Students     []TeacherAttendanceStudentSummaryResponse `json:"students"`
}

type TeacherLessonSummaryResponse struct {
	ID               string     `json:"id"`
	LessonID         string     `json:"lesson_id"`
	Topic            string     `json:"topic"`
	LessonContent    string     `json:"lesson_content"`
	ClassFeedback    string     `json:"class_feedback"`
	Homework         string     `json:"homework"`
	HomeworkDeadline *time.Time `json:"homework_deadline,omitempty"`
	TeacherNotes     string     `json:"teacher_notes"`
	CreatedByID      *string    `json:"created_by_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type GetTeacherLessonSummaryResponse struct {
	Lesson  TeacherLessonResponse         `json:"lesson"`
	Summary *TeacherLessonSummaryResponse `json:"summary,omitempty"`
}

type UpsertTeacherLessonSummaryRequest struct {
	Topic            string     `json:"topic"`
	LessonContent    string     `json:"lesson_content"`
	ClassFeedback    string     `json:"class_feedback"`
	Homework         string     `json:"homework"`
	HomeworkDeadline *time.Time `json:"homework_deadline"`
	TeacherNotes     string     `json:"teacher_notes"`
}

type TeacherAcademicRecordResponse struct {
	RecordID           *string                          `json:"record_id,omitempty"`
	LessonSummaryID    *string                          `json:"lesson_summary_id,omitempty"`
	Student            TeacherAttendanceStudentResponse `json:"student"`
	HomeworkCompleted  bool                             `json:"homework_completed"`
	HomeworkScore      float64                          `json:"homework_score"`
	AttitudeRating     int                              `json:"attitude_rating"`
	ParticipationScore float64                          `json:"participation_score"`
	PersonalComment    string                           `json:"personal_comment"`
	TotalScore         float64                          `json:"total_score"`
	IsCompleted        bool                             `json:"is_completed"`
	CreatedAt          *time.Time                       `json:"created_at,omitempty"`
	UpdatedAt          *time.Time                       `json:"updated_at,omitempty"`
}

type GetTeacherLessonAcademicRecordsResponse struct {
	Lesson  TeacherLessonResponse           `json:"lesson"`
	Records []TeacherAcademicRecordResponse `json:"records"`
}

type UpsertTeacherLessonAcademicRecordRequest struct {
	HomeworkCompleted  bool    `json:"homework_completed"`
	HomeworkScore      float64 `json:"homework_score"`
	AttitudeRating     int     `json:"attitude_rating"`
	ParticipationScore float64 `json:"participation_score"`
	PersonalComment    string  `json:"personal_comment"`
}

type TeacherAcademicRecordSaveResponse struct {
	SavedCount int `json:"saved_count"`
}

type TeacherAcademicRecordFinalizeResponse struct {
	FinalizedCount int `json:"finalized_count"`
}

type GetTeacherStudentAcademicRecordsResponse struct {
	TeacherID string                           `json:"teacher_id"`
	ClassID   string                           `json:"class_id"`
	StudentID string                           `json:"student_id"`
	Student   TeacherAttendanceStudentResponse `json:"student"`
	Records   []TeacherAcademicRecordResponse  `json:"records"`
}

type ListTeacherLeaveRequestsResponse struct {
	Requests []TeacherLeaveRequestResponse `json:"requests"`
}

type TeacherLeaveRequestStudentResponse struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"full_name"`
}

type TeacherLeaveRequestClassResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type TeacherLeaveRequestLessonResponse struct {
	ID        string    `json:"id"`
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
}

type TeacherLeaveRequestResponse struct {
	ID              string                             `json:"id"`
	Student         TeacherLeaveRequestStudentResponse `json:"student"`
	LeaveType       string                             `json:"leave_type"`
	ApplyDate       time.Time                          `json:"apply_date"`
	LateMinutes     int                                `json:"late_minutes"`
	EarlyMinutes    int                                `json:"early_minutes"`
	Reason          string                             `json:"reason"`
	Documents       []string                           `json:"documents"`
	Class           *TeacherLeaveRequestClassResponse  `json:"class,omitempty"`
	Lesson          *TeacherLeaveRequestLessonResponse `json:"lesson,omitempty"`
	Subject         string                             `json:"subject"`
	Status          string                             `json:"status"`
	ApprovedByID    *string                            `json:"approved_by_id,omitempty"`
	ApprovedAt      *time.Time                         `json:"approved_at,omitempty"`
	RejectionReason string                             `json:"rejection_reason"`
	CreatedAt       time.Time                          `json:"created_at"`
	UpdatedAt       time.Time                          `json:"updated_at"`
}

type RejectTeacherLeaveRequestRequest struct {
	RejectionReason string `json:"rejection_reason" binding:"required"`
}

type TeacherLeaveRequestStatusResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}
