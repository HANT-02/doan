package repositoryinterface

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
	"gorm.io/gorm"
	"time"
)

type UserRepository interface {
	repositories.BaseRepository[entities.User]
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	CreateTx(ctx context.Context, tx *gorm.DB, user *entities.User) error
	CreateOTPTx(ctx context.Context, tx *gorm.DB, userID string, purpose string, otpHash string, expiredAt time.Time, pendingPasswordHash *string) error
	GetActiveOTPByUserIDAndPurposeTx(ctx context.Context, tx *gorm.DB, userID string, purpose string) (*entities.UserOTP, error)
	MarkOTPUsedTx(ctx context.Context, tx *gorm.DB, otpID string, usedAt time.Time) error
	ActivateUserTx(ctx context.Context, tx *gorm.DB, userID string) error
	UpdatePasswordTx(ctx context.Context, tx *gorm.DB, userID string, passwordHash string) error
}
