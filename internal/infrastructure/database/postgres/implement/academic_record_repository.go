package implement

import (
	"context"

	"doan/internal/entities"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/base_struct"
	"doan/pkg/config"
	"doan/pkg/logger"

	"gorm.io/gorm"
)

type academicRecordRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.AcademicRecord]
	db *gorm.DB
}

func NewAcademicRecordRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.AcademicRecordRepository {
	modelRepo := postgres.NewBaseRepository[entities.AcademicRecord](log, manager, db, "academic_records")
	return &academicRecordRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (r *academicRecordRepository) ListByLessonSummaryID(ctx context.Context, lessonSummaryID string) ([]entities.AcademicRecord, error) {
	records := make([]entities.AcademicRecord, 0)
	err := r.db.WithContext(ctx).
		Preload("Student").
		Where("lesson_summary_id = ?", lessonSummaryID).
		Order("created_at ASC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *academicRecordRepository) ListByStudentID(ctx context.Context, studentID string) ([]entities.AcademicRecord, error) {
	records := make([]entities.AcademicRecord, 0)
	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("LessonSummary").
		Preload("LessonSummary.Lesson").
		Preload("LessonSummary.Lesson.Class").
		Preload("LessonSummary.Lesson.Teacher").
		Preload("LessonSummary.Lesson.Room").
		Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *academicRecordRepository) GetByLessonSummaryAndStudent(ctx context.Context, lessonSummaryID, studentID string) (*entities.AcademicRecord, error) {
	var record entities.AcademicRecord
	err := r.db.WithContext(ctx).
		Preload("Student").
		Where("lesson_summary_id = ? AND student_id = ?", lessonSummaryID, studentID).
		First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}
