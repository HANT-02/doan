package user

import (
	"context"
	_interface "doan/internal/repositories/interface"
)

type AuthService interface {
	CreateAuthToken(ctx context.Context, input CreateAuthTokenInput) (*CreateAuthTokenOutput, error)
}

type authService struct {
	userRepo _interface.UserRepository
}

func NewAuthService(userRepo _interface.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) CreateAuthToken(ctx context.Context, input CreateAuthTokenInput) (*CreateAuthTokenOutput, error) {
	// Business logic here
	return &CreateAuthTokenOutput{
		AccessToken:  "hehe",
		RefreshToken: "hihi",
	}, nil
}
