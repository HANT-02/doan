package entities

import (
	"gorm.io/gorm"
	"time"
)

type Teacher struct {
	ID              string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code            string         `gorm:"type:varchar(50);unique" json:"code"`
	FullName        string         `gorm:"type:varchar(255)" json:"full_name"`
	Email           string         `gorm:"type:varchar(255)" json:"email"`
	Phone           string         `gorm:"type:varchar(20)" json:"phone"`
	IsSchoolTeacher bool           `gorm:"default:false" json:"is_school_teacher"`
	SchoolName      string         `gorm:"type:varchar(255)" json:"school_name"`
	EmploymentType  string         `gorm:"type:varchar(50);default:'PART_TIME'" json:"employment_type"`
	Status          string         `gorm:"type:varchar(50);default:'ACTIVE'" json:"status"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
