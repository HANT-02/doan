package user

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/internal/services/security"
	"errors"
	"strings"
)

// Register

type RegisterInput struct {
	Email       string
	FullName    string
	PasswordEnc string
}

type RegisterOutput struct {
	UserID string
}

type RegisterUseCase interface {
	Execute(ctx context.Context, in RegisterInput) (*RegisterOutput, error)
}

type registerUseCase struct {
	userRepo repositoryinterface.UserRepository
	cipher   security.PasswordCipher
	hasher   security.PasswordHasher
}

func NewRegisterUseCase(repo repositoryinterface.UserRepository, cipher security.PasswordCipher, hasher security.PasswordHasher) RegisterUseCase {
	return &registerUseCase{userRepo: repo, cipher: cipher, hasher: hasher}
}

func (u *registerUseCase) Execute(ctx context.Context, in RegisterInput) (*RegisterOutput, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" || in.PasswordEnc == "" || strings.TrimSpace(in.FullName) == "" {
		return nil, errors.New("invalid payload")
	}
	// Check unique email
	cond := repositories.NewCommonCondition()
	cond.AddCondition("email", email, repositories.Equal)
	p, err := u.userRepo.GetByCondition(ctx, cond)
	if err == nil && p != nil && len(p.Data) > 0 {
		return nil, errors.New("email already registered")
	}

	plain, err := u.cipher.Decrypt(in.PasswordEnc)
	if err != nil {
		return nil, err
	}
	hash, err := u.hasher.Hash(plain)
	if err != nil {
		return nil, err
	}
	newUser := &entities.User{
		FullName: in.FullName,
		Email:    email,
		Password: hash,
	}
	created, err := u.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &RegisterOutput{UserID: created.ID}, nil
}
