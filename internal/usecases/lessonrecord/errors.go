package lessonrecord

import "errors"

var (
	ErrLessonNotFound     = errors.New("lesson not found")
	ErrLessonAccessDenied = errors.New("lesson access denied")
	ErrStudentNotFound    = errors.New("student not found")
	ErrInvalidRecordRow   = errors.New("invalid academic record row")
)
