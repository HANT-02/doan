package shift

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateShift(c *gin.Context)
	GetShift(c *gin.Context)
	UpdateShift(c *gin.Context)
	DeleteShift(c *gin.Context)
	ListShifts(c *gin.Context)
}

func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, configManager config.Manager) {
	v1 := router.Group("/v1/shifts")

	authMiddleware := middleware.AuthMiddleware(configManager)
	adminRole := middleware.RoleMiddleware("ADMIN")

	v1.Use(authMiddleware, adminRole)
	v1.GET("", ctrl.ListShifts)
	v1.GET("/:id", ctrl.GetShift)
	v1.POST("", ctrl.CreateShift)
	v1.PUT("/:id", ctrl.UpdateShift)
	v1.DELETE("/:id", ctrl.DeleteShift)
}
