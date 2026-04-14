package repositoryinterface

import (
	"context"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
)

type LessonRepository interface {
	repositories.BaseRepository[entities.Lesson]
	FindOverlappingLessons(
		ctx context.Context,
		from time.Time,
		to time.Time,
		classIDs []string,
		teacherIDs []string,
		roomIDs []string,
	) ([]entities.Lesson, error)
	GetLessonWithRelations(ctx context.Context, id string) (*entities.Lesson, error)
}
