package implement

import (
	"context"
	"doan/internal/entities"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/base_struct"
	"doan/pkg/config"
	"doan/pkg/logger"
	"gorm.io/gorm"
	"time"
)

type userRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.User]
	db *gorm.DB
}

func NewUserRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.UserRepository {
	modelRepo := postgres.NewBaseRepository[entities.User](log, manager, db, "users")
	return &userRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (u *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := u.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (u *userRepository) CreateTx(ctx context.Context, tx *gorm.DB, user *entities.User) error {
	return tx.WithContext(ctx).Create(user).Error
}

func (u *userRepository) CreateOTPTx(ctx context.Context, tx *gorm.DB, userID string, otpHash string, expiredAt time.Time) error {
	userOTP := &entities.UserOTP{
		UserID:    userID,
		OTPHash:   otpHash,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
	}
	return tx.WithContext(ctx).Create(userOTP).Error
}

func (u *userRepository) GetActiveOTPByUserIDTx(ctx context.Context, tx *gorm.DB, userID string) (*entities.UserOTP, error) {
	var otp entities.UserOTP

	err := tx.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("used_at IS NULL").
		Where("expired_at > ?", time.Now()).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (u *userRepository) MarkOTPUsedTx(ctx context.Context, tx *gorm.DB, otpID string, usedAt time.Time) error {
	return tx.WithContext(ctx).Model(&entities.UserOTP{}).
		Where("id = ?", otpID).
		Update("used_at", usedAt).Error
}

func (u *userRepository) ActivateUserTx(ctx context.Context, tx *gorm.DB, userID string) error {
	return tx.WithContext(ctx).Model(&entities.User{}).
		Where("id = ?", userID).
		Update("is_active", true).Error
}
