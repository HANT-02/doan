package entities

import "time"

type AuditLog struct {
	ID              string     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	MaterialID      string     `gorm:"type:uuid;not null" json:"material_id"`
	Material        Material   `gorm:"foreignKey:MaterialID;constraint:OnDelete:CASCADE" json:"-"`
	LabelID         *string    `json:"label_id"`
	Label           Label      `gorm:"foreignKey:LabelID" json:"-"`
	Status          string     `gorm:"type:varchar(50);not null;default:'COMPLETED'" json:"status"`
	Provider        string     `gorm:"type:varchar(100);not null" json:"provider"`
	RawOCRText      string     `gorm:"type:text" json:"raw_ocr_text"`
	ConfidenceScore float64    `gorm:"type:numeric(5,4)" json:"confidence_score"`
	Reasoning       string     `gorm:"type:text" json:"reasoning"`
	DetectedIssues  string     `gorm:"type:jsonb;default:'[]'" json:"detected_issues"`
	TriggeredAt     time.Time  `gorm:"default:now()" json:"triggered_at"`
	CompletedAt     *time.Time `json:"completed_at"`
	CreatedAt       time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
