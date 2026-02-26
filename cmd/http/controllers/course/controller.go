package course

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

// Controller defines the interface for course HTTP handlers
type Controller interface {
	CreateCourse(ctx *gin.Context)
	GetCourse(ctx *gin.Context)
	UpdateCourse(ctx *gin.Context)
	DeleteCourse(ctx *gin.Context)
	ListCourses(ctx *gin.Context)
}

// RegisterRoutesV1 registers course routes with the router
func RegisterRoutesV1(router *gin.RouterGroup, controller Controller, configManager config.Manager) {
	v1 := router.Group("/v1/courses")

	// Middleware
	authMiddleware := middleware.AuthMiddleware(configManager)
	adminRole := middleware.RoleMiddleware("ADMIN")

	// Admin-only routes (CRUD operations)
	v1.POST("", authMiddleware, adminRole, controller.CreateCourse)
	v1.PUT("/:id", authMiddleware, adminRole, controller.UpdateCourse)
	v1.DELETE("/:id", authMiddleware, adminRole, controller.DeleteCourse)

	// Public/authenticated routes (read operations)
	v1.GET("", controller.ListCourses)
	v1.GET("/:id", controller.GetCourse)
}
