package entities

import (
	"gorm.io/gorm"
	"time"
)

type ApprovalDecision struct {
	ID                  string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	MaterialID          string    `gorm:"type:uuid;not null" json:"material_id"`
	Material            Material  `gorm:"foreignKey:MaterialID;constraint:OnDelete:CASCADE" json:"-"`
	AuditLogID          *string   `json:"audit_log_id"`
	AuditLog            AuditLog  `gorm:"foreignKey:AuditLogID" json:"-"`
	ComplianceOfficerID string    `gorm:"type:uuid;not null" json:"compliance_officer_id"`
	Approved            bool      `gorm:"not null" json:"approved"`
	RejectReason        string    `gorm:"type:text" json:"reject_reason"`
	Notes               string    `gorm:"type:text" json:"notes"`
	DecidedAt           time.Time `gorm:"default:now()" json:"decided_at"`
	CreatedAt           time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DeletedAt           gorm.DeletedAt
}
