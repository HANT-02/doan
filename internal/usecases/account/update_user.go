package account

import (
	"context"
	"errors"
	"strings"

	repositoryinterface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type UpdateUserInput struct {
	ID       string
	FullName string
	Role     string
	IsActive *bool
}

type UpdateUserUseCase interface {
	Execute(ctx context.Context, input UpdateUserInput) error
}

type updateUserUseCase struct {
	userRepo repositoryinterface.UserRepository
}

func NewUpdateUserUseCase(userRepo repositoryinterface.UserRepository) UpdateUserUseCase {
	return &updateUserUseCase{userRepo: userRepo}
}

func (uc *updateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) error {
	ctxLogger := logger.NewLogger(ctx)

	updateData := make(map[string]interface{})
	if strings.TrimSpace(input.FullName) != "" {
		updateData["full_name"] = strings.TrimSpace(input.FullName)
	}
	if strings.TrimSpace(input.Role) != "" {
		updateData["role"] = strings.ToUpper(strings.TrimSpace(input.Role))
	}
	if input.IsActive != nil {
		updateData["is_active"] = *input.IsActive
	}
	if len(updateData) == 0 {
		return errors.New("no update field provided")
	}

	if err := uc.userRepo.Update(ctx, input.ID, updateData); err != nil {
		ctxLogger.Errorf("Failed to update user %s: %v", input.ID, err)
		return err
	}
	return nil
}
