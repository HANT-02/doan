package class

import (
	"time"
)

type CreateClassRequest struct {
	Code        string     `json:"code" binding:"required"`
	Name        string     `json:"name" binding:"required"`
	Notes       string     `json:"notes"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
	MaxStudents int        `json:"max_students" binding:"required,min=1"`
	Status      string     `json:"status" binding:"omitempty,oneof=OPEN CLOSED CANCELLED"`
	Price       float64    `json:"price"`
	ProgramID   *string    `json:"program_id"`
	CourseID    *string    `json:"course_id"`
	TeacherID   *string    `json:"teacher_id"`
}

type UpdateClassRequest struct {
	Code        string     `json:"code" binding:"required"`
	Name        string     `json:"name" binding:"required"`
	Notes       string     `json:"notes"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
	MaxStudents int        `json:"max_students" binding:"required,min=1"`
	Status      string     `json:"status" binding:"omitempty,oneof=OPEN CLOSED CANCELLED"`
	Price       float64    `json:"price"`
	ProgramID   *string    `json:"program_id"`
	CourseID    *string    `json:"course_id"`
	TeacherID   *string    `json:"teacher_id"`
}

type EnrollStudentsRequest struct {
	StudentIDs []string `json:"student_ids" binding:"required,min=1"`
}

type RemoveStudentsRequest struct {
	StudentIDs []string `json:"student_ids" binding:"required,min=1"`
}

type ReserveStudentRequest struct {
	Reason      string     `json:"reason"`
	EffectiveAt *time.Time `json:"effective_at"`
}

type ResumeStudentRequest struct {
	Reason      string     `json:"reason"`
	EffectiveAt *time.Time `json:"effective_at"`
}

type TransferStudentRequest struct {
	TargetClassID string     `json:"target_class_id" binding:"required"`
	Reason        string     `json:"reason"`
	EffectiveAt   *time.Time `json:"effective_at"`
}

type AssignTeacherRequest struct {
	TeacherID string `json:"teacher_id" binding:"required"`
}

type CreateClassScheduleRequest struct {
	ShiftID   string  `json:"shift_id" binding:"required"`
	DayOfWeek string  `json:"day_of_week" binding:"required"`
	RoomID    *string `json:"room_id"`
}

type ClassResponse struct {
	ID          string     `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Notes       string     `json:"notes"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	MaxStudents int        `json:"max_students"`
	Status      string     `json:"status"`
	Price       float64    `json:"price"`
	ProgramID   *string    `json:"program_id,omitempty"`
	CourseID    *string    `json:"course_id,omitempty"`
	TeacherID   *string    `json:"teacher_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
