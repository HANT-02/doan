package repositoryinterface

import (
	"doan/internal/entities"
	"doan/internal/repositories"
)

type ShiftRepository interface {
	repositories.BaseRepository[entities.Shift]
}
