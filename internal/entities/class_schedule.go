package entities

type ClassSchedule struct {
	ID        string  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClassID   string  `gorm:"not null" json:"class_id"`
	Class     Class   `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"-"`
	DayOfWeek string  `gorm:"type:varchar(20);not null" json:"day_of_week"`
	StartTime string  `gorm:"type:varchar(10);not null" json:"start_time"`
	EndTime   string  `gorm:"type:varchar(10);not null" json:"end_time"`
	RoomID    *string `json:"room_id"`
	Room      Room    `gorm:"foreignKey:RoomID" json:"room"`
}
