package user

import (
	"context"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/internal/services/security"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ConfirmChangePasswordOTPInput struct {
	UserID string
	OTP    string
}

type ConfirmChangePasswordOTPUseCase interface {
	Execute(ctx context.Context, in ConfirmChangePasswordOTPInput) error
}

type confirmChangePasswordOTPUseCase struct {
	db       *gorm.DB
	userRepo repositoryinterface.UserRepository
	hasher   security.PasswordHasher
}

func NewConfirmChangePasswordOTPUseCase(
	db *gorm.DB,
	userRepo repositoryinterface.UserRepository,
	hasher security.PasswordHasher,
) ConfirmChangePasswordOTPUseCase {
	return &confirmChangePasswordOTPUseCase{
		db:       db,
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (u *confirmChangePasswordOTPUseCase) Execute(ctx context.Context, in ConfirmChangePasswordOTPInput) error {
	if strings.TrimSpace(in.UserID) == "" || strings.TrimSpace(in.OTP) == "" {
		return errors.New("invalid payload")
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		otpEntity, err := u.userRepo.GetActiveOTPByUserIDAndPurposeTx(ctx, tx, in.UserID, OTPPurposeChangePassword)
		if err != nil {
			return errors.New("otp not found or expired")
		}
		if time.Now().After(otpEntity.ExpiredAt) {
			return errors.New("otp expired")
		}
		if hashErr := u.hasher.Compare(otpEntity.OTPHash, in.OTP); hashErr != nil {
			return errors.New("invalid otp")
		}
		if otpEntity.PendingPasswordHash == nil || strings.TrimSpace(*otpEntity.PendingPasswordHash) == "" {
			return errors.New("pending password not found")
		}

		if err := u.userRepo.UpdatePasswordTx(ctx, tx, in.UserID, *otpEntity.PendingPasswordHash); err != nil {
			return err
		}

		now := time.Now()
		return u.userRepo.MarkOTPUsedTx(ctx, tx, otpEntity.ID, now)
	})
}
