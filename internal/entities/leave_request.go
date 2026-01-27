package entities

import (
	"github.com/lib/pq"
	"time"
)

type LeaveRequest struct {
	ID              string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	StudentID       string         `gorm:"not null" json:"student_id"`
	Student         Student        `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student"`
	LeaveType       string         `gorm:"type:varchar(50);not null" json:"leave_type"` // LEAVE, LATE, EARLY
	ApplyDate       time.Time      `gorm:"not null" json:"apply_date"`
	LateMinutes     int            `json:"late_minutes"`
	EarlyMinutes    int            `json:"early_minutes"`
	Reason          string         `gorm:"type:text;not null" json:"reason"`
	Documents       pq.StringArray `gorm:"type:text[]" json:"documents"`
	ClassID         *string        `json:"class_id"`
	Class           Class          `gorm:"foreignKey:ClassID;constraint:OnDelete:SET NULL" json:"class"`
	LessonID        *string        `json:"lesson_id"`
	Lesson          Lesson         `gorm:"foreignKey:LessonID;constraint:OnDelete:SET NULL" json:"lesson"`
	Subject         string         `gorm:"type:varchar(255)" json:"subject"`
	Status          string         `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	ApprovedByID    *string        `json:"approved_by_id"`
	ApprovedBy      User           `gorm:"foreignKey:ApprovedByID" json:"approver"`
	ApprovedAt      *time.Time     `json:"approved_at"`
	RejectionReason string         `gorm:"type:text" json:"rejection_reason"`
	CreatedAt       time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
