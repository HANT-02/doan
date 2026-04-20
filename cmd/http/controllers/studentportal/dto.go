package studentportal

import "time"

type StudentTimetableShiftResponse struct {
	ID              string `json:"id"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	DurationMinutes int    `json:"duration_minutes"`
	SessionType     string `json:"session_type"`
}

type StudentTimetableTeacherResponse struct {
	ID       *string `json:"id,omitempty"`
	Code     *string `json:"code,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

type StudentTimetableLessonResponse struct {
	ID        string                          `json:"id"`
	ClassID   string                          `json:"class_id"`
	ClassName string                          `json:"class_name"`
	ClassCode string                          `json:"class_code"`
	Teacher   StudentTimetableTeacherResponse `json:"teacher"`
	RoomID    *string                         `json:"room_id,omitempty"`
	RoomName  *string                         `json:"room_name,omitempty"`
	DateStart time.Time                       `json:"date_start"`
	DateEnd   time.Time                       `json:"date_end"`
	Notes     string                          `json:"notes"`
	Shift     *StudentTimetableShiftResponse  `json:"shift,omitempty"`
}

type GetStudentTimetableResponse struct {
	StudentID string                           `json:"student_id"`
	Lessons   []StudentTimetableLessonResponse `json:"lessons"`
}

type StudentAttendanceSummaryResponse struct {
	TotalLessons   int     `json:"total_lessons"`
	MarkedCount    int     `json:"marked_count"`
	PresentCount   int     `json:"present_count"`
	AbsentCount    int     `json:"absent_count"`
	LateCount      int     `json:"late_count"`
	ExcusedCount   int     `json:"excused_count"`
	UnmarkedCount  int     `json:"unmarked_count"`
	AttendanceRate float64 `json:"attendance_rate"`
	AbsentRate     float64 `json:"absent_rate"`
	Warning        bool    `json:"warning"`
	WarningMessage string  `json:"warning_message,omitempty"`
}

type StudentAttendanceRecordResponse struct {
	Lesson   StudentTimetableLessonResponse `json:"lesson"`
	Status   *int                           `json:"status,omitempty"`
	Note     string                         `json:"note"`
	MarkedAt *time.Time                     `json:"marked_at,omitempty"`
}

type GetMyAttendanceResponse struct {
	StudentID string                            `json:"student_id"`
	ClassID   string                            `json:"class_id,omitempty"`
	Summary   StudentAttendanceSummaryResponse  `json:"summary"`
	Records   []StudentAttendanceRecordResponse `json:"records"`
}

type StudentAcademicRecordLessonSummaryResponse struct {
	ID       string `json:"id"`
	Topic    string `json:"topic"`
	Homework string `json:"homework"`
}

type StudentAcademicRecordLessonResponse struct {
	ID        string                                     `json:"id"`
	ClassID   string                                     `json:"class_id"`
	ClassName string                                     `json:"class_name"`
	ClassCode string                                     `json:"class_code"`
	Teacher   StudentTimetableTeacherResponse            `json:"teacher"`
	RoomID    *string                                    `json:"room_id,omitempty"`
	RoomName  *string                                    `json:"room_name,omitempty"`
	DateStart time.Time                                  `json:"date_start"`
	DateEnd   time.Time                                  `json:"date_end"`
	Notes     string                                     `json:"notes"`
	Shift     *StudentTimetableShiftResponse             `json:"shift,omitempty"`
	Summary   StudentAcademicRecordLessonSummaryResponse `json:"summary"`
}

type StudentAcademicRecordItemResponse struct {
	RecordID           string                              `json:"record_id"`
	LessonSummaryID    string                              `json:"lesson_summary_id"`
	Lesson             StudentAcademicRecordLessonResponse `json:"lesson"`
	HomeworkCompleted  bool                                `json:"homework_completed"`
	HomeworkScore      float64                             `json:"homework_score"`
	AttitudeRating     int                                 `json:"attitude_rating"`
	ParticipationScore float64                             `json:"participation_score"`
	PersonalComment    string                              `json:"personal_comment"`
	TotalScore         float64                             `json:"total_score"`
	IsCompleted        bool                                `json:"is_completed"`
	CreatedAt          time.Time                           `json:"created_at"`
	UpdatedAt          time.Time                           `json:"updated_at"`
}

type StudentAcademicClassSummaryResponse struct {
	ClassID           string  `json:"class_id"`
	ClassName         string  `json:"class_name"`
	ClassCode         string  `json:"class_code"`
	RecordsCount      int     `json:"records_count"`
	CompletedCount    int     `json:"completed_count"`
	AverageTotalScore float64 `json:"average_total_score"`
}

type GetMyAcademicRecordsResponse struct {
	StudentID      string                                `json:"student_id"`
	ClassID        string                                `json:"class_id,omitempty"`
	ClassSummaries []StudentAcademicClassSummaryResponse `json:"class_summaries"`
	Records        []StudentAcademicRecordItemResponse   `json:"records"`
}

type StudentAtRiskTopFeatureResponse struct {
	Key          string  `json:"key"`
	Label        string  `json:"label"`
	Value        float64 `json:"value"`
	DisplayValue string  `json:"display_value"`
}

type StudentAtRiskFeatureSummaryResponse struct {
	AttendanceRate28d         float64 `json:"attendance_rate_28d"`
	AbsenceCount28d           float64 `json:"absence_count_28d"`
	AverageTotalScore28d      float64 `json:"average_total_score_28d"`
	HomeworkCompletionRate28d float64 `json:"homework_completion_rate_28d"`
	ActiveEnrollmentCount28d  float64 `json:"active_enrollment_count_28d"`
	WeeklyLessonLoad28d       float64 `json:"weekly_lesson_load_28d"`
	ApprovedLeaveCount28d     float64 `json:"approved_leave_count_28d"`
	DaysSinceLastLesson       float64 `json:"days_since_last_lesson"`
}

type StudentAtRiskPredictionResponse struct {
	StudentID      string                              `json:"student_id"`
	StudentCode    string                              `json:"student_code"`
	StudentName    string                              `json:"student_name"`
	GradeLevel     string                              `json:"grade_level"`
	ClassID        string                              `json:"class_id"`
	ClassCode      string                              `json:"class_code"`
	ClassName      string                              `json:"class_name"`
	SnapshotAt     string                              `json:"snapshot_at"`
	RiskLabel      string                              `json:"risk_label"`
	RiskScore      float64                             `json:"risk_score"`
	RiskBand       string                              `json:"risk_band"`
	ModelName      string                              `json:"model_name"`
	ModelVersion   string                              `json:"model_version"`
	PrimaryReason  string                              `json:"primary_reason"`
	Reasons        []string                            `json:"reasons"`
	TopFeatures    []StudentAtRiskTopFeatureResponse   `json:"top_features"`
	FeatureSummary StudentAtRiskFeatureSummaryResponse `json:"feature_summary"`
}

type GetMyAtRiskPredictionResponse struct {
	StudentID  string                           `json:"student_id"`
	Prediction *StudentAtRiskPredictionResponse `json:"prediction,omitempty"`
}

type CreateMyLeaveRequestRequest struct {
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

type StudentLeaveRequestStudentResponse struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"full_name"`
}

type StudentLeaveRequestClassResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type StudentLeaveRequestLessonResponse struct {
	ID        string    `json:"id"`
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
}

type StudentLeaveRequestResponse struct {
	ID              string                             `json:"id"`
	Student         StudentLeaveRequestStudentResponse `json:"student"`
	LeaveType       string                             `json:"leave_type"`
	ApplyDate       time.Time                          `json:"apply_date"`
	LateMinutes     int                                `json:"late_minutes"`
	EarlyMinutes    int                                `json:"early_minutes"`
	Reason          string                             `json:"reason"`
	Documents       []string                           `json:"documents"`
	Class           *StudentLeaveRequestClassResponse  `json:"class,omitempty"`
	Lesson          *StudentLeaveRequestLessonResponse `json:"lesson,omitempty"`
	Subject         string                             `json:"subject"`
	Status          string                             `json:"status"`
	ApprovedByID    *string                            `json:"approved_by_id,omitempty"`
	ApprovedAt      *time.Time                         `json:"approved_at,omitempty"`
	RejectionReason string                             `json:"rejection_reason"`
	CreatedAt       time.Time                          `json:"created_at"`
	UpdatedAt       time.Time                          `json:"updated_at"`
}

type ListMyLeaveRequestsResponse struct {
	Requests []StudentLeaveRequestResponse `json:"requests"`
}

type CancelMyLeaveRequestResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}
