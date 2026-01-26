package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id" db:"id"`
	UserName  string         `gorm:"type:varchar(255);unique;not null;column:username" json:"username" db:"username"` // Đảm bảo GORM tạo unique constraint cho username
	PassWord  string         `gorm:"type:text;not null;column:password" json:"-" db:"password"`
	Status    string         `gorm:"type:varchar(50);default:'active';not null" json:"status" db:"status"`
	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at" db:"created_at"`
	UpdatedAt *time.Time     `gorm:"default:null" json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;default:null" json:"deleted_at,omitempty" db:"deleted_at"`
}
