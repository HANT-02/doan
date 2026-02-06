package implement

import (
	"context"
	"doan/internal/entities"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
	"time"

	"gorm.io/gorm"
)

type passwordResetRepository struct {
	db  *gorm.DB
	log logger.Logger
}

func NewPasswordResetRepository(db *gorm.DB, log logger.Logger) repositoryinterface.PasswordResetRepository {
	return &passwordResetRepository{
		db:  db,
		log: log,
	}
}

func (r *passwordResetRepository) CreateTx(ctx context.Context, tx *gorm.DB, reset *entities.PasswordReset) error {
	return tx.WithContext(ctx).Create(reset).Error
}

func (r *passwordResetRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entities.PasswordReset, error) {
	var reset entities.PasswordReset
	// Find token that matches hash, is NOT used, and is NOT expired
	// Note: We intentionally don't filter by expires_at in query to distinguish between "invalid" and "expired" if we wanted,
	// but strictly speaking, an expired token is invalid for use.
	// For security, strict filtering is better.
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		Where("used_at IS NULL").
		Where("expires_at > ?", time.Now()).
		First(&reset).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if valid token not found
		}
		return nil, err
	}
	return &reset, nil
}

func (r *passwordResetRepository) MarkAsUsedTx(ctx context.Context, tx *gorm.DB, id string, usedAt time.Time) error {
	return tx.WithContext(ctx).
		Model(&entities.PasswordReset{}).
		Where("id = ?", id).
		Update("used_at", usedAt).Error
}
