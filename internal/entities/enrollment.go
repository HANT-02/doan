package entities

import "time"

// Table 3.12
type Enrollment struct {
	ID         string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClassID    string     `gorm:"not null" json:"class_id"`
	Class      Class      `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"class"`
	StudentID  string     `gorm:"not null" json:"student_id"`
	Student    Student    `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student"`
	Status     string     `gorm:"type:varchar(50);default:'APPLIED'" json:"status"`
	ApprovedAt *time.Time `json:"approved_at"`
	RejectedAt *time.Time `json:"rejected_at"`
	CreatedAt  time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
