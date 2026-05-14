package repositoryinterface

import (
	"doan/internal/entities"
	"doan/internal/repositories"
)

type CampusTravelTimeRepository interface {
	repositories.BaseRepository[entities.CampusTravelTime]
}
