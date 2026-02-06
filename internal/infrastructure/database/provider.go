package database

import (
	"context"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/infrastructure/database/postgres/implement"
	"doan/pkg/config"
	"doan/pkg/logger"

	"github.com/google/wire"
	"gorm.io/gorm"
)

var DBProvider = wire.NewSet(
	ProvideDB,
	postgres.NewMigration,
	implement.NewUserRepository,
	implement.NewPasswordResetRepository,
)

// ProvideDB wraps GetDBContext and panics on error (for Wire)
func ProvideDB(ctx context.Context, log logger.Logger, cfg config.Manager) *gorm.DB {
	db, err := postgres.GetDBContext(ctx, log, cfg)
	if err != nil {
		panic(err)
	}
	return db
}
