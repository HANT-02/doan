package entities

import "time"

type Lesson struct {
	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClassID   string    `gorm:"not null" json:"class_id"`
	Class     Class     `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"class"`
	TeacherID *string   `json:"teacher_id"`
	Teacher   Teacher   `gorm:"foreignKey:TeacherID" json:"teacher"`
	DateStart time.Time `gorm:"not null" json:"date_start"`
	DateEnd   time.Time `gorm:"not null" json:"date_end"`
	RoomID    *string   `json:"room_id"`
	Room      Room      `gorm:"foreignKey:RoomID" json:"room"`
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
