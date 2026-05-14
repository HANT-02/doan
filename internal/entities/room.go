package entities

import (
	"gorm.io/gorm"
	"time"
)

type Room struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code      string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Capacity  int            `json:"capacity"`
	Address   string         `gorm:"type:text" json:"address"`
	CampusID  *string        `json:"campus_id,omitempty"`
	Campus    Campus         `gorm:"foreignKey:CampusID" json:"campus,omitempty"`
	CreatedAt time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
