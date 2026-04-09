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

type shiftRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Shift]
	db *gorm.DB
}

func NewShiftRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.ShiftRepository {
	modelRepo := postgres.NewBaseRepository[entities.Shift](log, manager, db, "shifts")
	return &shiftRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		db:             db,
		BaseRepository: modelRepo,
	}
}
