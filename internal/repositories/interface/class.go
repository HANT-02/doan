package repositoryinterface

import (
	"doan/internal/entities"
	"doan/internal/repositories"
)

type ClassRepository interface {
	repositories.BaseRepository[entities.Class]
}
