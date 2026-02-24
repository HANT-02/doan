package repositoryinterface

import (
	"doan/internal/entities"
	"doan/internal/repositories"
)

type RoomRepository interface {
	repositories.BaseRepository[entities.Room]
}
