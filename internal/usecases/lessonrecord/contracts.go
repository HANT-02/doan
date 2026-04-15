package lessonrecord

import "doan/internal/entities"

type LessonActor struct {
	Role   string
	Email  string
	UserID string
}

type LessonAcademicRecordItem struct {
	Student entities.Student         `json:"student"`
	Record  *entities.AcademicRecord `json:"record,omitempty"`
}
