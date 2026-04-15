package lessonactivity

import "doan/internal/entities"

const (
	AttendanceStatusPresent = 1
	AttendanceStatusAbsent  = 2
	AttendanceStatusExcused = 3
	AttendanceStatusLate    = 4
	AttendanceStatusEarly   = 5
)

type LessonActor struct {
	Role   string
	Email  string
	UserID string
}

type AttendanceRecordItem struct {
	Student    entities.Student     `json:"student"`
	Attendance *entities.Attendance `json:"attendance,omitempty"`
	Status     int                  `json:"status"`
	Note       string               `json:"note"`
}

func IsValidAttendanceStatus(status int) bool {
	switch status {
	case AttendanceStatusPresent, AttendanceStatusAbsent, AttendanceStatusExcused, AttendanceStatusLate, AttendanceStatusEarly:
		return true
	default:
		return false
	}
}
