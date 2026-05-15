package entities

import "time"

type TeacherSkill struct {
	TeacherID string    `gorm:"type:uuid;primaryKey" json:"teacher_id"`
	SkillID   string    `gorm:"type:uuid;primaryKey" json:"skill_id"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
}

func (TeacherSkill) TableName() string {
	return "teacher_skills"
}
