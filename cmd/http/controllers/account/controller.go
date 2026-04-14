package account

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	ListUsers(c *gin.Context)
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	ResetPassword(c *gin.Context)
}

func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, manager config.Manager) {
	accountRoutes := router.Group("/v1/users")
	accountRoutes.Use(middleware.AuthMiddleware(manager), middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"))
	{
		accountRoutes.GET("", ctrl.ListUsers)
		accountRoutes.GET("/:id", ctrl.GetUser)
		accountRoutes.POST("", ctrl.CreateUser)
		accountRoutes.PUT("/:id", ctrl.UpdateUser)
		accountRoutes.POST("/:id/reset-password", ctrl.ResetPassword)
	}
}
