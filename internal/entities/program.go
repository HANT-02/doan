package entities

import (
	"time"

	"gorm.io/gorm"
)

type Program struct {
	ID            string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code          string         `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name          string         `gorm:"type:varchar(255)" json:"name"`
	Track         string         `gorm:"type:varchar(50)" json:"track"` // SUPPORT, BASIC, ADVANCED
	EffectiveFrom *time.Time     `json:"effective_from"`
	EffectiveTo   *time.Time     `json:"effective_to"`
	CreatedByID   *string        `gorm:"type:uuid" json:"created_by_id"`
	ApprovedByID  *string        `gorm:"type:uuid" json:"approved_by_id"`
	ApprovalNote  string         `gorm:"type:text" json:"approval_note"`
	PublishedAt   *time.Time     `json:"published_at"`
	ArchivedAt    *time.Time     `json:"archived_at"`
	CreatedAt     time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Relationships
	Courses []Course `gorm:"many2many:program_courses;" json:"courses"`
}
