package teacherportal

import (
	"errors"

	lessonactivity "doan/internal/usecases/lessonactivity"
	lessonrecord "doan/internal/usecases/lessonrecord"
)

const (
	TeacherAttendanceStatusUnmarked = -1
	TeacherAttendanceStatusAbsent   = 0
	TeacherAttendanceStatusPresent  = 1
	TeacherAttendanceStatusLate     = 2
	TeacherAttendanceStatusExcused  = 3
)

var ErrInvalidTeacherAttendanceStatus = errors.New("invalid teacher attendance status")

func IsValidTeacherAttendanceStatus(status int) bool {
	switch status {
	case TeacherAttendanceStatusAbsent, TeacherAttendanceStatusPresent, TeacherAttendanceStatusLate, TeacherAttendanceStatusExcused:
		return true
	default:
		return false
	}
}

func buildLessonActor(actor Actor) lessonactivity.LessonActor {
	return lessonactivity.LessonActor{
		Role:   actor.Role,
		Email:  actor.Email,
		UserID: actor.UserID,
	}
}

func buildLessonRecordActor(actor Actor) lessonrecord.LessonActor {
	return lessonrecord.LessonActor{
		Role:   actor.Role,
		Email:  actor.Email,
		UserID: actor.UserID,
	}
}

func mapTeacherAttendanceStatusToInternal(status int) (int, error) {
	switch status {
	case TeacherAttendanceStatusAbsent:
		return lessonactivity.AttendanceStatusAbsent, nil
	case TeacherAttendanceStatusPresent:
		return lessonactivity.AttendanceStatusPresent, nil
	case TeacherAttendanceStatusLate:
		return lessonactivity.AttendanceStatusLate, nil
	case TeacherAttendanceStatusExcused:
		return lessonactivity.AttendanceStatusExcused, nil
	default:
		return 0, ErrInvalidTeacherAttendanceStatus
	}
}

func mapInternalAttendanceStatusToTeacher(status int) int {
	switch status {
	case lessonactivity.AttendanceStatusAbsent:
		return TeacherAttendanceStatusAbsent
	case lessonactivity.AttendanceStatusPresent:
		return TeacherAttendanceStatusPresent
	case lessonactivity.AttendanceStatusLate:
		return TeacherAttendanceStatusLate
	case lessonactivity.AttendanceStatusExcused:
		return TeacherAttendanceStatusExcused
	case lessonactivity.AttendanceStatusEarly:
		// Legacy value retained as "late/irregular" in the teacher portal to avoid breaking old rows.
		return TeacherAttendanceStatusLate
	default:
		return TeacherAttendanceStatusUnmarked
	}
}

func isTeacherAttendancePresent(status int) bool {
	return status == TeacherAttendanceStatusPresent || status == TeacherAttendanceStatusLate
}
