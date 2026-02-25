package repositoryinterface

import (
	"doan/internal/entities"
	"doan/internal/repositories"
)

type StudentRepository interface {
	repositories.BaseRepository[entities.Student]
}
