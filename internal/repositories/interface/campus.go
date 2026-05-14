package repositoryinterface

import (
	"doan/internal/entities"
	"doan/internal/repositories"
)

type CampusRepository interface {
	repositories.BaseRepository[entities.Campus]
}
