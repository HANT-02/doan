package entities

import "time"

type CourseRequiredSkill struct {
	CourseID  string    `gorm:"type:uuid;primaryKey" json:"course_id"`
	SkillID   string    `gorm:"type:uuid;primaryKey" json:"skill_id"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
}

func (CourseRequiredSkill) TableName() string {
	return "course_required_skills"
}
