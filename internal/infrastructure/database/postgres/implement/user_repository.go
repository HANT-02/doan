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

type userRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.User]
	db *gorm.DB
}

func NewUserRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.UserRepository {
	modelRepo := postgres.NewBaseRepository[entities.User](log, manager, db, "users")
	return &userRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}
