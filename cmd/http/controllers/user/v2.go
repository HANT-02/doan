package user

import (
	"doan/cmd/http/rest"
	"doan/internal/usecases/user"
	"doan/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var _ Controller = (*ControllerV2)(nil)

type ControllerV2 struct {
	loginUseCase        user.LoginUseCase
	logoutUseCase       user.LogoutUseCase
	refreshTokenUseCase user.RefreshTokenUseCase
}

func NewUserControllerV2(
	loginUseCase user.LoginUseCase,
	logoutUseCase user.LogoutUseCase,
	refreshTokenUseCase user.RefreshTokenUseCase,
) *ControllerV2 {
	return &ControllerV2{
		loginUseCase:        loginUseCase,
		logoutUseCase:       logoutUseCase,
		refreshTokenUseCase: refreshTokenUseCase,
	}
}

// Login godoc
// @Summary User login (v2)
// @Description Authenticate user and return JWT tokens
// @Tags Authentication v2
// @Accept json
// @Produce json
// @Param payload body LoginRequest true "Login credentials"
// @Success 200 {object} rest.BaseResponse{data=LoginResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 401 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v2/auth/login [post]
func (c *ControllerV2) Login(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Execute login use case
	output, err := c.loginUseCase.Execute(ctx, user.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to login: %v", err)
		rest.ResponseError(ctx, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	// Map to response
	response := LoginResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		User: UserResponse{
			ID:       output.User.ID,
			Code:     output.User.Code,
			FullName: output.User.FullName,
			Email:    output.User.Email,
			Role:     output.User.Role,
			IsActive: output.User.IsActive,
		},
	}

	rest.ResponseSuccess(ctx, http.StatusOK, "Login successful", response)
}

// Logout godoc
// @Summary User logout (v2)
// @Description Invalidate user session
// @Tags Authentication v2
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body LogoutRequest true "Logout request"
// @Success 200 {object} rest.BaseResponse{data=LogoutResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 401 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v2/auth/logout [post]
func (c *ControllerV2) Logout(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	var req LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Execute logout use case
	output, err := c.logoutUseCase.Execute(ctx, user.LogoutInput{
		Token: req.Token,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to logout: %v", err)
		rest.ResponseError(ctx, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// Map to response
	response := LogoutResponse{
		Message: output.Message,
	}

	rest.ResponseSuccess(ctx, http.StatusOK, "Logout successful", response)
}

// RefreshToken godoc
// @Summary Refresh access token (v2)
// @Description Get a new access token using refresh token
// @Tags Authentication v2
// @Accept json
// @Produce json
// @Param payload body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} rest.BaseResponse{data=RefreshTokenResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 401 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v2/auth/refresh [post]
func (c *ControllerV2) RefreshToken(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Execute refresh token use case
	output, err := c.refreshTokenUseCase.Execute(ctx, user.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to refresh token: %v", err)
		rest.ResponseError(ctx, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// Map to response
	response := RefreshTokenResponse{
		AccessToken: output.AccessToken,
	}

	rest.ResponseSuccess(ctx, http.StatusOK, "Token refreshed successfully", response)
}

func (c *ControllerV2) Register(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c *ControllerV2) ForgotPassword(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c *ControllerV2) ResetPassword(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c *ControllerV2) ChangePassword(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}
