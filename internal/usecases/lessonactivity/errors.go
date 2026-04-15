package lessonactivity

import "errors"

var (
	ErrLessonNotFound        = errors.New("lesson not found")
	ErrLessonAccessDenied    = errors.New("lesson access denied")
	ErrTeacherProfileMissing = errors.New("teacher profile not found")
	ErrInvalidAttendanceRow  = errors.New("invalid attendance row")
)
