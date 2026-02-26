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

type courseRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Course]
	db *gorm.DB
}

func NewCourseRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.CourseRepository {
	modelRepo := postgres.NewBaseRepository[entities.Course](log, manager, db, "courses")
	return &courseRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}
