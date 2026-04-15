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

type lessonSummaryRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.LessonSummary]
	db *gorm.DB
}

func NewLessonSummaryRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.LessonSummaryRepository {
	modelRepo := postgres.NewBaseRepository[entities.LessonSummary](log, manager, db, "lesson_summaries")
	return &lessonSummaryRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (r *lessonSummaryRepository) GetByLessonID(ctx context.Context, lessonID string) (*entities.LessonSummary, error) {
	var summary entities.LessonSummary
	err := r.db.WithContext(ctx).
		Preload("CreatedBy").
		Where("lesson_id = ?", lessonID).
		First(&summary).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &summary, nil
}
