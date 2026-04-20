package studentportal

import "errors"

var (
	ErrStudentAccessDenied = errors.New("student access denied")
	ErrStudentNotFound     = errors.New("student not found")
)
