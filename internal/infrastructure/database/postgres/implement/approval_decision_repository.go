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

type approvalDecisionRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.ApprovalDecision]
	db *gorm.DB
}

func NewApprovalDecisionRepository(db *gorm.DB, log logger.Logger, manager config.Manager) repositoryinterface.ApprovalDecisionRepository {
	modelRepo := postgres.NewBaseRepository[entities.ApprovalDecision](log, manager, db, "approval_decisions")
	return &approvalDecisionRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}
