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

type campusRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Campus]
	db *gorm.DB
}

func NewCampusRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repositoryinterface.CampusRepository {
	modelRepo := postgres.NewBaseRepository[entities.Campus](log, manager, db, "campuses")
	return &campusRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}
