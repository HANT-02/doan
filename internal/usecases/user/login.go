package user

import (
	"context"
	"doan/internal/services/user"
	"doan/pkg/logger"
)

// LoginInput represents the input of the LoginUseCase
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUseCaseUserOutput struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

// LoginOutput represents the output of the LoginUseCase
type LoginOutput struct {
	AccessToken  string                 `json:"access_token"`
	RefreshToken string                 `json:"refresh_token"`
	User         LoginUseCaseUserOutput `json:"user"`
}

// LoginUseCase is a use case for login
type LoginUseCase interface {
	// Execute is a function to login
	// Param: LoginInput
	// Return: LoginOutput, error
	Execute(ctx context.Context, input LoginInput) (*LoginOutput, error)
}

// loginUseCase implements LoginUseCase
type loginUseCase struct {
	authService user.AuthService
}

func NewLoginUseCase(authService user.AuthService) LoginUseCase {
	return &loginUseCase{authService: authService}
}

func (u *loginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	token, err := u.authService.CreateAuthToken(ctx, user.CreateAuthTokenInput{
		Username: input.Username,
		Password: input.Password,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to create auth token: %v", err)
		return nil, err
	}

	loginUser := LoginUseCaseUserOutput{
		ID:       token.User.ID,
		Code:     token.User.Code,
		FullName: token.User.FullName,
		Email:    token.User.Email,
		Role:     token.User.Role,
		IsActive: token.User.IsActive,
	}

	return &LoginOutput{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         loginUser,
	}, nil
}
