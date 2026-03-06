package entities

import "time"

type Label struct {
	ID          string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code        string    `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Severity    string    `gorm:"type:varchar(50);not null" json:"severity"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
