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

type enrollmentRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Enrollment]
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB, log logger.Logger, manager config.Manager) repointerface.EnrollmentRepository {
	modelRepo := postgres.NewBaseRepository[entities.Enrollment](log, manager, db, "enrollments")
	return &enrollmentRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (r *enrollmentRepository) ListByClassID(ctx context.Context, classID string) ([]entities.Enrollment, error) {
	var enrollments []entities.Enrollment
	err := r.db.WithContext(ctx).Where("class_id = ?", classID).Find(&enrollments).Error
	if err != nil {
		return nil, err
	}
	return enrollments, nil
}
