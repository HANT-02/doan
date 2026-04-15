package repositoryinterface

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
)

type LessonSummaryRepository interface {
	repositories.BaseRepository[entities.LessonSummary]
	GetByLessonID(ctx context.Context, lessonID string) (*entities.LessonSummary, error)
}
