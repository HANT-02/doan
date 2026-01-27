package user

import (
	"context"
	"doan/internal/services/user"
	"doan/pkg/logger"
)

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenOutput struct {
	AccessToken string `json:"access_token"`
}

type RefreshTokenUseCase interface {
	Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error)
}

type refreshTokenUseCase struct {
	authService user.AuthService
}

func NewRefreshTokenUseCase(authService user.AuthService) RefreshTokenUseCase {
	return &refreshTokenUseCase{authService: authService}
}

func (u *refreshTokenUseCase) Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	// Refresh access token
	accessToken, err := u.authService.RefreshAccessToken(ctx, input.RefreshToken)
	if err != nil {
		ctxLogger.Errorf("Failed to refresh token: %v", err)
		return nil, err
	}

	ctxLogger.Info("Token refreshed successfully")

	return &RefreshTokenOutput{
		AccessToken: accessToken,
	}, nil
}
