package room

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateRoom(c *gin.Context)
	GetRoom(c *gin.Context)
	UpdateRoom(c *gin.Context)
	DeleteRoom(c *gin.Context)
	ListRooms(c *gin.Context)
}

// RegisterRoutesV1 registers room routes with the router
func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, configManager config.Manager) {
	v1 := router.Group("/v1/rooms")

	// Middleware
	authMiddleware := middleware.AuthMiddleware(configManager)
	adminRole := middleware.RoleMiddleware("ADMIN")

	// Routes
	v1.GET("", ctrl.ListRooms)
	v1.GET("/:id", ctrl.GetRoom)

	// Admin-only operations
	v1.POST("", authMiddleware, adminRole, ctrl.CreateRoom)
	v1.PUT("/:id", authMiddleware, adminRole, ctrl.UpdateRoom)
	v1.DELETE("/:id", authMiddleware, adminRole, ctrl.DeleteRoom)
}
