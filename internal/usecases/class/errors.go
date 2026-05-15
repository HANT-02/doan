package class

import "errors"

var (
	ErrEnrollmentNotFound         = errors.New("active enrollment not found")
	ErrReservedEnrollmentNotFound = errors.New("reserved enrollment not found")
	ErrTargetClassRequired        = errors.New("target_class_id is required")
	ErrTransferSameClass          = errors.New("target class must be different from source class")
	ErrTargetClassNotOpen         = errors.New("target class is not open")
	ErrTargetClassIncompatible    = errors.New("target class is not compatible with source class")
	ErrTargetClassFull            = errors.New("target class has no remaining capacity")
	ErrStudentAlreadyInTarget     = errors.New("student already belongs to target class")
	ErrStudentScheduleConflict    = errors.New("target class conflicts with student timetable")
	ErrStudentTravelConflict      = errors.New("target class violates student travel gap")
)
