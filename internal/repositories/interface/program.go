package repositoryinterface

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
)

type ProgramRepository interface {
	repositories.BaseRepository[entities.Program]
	AddCourses(ctx context.Context, programID string, courseIDs []string) error
	RemoveCourses(ctx context.Context, programID string, courseIDs []string) error
	GetProgramWithCourses(ctx context.Context, id string) (*entities.Program, error)
}
