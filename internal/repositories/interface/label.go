package repositoryinterface

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
)

type LabelRepository interface {
	repositories.BaseRepository[entities.Label]
	GetByCode(ctx context.Context, code string) (*entities.Label, error)
}
