package student

import (
	"time"
)

type CreateStudentRequest struct {
	Code          string     `json:"code" binding:"required"`
	FullName      string     `json:"full_name" binding:"required"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	GuardianPhone string     `json:"guardian_phone"`
	GradeLevel    string     `json:"grade_level"`
	Status        string     `json:"status"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	Gender        string     `json:"gender"`
	Address       string     `json:"address"`
}

type UpdateStudentRequest struct {
	Code          string     `json:"code"`
	FullName      string     `json:"full_name"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	GuardianPhone string     `json:"guardian_phone"`
	GradeLevel    string     `json:"grade_level"`
	Status        string     `json:"status"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	Gender        string     `json:"gender"`
	Address       string     `json:"address"`
}
