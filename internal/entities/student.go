package entities

import (
	"gorm.io/gorm"
	"time"
)

type Student struct {
	ID            string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code          string         `gorm:"type:varchar(50);unique" json:"code"`
	FullName      string         `gorm:"type:varchar(255)" json:"full_name"`
	Email         string         `gorm:"type:varchar(255)" json:"email"`
	Phone         string         `gorm:"type:varchar(20)" json:"phone"`
	GuardianPhone string         `gorm:"type:varchar(20)" json:"guardian_phone"`
	GradeLevel    string         `gorm:"type:varchar(50)" json:"grade_level"`
	Status        string         `gorm:"type:varchar(50);default:'ACTIVE'" json:"status"`
	DateOfBirth   *time.Time     `json:"date_of_birth"`
	Gender        string         `gorm:"type:varchar(20)" json:"gender"`
	Address       string         `gorm:"type:text" json:"address"`
	CreatedAt     time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
