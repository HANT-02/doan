package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id" db:"id"`
	Code      string         `gorm:"type:varchar(50);unique" json:"code"`
	FullName  string         `gorm:"type:varchar(255)" json:"fullName"`
	Email     string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string         `gorm:"type:text" json:"-"`
	Role      string         `gorm:"type:varchar(50);default:'STUDENT'" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"isActive"`
	CreatedAt time.Time      `gorm:"default:now()" json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
