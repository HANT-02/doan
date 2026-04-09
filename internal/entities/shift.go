package entities

import "time"

type Shift struct {
	ID              string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code            string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name            string    `gorm:"type:varchar(255);not null" json:"name"`
	StartTime       string    `gorm:"type:varchar(10);not null" json:"start_time"`
	EndTime         string    `gorm:"type:varchar(10);not null" json:"end_time"`
	DurationMinutes int       `gorm:"not null" json:"duration_minutes"`
	SessionType     string    `gorm:"type:varchar(50);not null" json:"session_type"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	Notes           string    `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
