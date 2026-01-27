package entities

import "time"

// Table 3.15
type LessonSummary struct {
	ID               string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	LessonID         string    `gorm:"unique;not null" json:"lesson_id"`
	Lesson           Lesson    `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"lesson"`
	Topic            string    `gorm:"type:text" json:"topic"`
	LessonContent    string    `gorm:"type:text" json:"lesson_content"`
	ClassFeedback    string    `gorm:"type:text" json:"class_feedback"`
	Homework         string    `gorm:"type:text" json:"homework"`
	HomeworkDeadline time.Time `json:"homework_deadline"`
	TeacherNotes     string    `gorm:"type:text" json:"teacher_notes"`
	CreatedByID      *string   `json:"created_by_id"`
	CreatedBy        User      `gorm:"foreignKey:CreatedByID" json:"created_by"`
	CreatedAt        time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
