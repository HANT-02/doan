package entities

type ProgramCourse struct {
	ID        string  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ProgramID string  `gorm:"not null" json:"program_id"`
	Program   Program `gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE" json:"-"`
	CourseID  string  `gorm:"not null" json:"course_id"`
	Course    Course  `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"course"`
}
