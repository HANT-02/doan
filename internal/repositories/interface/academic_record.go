package repositoryinterface

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
)

type AcademicRecordRepository interface {
	repositories.BaseRepository[entities.AcademicRecord]
	ListByLessonSummaryID(ctx context.Context, lessonSummaryID string) ([]entities.AcademicRecord, error)
	ListByStudentID(ctx context.Context, studentID string) ([]entities.AcademicRecord, error)
	GetByLessonSummaryAndStudent(ctx context.Context, lessonSummaryID, studentID string) (*entities.AcademicRecord, error)
}
