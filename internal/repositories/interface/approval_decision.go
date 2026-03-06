package repositoryinterface

import (
	"doan/internal/entities"
	"doan/internal/repositories"
)

type ApprovalDecisionRepository interface {
	repositories.BaseRepository[entities.ApprovalDecision]
}
