package teacher

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

// RegisterRoutesV1 registers teacher routes with the router
func RegisterRoutesV1(router *gin.RouterGroup, controller Controller, configManager config.Manager) {
	v1 := router.Group("/v1/teachers")

	// Middleware
	authMiddleware := middleware.AuthMiddleware(configManager)
	adminRole := middleware.RoleMiddleware("ADMIN")

	// Admin-only routes (CRUD operations)
	v1.POST("", authMiddleware, adminRole, controller.CreateTeacher)
	v1.PUT("/:id", authMiddleware, adminRole, controller.UpdateTeacher)
	v1.DELETE("/:id", authMiddleware, adminRole, controller.DeleteTeacher)

	// Public/authenticated routes (read operations)
	v1.GET("", controller.ListTeachers)
	v1.GET("/:id", controller.GetTeacher)
	v1.GET("/:id/timetable", controller.GetTeacherTimetable)
	v1.GET("/:id/stats/teaching-hours", controller.GetTeachingHoursStats)
}
