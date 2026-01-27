package entities

import "time"

type Consultation struct {
	ID         string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	FullName   string    `gorm:"type:varchar(255);not null" json:"full_name"`
	Phone      string    `gorm:"type:varchar(20);not null" json:"phone"`
	GradeLevel string    `gorm:"type:varchar(50);not null" json:"grade_level"`
	Notes      string    `gorm:"type:text" json:"notes"`
	Status     string    `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	CreatedAt  time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
