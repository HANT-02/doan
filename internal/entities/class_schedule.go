package entities

type ClassSchedule struct {
	ID        string  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClassID   string  `gorm:"not null" json:"class_id"`
	Class     Class   `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"-"`
	DayOfWeek string  `gorm:"type:varchar(20);not null" json:"day_of_week"`
	ShiftID   string  `gorm:"type:uuid;not null" json:"shift_id"`
	Shift     Shift   `gorm:"foreignKey:ShiftID" json:"shift"`
	RoomID    *string `json:"room_id"`
	Room      Room    `gorm:"foreignKey:RoomID" json:"room"`
}
