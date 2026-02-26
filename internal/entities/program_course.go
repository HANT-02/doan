package entities

type ProgramCourse struct {
	ID        string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProgramID string `gorm:"type:uuid;not null" json:"program_id"`
	CourseID  string `gorm:"type:uuid;not null" json:"course_id"`
}
