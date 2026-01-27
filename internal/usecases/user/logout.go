package user

import (
	"context"
	"doan/internal/services/user"
	"doan/pkg/logger"
)

// LogoutInput represents the input of the LogoutUseCase
type LogoutInput struct {
	Token string `json:"token"`
}

// LogoutOutput represents the output of the LogoutUseCase
type LogoutOutput struct {
	Message string `json:"message"`
}

// LogoutUseCase is a use case for logout
type LogoutUseCase interface {
	// Execute is a function to logout
	// Param: LogoutInput
	// Return: LogoutOutput, error
	Execute(ctx context.Context, input LogoutInput) (*LogoutOutput, error)
}

// logoutUseCase implements LogoutUseCase
type logoutUseCase struct {
	authService user.AuthService
}

func NewLogoutUseCase(authService user.AuthService) LogoutUseCase {
	return &logoutUseCase{authService: authService}
}

func (u *logoutUseCase) Execute(ctx context.Context, input LogoutInput) (*LogoutOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	// Validate token to ensure it's valid
	_, err := u.authService.ValidateToken(ctx, input.Token)
	if err != nil {
		ctxLogger.Errorf("Failed to validate token: %v", err)
		return nil, err
	}

	// In a production system, you would:
	// 1. Add the token to a blacklist/revocation list (Redis)
	// 2. Or invalidate the session in your session store
	// For now, we just return success

	ctxLogger.Info("User logged out successfully")

	return &LogoutOutput{
		Message: "Logged out successfully",
	}, nil
}
