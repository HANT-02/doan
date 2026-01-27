package entities

type Objective struct {
	ID        string  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code      string  `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name      string  `gorm:"type:text;not null" json:"name"`
	ProgramID string  `gorm:"not null" json:"program_id"`
	Program   Program `gorm:"foreignKey:ProgramID;constraint:OnDelete:CASCADE" json:"-"`
}
