package user

import "github.com/gin-gonic/gin"

type Controller interface {
	Login(ctx *gin.Context)
}

// RegisterRoutesV1 register routes for version 1
func RegisterRoutesV1(router *gin.Engine, controller Controller) {
	v1 := router.Group("/v1")
	v1.POST("/login", controller.Login)
}

// RegisterRoutesV2 register routes for version 2
func RegisterRoutesV2(router *gin.Engine, controller Controller) {
	v2 := router.Group("/v2")
	v2.POST("/login", controller.Login)
}
