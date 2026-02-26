package entities

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID          string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code        string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Notes       string         `gorm:"type:text" json:"notes"`
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	MaxStudents int            `json:"max_students"`
	Status      string         `gorm:"type:varchar(50);default:'OPEN'" json:"status"`
	Price       float64        `gorm:"type:numeric(10,2)" json:"price"`
	ProgramID   *string        `json:"program_id"`
	Program     Program        `gorm:"foreignKey:ProgramID" json:"program"`
	CourseID    *string        `json:"course_id"`
	Course      Course         `gorm:"foreignKey:CourseID" json:"course"`
	TeacherID   *string        `json:"teacher_id"`
	Teacher     Teacher        `gorm:"foreignKey:TeacherID" json:"teacher"`
	CreatedAt   time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
