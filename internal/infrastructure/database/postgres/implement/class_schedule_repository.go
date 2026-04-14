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

type classScheduleRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.ClassSchedule]
	db *gorm.DB
}

func NewClassScheduleRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.ClassScheduleRepository {
	modelRepo := postgres.NewBaseRepository[entities.ClassSchedule](log, manager, db, "class_schedules")
	return &classScheduleRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (r *classScheduleRepository) GetSchedulesByClassID(ctx context.Context, classID string) ([]entities.ClassSchedule, error) {
	var schedules []entities.ClassSchedule
	err := r.db.WithContext(ctx).
		Where("class_id = ?", classID).
		Preload("Shift").
		Preload("Room").
		Find(&schedules).Error
	return schedules, err
}
