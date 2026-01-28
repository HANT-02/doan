package entities

import (
	"time"
)

// PasswordReset represents a password reset request stored in DB
// We only store a hash of the plaintext token for security.
// Token is one-time-use and expires after ExpiresAt.
//
// Table name: password_resets
//
// Note: GORM tags assume table already exists via SQL migration.

type PasswordReset struct {
	ID        string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    string     `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string     `gorm:"type:varchar(255);unique;not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedIP *string    `gorm:"type:varchar(100)" json:"created_ip"`
	CreatedUA *string    `gorm:"type:varchar(255)" json:"created_ua"`
	CreatedAt time.Time  `gorm:"default:now()" json:"created_at"`
}
