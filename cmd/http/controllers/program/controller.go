package program

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

// Controller defines the interface for program HTTP handlers
type Controller interface {
	CreateProgram(ctx *gin.Context)
	GetProgram(ctx *gin.Context)
	UpdateProgram(ctx *gin.Context)
	DeleteProgram(ctx *gin.Context)
	ListPrograms(ctx *gin.Context)
	AddCourses(ctx *gin.Context)
	RemoveCourses(ctx *gin.Context)
}

// RegisterRoutesV1 registers program routes with the router
func RegisterRoutesV1(router *gin.RouterGroup, controller Controller, configManager config.Manager) {
	v1 := router.Group("/v1/programs")

	// Middleware
	authMiddleware := middleware.AuthMiddleware(configManager)
	adminRole := middleware.RoleMiddleware("ADMIN")

	// Admin-only routes (CRUD operations)
	v1.POST("", authMiddleware, adminRole, controller.CreateProgram)
	v1.PUT("/:id", authMiddleware, adminRole, controller.UpdateProgram)
	v1.DELETE("/:id", authMiddleware, adminRole, controller.DeleteProgram)

	v1.POST("/:id/courses", authMiddleware, adminRole, controller.AddCourses)
	v1.DELETE("/:id/courses", authMiddleware, adminRole, controller.RemoveCourses)

	// Public/authenticated routes (read operations)
	v1.GET("", controller.ListPrograms)
	v1.GET("/:id", controller.GetProgram)
}
