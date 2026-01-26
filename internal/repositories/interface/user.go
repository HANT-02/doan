package repositoryinterface

import (
	"doan/internal/entities"
	"doan/internal/repositories"
)

type UserRepository interface {
	repositories.BaseRepository[entities.User]
}
