package class

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateClass(c *gin.Context)
	GetClass(c *gin.Context)
	UpdateClass(c *gin.Context)
	DeleteClass(c *gin.Context)
	ListClasses(c *gin.Context)
}

// RegisterRoutesV1 registers class routes with the router
func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, configManager config.Manager) {
	v1 := router.Group("/v1/classes")

	// Middleware
	authMiddleware := middleware.AuthMiddleware(configManager)
	adminRole := middleware.RoleMiddleware("ADMIN")

	// Routes
	v1.GET("", ctrl.ListClasses)
	v1.GET("/:id", ctrl.GetClass)

	// Admin-only operations
	v1.POST("", authMiddleware, adminRole, ctrl.CreateClass)
	v1.PUT("/:id", authMiddleware, adminRole, ctrl.UpdateClass)
	v1.DELETE("/:id", authMiddleware, adminRole, ctrl.DeleteClass)
}
