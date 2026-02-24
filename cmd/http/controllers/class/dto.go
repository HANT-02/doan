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
