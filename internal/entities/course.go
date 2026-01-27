package entities

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	ID                     string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code                   string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name                   string         `gorm:"type:varchar(255)" json:"name"`
	Description            string         `gorm:"type:text" json:"description"`
	GradeLevel             string         `gorm:"type:varchar(50)" json:"grade_level"`
	Subject                string         `gorm:"type:varchar(255)" json:"subject"`
	SessionCount           int            `json:"session_count"`
	SessionDurationMinutes int            `json:"session_duration_minutes"`
	TotalHours             float64        `gorm:"type:numeric(8,2)" json:"total_hours"`
	Price                  float64        `gorm:"type:numeric(10,2)" json:"price"`
	Status                 string         `gorm:"type:varchar(50);default:'ACTIVE'" json:"status"`
	CreatedAt              time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
