package implement

import (
	"doan/internal/entities"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/base_struct"
	"doan/pkg/config"
	"doan/pkg/logger"

	"gorm.io/gorm"
)

type classRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Class]
	db *gorm.DB
}

func NewClassRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.ClassRepository {
	modelRepo := postgres.NewBaseRepository[entities.Class](log, manager, db, "classes")
	return &classRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}
