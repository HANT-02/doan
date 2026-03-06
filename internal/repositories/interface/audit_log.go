package repositoryinterface

import (
	"doan/internal/entities"
	"doan/internal/repositories"
)

type AuditLogRepository interface {
	repositories.BaseRepository[entities.AuditLog]
}
