package entities

import (
	"time"

	"gorm.io/gorm"
)

type PasswordReset struct {
	ID          string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      string         `gorm:"type:uuid;not null" json:"user_id"`
	TokenHash   string         `gorm:"type:text;not null" json:"-"` // Store hashed token
	ExpiresAt   time.Time      `gorm:"not null" json:"expires_at"`
	UsedAt      *time.Time     `json:"used_at"` // Nullable, to mark token as used
	RequestedIP string         `gorm:"type:varchar(45)" json:"requested_ip"`
	UserAgent   string         `gorm:"type:text" json:"user_agent"`
	CreatedAt   time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"default:now()" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"` // If soft delete is consistent
}
