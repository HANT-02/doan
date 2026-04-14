package account

import (
	"context"
	"errors"
	"strings"

	repositoryinterface "doan/internal/repositories/interface"
	"doan/internal/services/security"
	"doan/pkg/logger"
)

type ResetUserPasswordInput struct {
	ID          string
	NewPassword string
}

type ResetUserPasswordUseCase interface {
	Execute(ctx context.Context, input ResetUserPasswordInput) error
}

type resetUserPasswordUseCase struct {
	userRepo repositoryinterface.UserRepository
	hasher   security.PasswordHasher
}

func NewResetUserPasswordUseCase(userRepo repositoryinterface.UserRepository, hasher security.PasswordHasher) ResetUserPasswordUseCase {
	return &resetUserPasswordUseCase{userRepo: userRepo, hasher: hasher}
}

func (uc *resetUserPasswordUseCase) Execute(ctx context.Context, input ResetUserPasswordInput) error {
	ctxLogger := logger.NewLogger(ctx)

	if strings.TrimSpace(input.NewPassword) == "" {
		return errors.New("new password is required")
	}

	hash, err := uc.hasher.Hash(input.NewPassword)
	if err != nil {
		return err
	}

	if err := uc.userRepo.Update(ctx, input.ID, map[string]interface{}{"password": hash}); err != nil {
		ctxLogger.Errorf("Failed to reset password for user %s: %v", input.ID, err)
		return err
	}
	return nil
}
