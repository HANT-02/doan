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

type studentRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Student]
	db *gorm.DB
}

func NewStudentRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.StudentRepository {
	modelRepo := postgres.NewBaseRepository[entities.Student](log, manager, db, "students")
	return &studentRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}
