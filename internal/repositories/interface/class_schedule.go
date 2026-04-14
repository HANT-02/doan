package repositoryinterface

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
)

type ClassScheduleRepository interface {
	repositories.BaseRepository[entities.ClassSchedule]
	GetSchedulesByClassID(ctx context.Context, classID string) ([]entities.ClassSchedule, error)
}
