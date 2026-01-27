package user

import "github.com/gin-gonic/gin"

type Controller interface {
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

// RegisterRoutesV1 register routes for version 1
func RegisterRoutesV1(router *gin.Engine, controller Controller) {
	v1 := router.Group("/v1/auth")
	{
		v1.POST("/login", controller.Login)
		v1.POST("/logout", controller.Logout)
		v1.POST("/refresh", controller.RefreshToken)
	}
}

// RegisterRoutesV2 register routes for version 2
func RegisterRoutesV2(router *gin.Engine, controller Controller) {
	v2 := router.Group("/v2/auth")
	{
		v2.POST("/login", controller.Login)
		v2.POST("/logout", controller.Logout)
		v2.POST("/refresh", controller.RefreshToken)
	}
}
