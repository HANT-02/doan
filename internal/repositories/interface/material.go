package repositoryinterface

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
)

type MaterialFilter struct {
	TeacherID string
	Status    string
	Queue     string
}

type MaterialRepository interface {
	repositories.BaseRepository[entities.Material]
	GetDetailed(ctx context.Context, id string) (*entities.Material, error)
	ListDetailed(ctx context.Context, filter MaterialFilter) ([]entities.Material, error)
}
