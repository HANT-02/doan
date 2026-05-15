package entities

import (
	"time"

	"gorm.io/gorm"
)

type CampusTravelTime struct {
	ID            string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	FromCampusID  string         `gorm:"not null;index:idx_campus_travel_pair,unique" json:"from_campus_id"`
	FromCampus    Campus         `gorm:"foreignKey:FromCampusID;constraint:OnDelete:CASCADE" json:"from_campus"`
	ToCampusID    string         `gorm:"not null;index:idx_campus_travel_pair,unique" json:"to_campus_id"`
	ToCampus      Campus         `gorm:"foreignKey:ToCampusID;constraint:OnDelete:CASCADE" json:"to_campus"`
	TravelMinutes int            `gorm:"not null" json:"travel_minutes"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
