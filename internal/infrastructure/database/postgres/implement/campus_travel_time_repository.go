package implement

import (
	"doan/internal/entities"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/pkg/base_struct"
	"doan/pkg/config"
	"doan/pkg/logger"

	"gorm.io/gorm"
)

type campusTravelTimeRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.CampusTravelTime]
	db *gorm.DB
}

func NewCampusTravelTimeRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repositoryinterface.CampusTravelTimeRepository {
	modelRepo := postgres.NewBaseRepository[entities.CampusTravelTime](log, manager, db, "campus_travel_times")
	return &campusTravelTimeRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}
