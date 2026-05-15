package entities

import (
	"time"

	"gorm.io/gorm"
)

type Skill struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code      string         `gorm:"type:varchar(120);uniqueIndex;not null" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Status    string         `gorm:"type:varchar(50);default:'ACTIVE'" json:"status"`
	CreatedAt time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
