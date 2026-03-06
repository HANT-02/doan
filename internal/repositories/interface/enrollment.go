package repositoryinterface

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
)

type EnrollmentRepository interface {
	repositories.BaseRepository[entities.Enrollment]
	ListByClassID(ctx context.Context, classID string) ([]entities.Enrollment, error)
}
