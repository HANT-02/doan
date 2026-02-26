package repositoryinterface

import (
	"doan/internal/entities"
	"doan/internal/repositories"
)

type CourseRepository interface {
	repositories.BaseRepository[entities.Course]
}
