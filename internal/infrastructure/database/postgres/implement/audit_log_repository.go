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

type auditLogRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.AuditLog]
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB, log logger.Logger, manager config.Manager) repositoryinterface.AuditLogRepository {
	modelRepo := postgres.NewBaseRepository[entities.AuditLog](log, manager, db, "audit_logs")
	return &auditLogRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}
