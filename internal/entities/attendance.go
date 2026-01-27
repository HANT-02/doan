package entities

import "time"

type Attendance struct {
	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	LessonID  string    `gorm:"not null" json:"lesson_id"`
	Lesson    Lesson    `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"lesson"`
	StudentID string    `gorm:"not null" json:"student_id"`
	Student   Student   `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student"`
	Status    int       `gorm:"not null" json:"status"`
	Note      string    `gorm:"type:text" json:"note"`
	MarkedAt  time.Time `json:"marked_at"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
