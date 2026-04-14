package account

import (
	"context"
	"fmt"

	"doan/internal/entities"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetUserOutput struct {
	User entities.User `json:"user"`
}

type GetUserUseCase interface {
	Execute(ctx context.Context, id string) (*GetUserOutput, error)
}

type getUserUseCase struct {
	userRepo repositoryinterface.UserRepository
}

func NewGetUserUseCase(userRepo repositoryinterface.UserRepository) GetUserUseCase {
	return &getUserUseCase{userRepo: userRepo}
}

func (uc *getUserUseCase) Execute(ctx context.Context, id string) (*GetUserOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		ctxLogger.Errorf("Failed to get user %s: %v", id, err)
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user %s not found", id)
	}
	return &GetUserOutput{User: *user}, nil
}
