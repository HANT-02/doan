package user

import (
	"context"
	_interface "doan/internal/repositories/interface"
	"doan/internal/services/user"
	"doan/pkg/logger"
)

// LoginInput represents the input of the LoginUseCase
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginOutput represents the output of the LoginUseCase
type LoginOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// LoginUseCase is a use case for login
type LoginUseCase interface {
	// Execute is a function to login
	// Param: LoginInput
	// Return: LoginOutput, error
	Execute(ctx context.Context, input LoginInput) (*LoginOutput, error)
}

// loginUseCase implements LoginUseCase
// Dependencies: AuthService, UserRepository
type loginUseCase struct {
	authService user.AuthService
	userRepo    _interface.UserRepository
}

func NewLoginUseCase(authService user.AuthService, userRepo _interface.UserRepository) LoginUseCase {
	return &loginUseCase{authService: authService, userRepo: userRepo}
}

func (u *loginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	// Generate JWT Token
	token, err := u.authService.CreateAuthToken(ctx, user.CreateAuthTokenInput{
		Username: input.Username,
		Password: input.Password,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to create auth token: %v", err)
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
