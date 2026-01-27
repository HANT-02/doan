package user

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
	_interface "doan/internal/repositories/interface"
	"doan/pkg/config"
	"doan/pkg/logger"
	"doan/pkg/types"
	"doan/pkg/utils"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	CreateAuthToken(ctx context.Context, input CreateAuthTokenInput) (*CreateAuthTokenOutput, error)
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type authService struct {
	userRepo      _interface.UserRepository
	configManager config.Manager
	log           logger.Logger
}

func NewAuthService(
	userRepo _interface.UserRepository,
	configManager config.Manager,
	log logger.Logger,
) AuthService {
	return &authService{
		userRepo:      userRepo,
		configManager: configManager,
		log:           log,
	}
}

func (s *authService) CreateAuthToken(ctx context.Context, input CreateAuthTokenInput) (*CreateAuthTokenOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	// Get user by email
	userCondition := repositories.NewCommonCondition()
	userCondition.AddCondition("email", input.Username, repositories.Equal)
	userPagination, err := s.userRepo.GetByCondition(ctx, userCondition)
	if err != nil {
		ctxLogger.Errorf("failed to get user: %v", err)
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if userPagination == nil || len(userPagination.Data) == 0 {
		ctxLogger.Info("user not found")
		return nil, errors.New("user not found")
	}

	user := userPagination.Data[0]

	// Check if user is active
	if !user.IsActive {
		ctxLogger.Info("user is not active")
		return nil, errors.New("user is not active")
	}

	// Verify password
	if compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); compareErr != nil {
		ctxLogger.Errorf("password does not match: %v", compareErr)
		return nil, errors.New("invalid credentials")
	}

	// Get JWT config
	jwtConfig := types.JWTConfig{}
	if err := s.configManager.UnmarshalKey("jwt", &jwtConfig); err != nil {
		ctxLogger.Errorf("failed to get jwt config: %v", err)
		return nil, fmt.Errorf("failed to get jwt config: %w", err)
	}

	// Parse token durations
	accessTokenDuration, err := time.ParseDuration(jwtConfig.AccessTokenDuration)
	if err != nil {
		accessTokenDuration = 24 * time.Hour // Default to 24 hours
	}

	refreshTokenDuration, err := time.ParseDuration(jwtConfig.RefreshTokenDuration)
	if err != nil {
		refreshTokenDuration = 168 * time.Hour // Default to 7 days
	}

	// Generate access token (JWT)
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role, jwtConfig.Secret, accessTokenDuration)
	if err != nil {
		ctxLogger.Errorf("failed to generate access token: %v", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateToken(user.ID, user.Email, user.Role, jwtConfig.Secret, refreshTokenDuration)
	if err != nil {
		ctxLogger.Errorf("failed to generate refresh token: %v", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	ctxLogger.Infof("user %s logged in successfully", user.Email)

	return &CreateAuthTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         s.mapUserToOutput(user),
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	ctxLogger := logger.NewLogger(ctx)

	// Get JWT config
	jwtConfig := types.JWTConfig{}
	if err := s.configManager.UnmarshalKey("jwt", &jwtConfig); err != nil {
		ctxLogger.Errorf("failed to get jwt config: %v", err)
		return nil, fmt.Errorf("failed to get jwt config: %w", err)
	}

	// Validate token
	claims, err := utils.ValidateToken(token, jwtConfig.Secret)
	if err != nil {
		ctxLogger.Errorf("failed to validate token: %v", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return &TokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	ctxLogger := logger.NewLogger(ctx)

	// Get JWT config
	jwtConfig := types.JWTConfig{}
	if err := s.configManager.UnmarshalKey("jwt", &jwtConfig); err != nil {
		ctxLogger.Errorf("failed to get jwt config: %v", err)
		return "", fmt.Errorf("failed to get jwt config: %w", err)
	}

	// Parse access token duration
	accessTokenDuration, err := time.ParseDuration(jwtConfig.AccessTokenDuration)
	if err != nil {
		accessTokenDuration = 24 * time.Hour // Default to 24 hours
	}

	// Refresh token
	newAccessToken, err := utils.RefreshToken(refreshToken, jwtConfig.Secret, accessTokenDuration)
	if err != nil {
		ctxLogger.Errorf("failed to refresh token: %v", err)
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	return newAccessToken, nil
}

func (s *authService) mapUserToOutput(user *entities.User) UserOutput {
	return UserOutput{
		ID:       user.ID,
		Code:     user.Code,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}
}
