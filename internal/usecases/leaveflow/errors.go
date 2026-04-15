package leaveflow

import "errors"

var (
	ErrStudentNotFound        = errors.New("student not found")
	ErrTeacherNotFound        = errors.New("teacher not found")
	ErrLeaveRequestNotFound   = errors.New("leave request not found")
	ErrLeaveRequestForbidden  = errors.New("leave request access denied")
	ErrInvalidLeaveType       = errors.New("invalid leave type")
	ErrLeaveRequestNotPending = errors.New("leave request is not pending")
	ErrStudentNotInClass      = errors.New("student is not enrolled in class")
)
