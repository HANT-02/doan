package implement

import (
	"context"

	"doan/internal/entities"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/pkg/base_struct"
	"doan/pkg/config"
	"doan/pkg/logger"

	"gorm.io/gorm"
)

type labelRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Label]
	db *gorm.DB
}

func NewLabelRepository(db *gorm.DB, log logger.Logger, manager config.Manager) repositoryinterface.LabelRepository {
	modelRepo := postgres.NewBaseRepository[entities.Label](log, manager, db, "labels")
	return &labelRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (r *labelRepository) GetByCode(ctx context.Context, code string) (*entities.Label, error) {
	var label entities.Label
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&label).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}
