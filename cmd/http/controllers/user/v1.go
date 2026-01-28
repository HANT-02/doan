package user

import (
	"doan/cmd/http/rest"
	"doan/internal/usecases/user"
	"doan/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	loginUseCase          user.LoginUseCase
	logoutUseCase         user.LogoutUseCase
	refreshTokenUseCase   user.RefreshTokenUseCase
	registerUseCase       user.RegisterUseCase
	forgotPasswordUseCase user.ForgotPasswordUseCase
	resetPasswordUseCase  user.ResetPasswordUseCase
	changePasswordUseCase user.ChangePasswordUseCase
}

func NewUserControllerV1(
	loginUseCase user.LoginUseCase,
	logoutUseCase user.LogoutUseCase,
	refreshTokenUseCase user.RefreshTokenUseCase,
	registerUseCase user.RegisterUseCase,
	forgotPasswordUseCase user.ForgotPasswordUseCase,
	resetPasswordUseCase user.ResetPasswordUseCase,
	changePasswordUseCase user.ChangePasswordUseCase,
) *ControllerV1 {
	return &ControllerV1{
		loginUseCase:          loginUseCase,
		logoutUseCase:         logoutUseCase,
		refreshTokenUseCase:   refreshTokenUseCase,
		registerUseCase:       registerUseCase,
		forgotPasswordUseCase: forgotPasswordUseCase,
		resetPasswordUseCase:  resetPasswordUseCase,
		changePasswordUseCase: changePasswordUseCase,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body LoginRequest true "Login credentials"
// @Success 200 {object} rest.BaseResponse{data=LoginResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 401 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/auth/login [post]
func (c *ControllerV1) Login(ctx *gin.Context) {
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
// @Summary User logout
// @Description Invalidate user session
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body LogoutRequest true "Logout request"
// @Success 200 {object} rest.BaseResponse{data=LogoutResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 401 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/auth/logout [post]
func (c *ControllerV1) Logout(ctx *gin.Context) {
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
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} rest.BaseResponse{data=RefreshTokenResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 401 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/auth/refresh [post]
func (c *ControllerV1) RefreshToken(ctx *gin.Context) {
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

// Register godoc
// @Summary User registration
// @Description Create a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body RegisterRequest true "Register request"
// @Success 201 {object} rest.BaseResponse{data=RegisterResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 409 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/auth/register [post]
func (c *ControllerV1) Register(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	out, err := c.registerUseCase.Execute(ctx, user.RegisterInput{
		Email:       req.Email,
		FullName:    req.FullName,
		PasswordEnc: req.PasswordEnc,
	})
	if err != nil {
		ctxLogger.Errorf("Failed to register: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Registration failed", err)
		return
	}

	resp := RegisterResponse{UserID: out.UserID}
	rest.ResponseSuccess(ctx, http.StatusCreated, "Registered successfully", resp)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send a password reset email if the account exists
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} rest.BaseResponse{data=MessageResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/auth/forgot-password [post]
func (c *ControllerV1) ForgotPassword(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	var req ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := c.forgotPasswordUseCase.Execute(ctx, user.ForgotPasswordInput{Email: req.Email}); err != nil {
		ctxLogger.Errorf("Failed to process forgot password: %v", err)
		rest.ResponseError(ctx, http.StatusInternalServerError, "Failed to process request", err)
		return
	}

	rest.ResponseSuccess(ctx, http.StatusOK, "If the email exists, a reset link has been sent", MessageResponse{Message: "If the email exists, a reset link has been sent"})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using the provided token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param payload body ResetPasswordRequest true "Reset password request"
// @Success 200 {object} rest.BaseResponse{data=MessageResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/auth/reset-password [post]
func (c *ControllerV1) ResetPassword(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	var req ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := c.resetPasswordUseCase.Execute(ctx, user.ResetPasswordInput{
		Token:          req.Token,
		NewPasswordEnc: req.NewPasswordEnc,
	}); err != nil {
		ctxLogger.Errorf("Failed to reset password: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to reset password", err)
		return
	}

	rest.ResponseSuccess(ctx, http.StatusOK, "Password reset successfully", MessageResponse{Message: "Password reset successfully"})
}

// ChangePassword godoc
// @Summary Change password
// @Description Change current user's password
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body ChangePasswordRequest true "Change password request"
// @Success 200 {object} rest.BaseResponse{data=MessageResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 401 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/auth/change-password [post]
func (c *ControllerV1) ChangePassword(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		rest.ResponseError(ctx, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	userID, _ := userIDVal.(string)

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := c.changePasswordUseCase.Execute(ctx, user.ChangePasswordInput{
		UserID:         userID,
		OldPasswordEnc: req.OldPasswordEnc,
		NewPasswordEnc: req.NewPasswordEnc,
	}); err != nil {
		ctxLogger.Errorf("Failed to change password: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to change password", err)
		return
	}

	rest.ResponseSuccess(ctx, http.StatusOK, "Password changed successfully", MessageResponse{Message: "Password changed successfully"})
}
