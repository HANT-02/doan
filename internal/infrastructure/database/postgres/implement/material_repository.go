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

type materialRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Material]
	db *gorm.DB
}

func NewMaterialRepository(db *gorm.DB, log logger.Logger, manager config.Manager) repositoryinterface.MaterialRepository {
	modelRepo := postgres.NewBaseRepository[entities.Material](log, manager, db, "materials")
	return &materialRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (r *materialRepository) GetDetailed(ctx context.Context, id string) (*entities.Material, error) {
	var material entities.Material
	err := r.db.WithContext(ctx).
		Preload("LatestLabel").
		Preload("AuditLogs", func(db *gorm.DB) *gorm.DB { return db.Order("triggered_at DESC") }).
		Preload("AuditLogs.Label").
		Preload("ApprovalDecisions", func(db *gorm.DB) *gorm.DB { return db.Order("decided_at DESC") }).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&material).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (r *materialRepository) ListDetailed(ctx context.Context, filter repositoryinterface.MaterialFilter) ([]entities.Material, error) {
	query := r.db.WithContext(ctx).
		Model(&entities.Material{}).
		Preload("LatestLabel").
		Preload("AuditLogs", func(db *gorm.DB) *gorm.DB { return db.Order("triggered_at DESC") }).
		Preload("AuditLogs.Label").
		Preload("ApprovalDecisions", func(db *gorm.DB) *gorm.DB { return db.Order("decided_at DESC") }).
		Where("deleted_at IS NULL")

	if filter.TeacherID != "" {
		query = query.Where("teacher_id = ?", filter.TeacherID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Queue == "flagged" {
		query = query.Where("status = ?", "AI_REVIEWED")
	}

	var materials []entities.Material
	err := query.Order("uploaded_at DESC").Find(&materials).Error
	if err != nil {
		return nil, err
	}
	return materials, nil
}
