package user

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	Register(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	VerifyOTP(ctx *gin.Context)
	GetMe(ctx *gin.Context)
}

// RegisterRoutesV1 register routes for version 1
func RegisterRoutesV1(router *gin.RouterGroup, controller Controller) {
	v1 := router.Group("/v1/auth")
	{
		v1.POST("/login", controller.Login)
		v1.POST("/logout", controller.Logout)
		v1.POST("/refresh", controller.RefreshToken)
		v1.POST("/register", controller.Register)
		v1.POST("/forgot-password", controller.ForgotPassword)
		v1.POST("/reset-password", controller.ResetPassword)
		v1.POST("/change-password", controller.ChangePassword)
		v1.POST("/verify-otp", controller.VerifyOTP)
		v1.GET("/me", controller.GetMe)
	}
}

// RegisterRoutesV2 register routes for version 2
func RegisterRoutesV2(router *gin.RouterGroup, controller Controller) {
	v2 := router.Group("/v2/auth")
	{
		v2.POST("/login", controller.Login)
		v2.POST("/logout", controller.Logout)
		v2.POST("/refresh", controller.RefreshToken)
	}
}
