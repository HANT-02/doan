package repositoryinterface

import (
	"context"
	"doan/internal/entities"
	"time"

	"gorm.io/gorm"
)

type PasswordResetRepository interface {
	CreateTx(ctx context.Context, tx *gorm.DB, reset *entities.PasswordReset) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*entities.PasswordReset, error)
	MarkAsUsedTx(ctx context.Context, tx *gorm.DB, id string, usedAt time.Time) error
}
