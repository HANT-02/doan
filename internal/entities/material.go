package entities

import (
	"time"

	"gorm.io/gorm"
)

type Material struct {
	ID                string             `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TeacherID         string             `gorm:"type:uuid;not null" json:"teacher_id"`
	Teacher           Teacher            `gorm:"foreignKey:TeacherID" json:"-"`
	Title             string             `gorm:"type:varchar(255);not null" json:"title"`
	Description       string             `gorm:"type:text" json:"description"`
	FileName          string             `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath          string             `gorm:"type:text;not null" json:"file_path"`
	FileType          string             `gorm:"type:varchar(100);not null" json:"file_type"`
	Status            string             `gorm:"type:varchar(50);not null;default:'UPLOADED'" json:"status"`
	LatestLabelID     *string            `json:"latest_label_id"`
	LatestLabel       Label              `gorm:"foreignKey:LatestLabelID" json:"-"`
	UploadedAt        time.Time          `gorm:"default:now()" json:"uploaded_at"`
	AuditLogs         []AuditLog         `gorm:"foreignKey:MaterialID" json:"-"`
	ApprovalDecisions []ApprovalDecision `gorm:"foreignKey:MaterialID" json:"-"`
	CreatedAt         time.Time          `gorm:"default:now()" json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	DeletedAt         gorm.DeletedAt     `gorm:"index" json:"-"`
}
