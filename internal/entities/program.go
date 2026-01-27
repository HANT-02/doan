package entities

import (
	"gorm.io/gorm"
	"time"
)

type Program struct {
	ID            string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code          string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name          string         `gorm:"type:varchar(255)" json:"name"`
	Track         string         `gorm:"type:varchar(50)" json:"track"`
	EffectiveFrom *time.Time     `json:"effective_from"`
	EffectiveTo   *time.Time     `json:"effective_to"`
	CreatedByID   *string        `json:"created_by_id"`
	CreatedBy     User           `gorm:"foreignKey:CreatedByID" json:"created_by"`
	ApprovedByID  *uint          `json:"approved_by_id"`
	ApprovedBy    User           `gorm:"foreignKey:ApprovedByID" json:"approved_by"`
	ApprovalNote  string         `gorm:"type:text" json:"approval_note"`
	PublishedAt   *time.Time     `json:"published_at"`
	ArchivedAt    *time.Time     `json:"archived_at"`
	CreatedAt     time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
