package entities

import (
	"time"

	"gorm.io/gorm"
)

type Campus struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code      string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Address   string         `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
