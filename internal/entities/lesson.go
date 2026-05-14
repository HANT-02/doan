package entities

import (
	"time"

	"gorm.io/gorm"
)

const (
	LessonStatusDraft     = "DRAFT"
	LessonStatusPublished = "PUBLISHED"
	LessonStatusHistory   = "HISTORY"
	LessonStatusUnplanned = "UNPLANNED"
)

const (
	LessonChangeReasonInitialSchedulingCommit = "INITIAL_SCHEDULING_COMMIT"
	LessonChangeReasonLegacyBackfill          = "LEGACY_BACKFILL"
)

type Lesson struct {
	ID               string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClassID          string         `gorm:"not null" json:"class_id"`
	Class            Class          `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"class"`
	TeacherID        *string        `json:"teacher_id"`
	Teacher          Teacher        `gorm:"foreignKey:TeacherID" json:"teacher"`
	DateStart        time.Time      `gorm:"not null" json:"date_start"`
	DateEnd          time.Time      `gorm:"not null" json:"date_end"`
	RoomID           *string        `json:"room_id"`
	Room             Room           `gorm:"foreignKey:RoomID" json:"room"`
	Status           string         `gorm:"type:varchar(32);not null;default:'PUBLISHED'" json:"status"`
	PublishedAt      *time.Time     `json:"published_at,omitempty"`
	SourcePreviewRun *string        `gorm:"column:source_preview_run_id;type:varchar(255)" json:"source_preview_run_id,omitempty"`
	ChangeReason     string         `gorm:"type:text" json:"change_reason"`
	Notes            string         `gorm:"type:text" json:"notes"`
	CreatedAt        time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
