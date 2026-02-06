package user

import (
	"context"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/internal/services/security"
	"errors"
	"gorm.io/gorm"
	"time"
)

type VerifyOTPInput struct {
	UserID string
	OTP    string
}

type VerifyOTPUseCase interface {
	Execute(ctx context.Context, in VerifyOTPInput) error
}

type verifyOTPUseCase struct {
	db       *gorm.DB
	userRepo repositoryinterface.UserRepository
	hasher   security.PasswordHasher
}

func NewVerifyOTPUseCase(
	db *gorm.DB,
	userRepo repositoryinterface.UserRepository,
	hasher security.PasswordHasher,
) VerifyOTPUseCase {
	return &verifyOTPUseCase{
		db:       db,
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (u *verifyOTPUseCase) Execute(ctx context.Context, in VerifyOTPInput) error {
	if in.OTP == "" {
		return errors.New("otp is required")
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		otpEntity, err := u.userRepo.GetActiveOTPByUserIDTx(ctx, tx, in.UserID)
		if err != nil {
			return errors.New("otp not found or expired")
		}

		// check expired
		if time.Now().After(otpEntity.ExpiredAt) {
			return errors.New("otp expired")
		}

		// compare otp
		if hashErr := u.hasher.Compare(otpEntity.OTPHash, in.OTP); hashErr != nil {
			return errors.New("invalid otp")
		}

		now := time.Now()

		// mark otp used
		if markOtpErr := u.userRepo.MarkOTPUsedTx(ctx, tx, otpEntity.ID, now); markOtpErr != nil {
			return markOtpErr
		}

		// activate user
		if userErr := u.userRepo.ActivateUserTx(ctx, tx, in.UserID); userErr != nil {
			return userErr
		}

		return nil
	})
}
