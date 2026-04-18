package teacherportal

import "errors"

var (
	ErrTeacherAccessDenied   = errors.New("teacher access denied")
	ErrTeacherProfileMissing = errors.New("teacher profile missing")
)
