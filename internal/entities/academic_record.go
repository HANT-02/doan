package entities

import "time"

type AcademicRecord struct {
	ID                 string        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	LessonSummaryID    string        `gorm:"not null" json:"lesson_summary_id"`
	LessonSummary      LessonSummary `gorm:"foreignKey:LessonSummaryID;constraint:OnDelete:CASCADE" json:"-"`
	StudentID          string        `gorm:"not null" json:"student_id"`
	Student            Student       `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student"`
	HomeworkCompleted  bool          `gorm:"default:false" json:"homework_completed"`
	HomeworkScore      float64       `gorm:"type:numeric(5,2)" json:"homework_score"`
	AttitudeRating     int           `json:"attitude_rating"`
	ParticipationScore float64       `gorm:"type:numeric(5,2)" json:"participation_score"`
	PersonalComment    string        `gorm:"type:text" json:"personal_comment"`
	TotalScore         float64       `gorm:"type:numeric(5,2)" json:"total_score"`
	IsCompleted        bool          `gorm:"default:false" json:"is_completed"`
	CreatedAt          time.Time     `gorm:"default:now()" json:"created-at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}
