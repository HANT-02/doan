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

type roomRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Room]
	db *gorm.DB
}

func NewRoomRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.RoomRepository {
	modelRepo := postgres.NewBaseRepository[entities.Room](log, manager, db, "rooms")
	return &roomRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}
