package account

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"doan/internal/entities"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/internal/services/security"
	"doan/pkg/logger"
)

type CreateUserInput struct {
	Code     string
	FullName string
	Email    string
	Role     string
	IsActive bool
	Password string
}

type CreateUserOutput struct {
	User entities.User `json:"user"`
}

type CreateUserUseCase interface {
	Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error)
}

type createUserUseCase struct {
	userRepo repositoryinterface.UserRepository
	hasher   security.PasswordHasher
}

func NewCreateUserUseCase(userRepo repositoryinterface.UserRepository, hasher security.PasswordHasher) CreateUserUseCase {
	return &createUserUseCase{userRepo: userRepo, hasher: hasher}
}

func (uc *createUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || strings.TrimSpace(input.Password) == "" || strings.TrimSpace(input.FullName) == "" {
		return nil, errors.New("full_name, email and password are required")
	}

	exists, err := uc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		ctxLogger.Errorf("Failed to check email %s: %v", email, err)
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	passwordHash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	code := strings.TrimSpace(input.Code)
	if code == "" {
		code = fmt.Sprintf("USR-%s", time.Now().Format("20060102150405"))
	}

	role := strings.ToUpper(strings.TrimSpace(input.Role))
	if role == "" {
		role = "STUDENT"
	}

	newUser := &entities.User{
		Code:      code,
		FullName:  strings.TrimSpace(input.FullName),
		Email:     email,
		Role:      role,
		IsActive:  input.IsActive,
		Password:  passwordHash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := uc.userRepo.Create(ctx, newUser)
	if err != nil {
		ctxLogger.Errorf("Failed to create user %s: %v", email, err)
		return nil, err
	}

	return &CreateUserOutput{User: *created}, nil
}
