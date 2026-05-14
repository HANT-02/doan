package class

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateClass(c *gin.Context)
	GetClass(c *gin.Context)
	GetClassRoster(c *gin.Context)
	UpdateClass(c *gin.Context)
	DeleteClass(c *gin.Context)
	ListClasses(c *gin.Context)
	EnrollStudents(c *gin.Context)
	RemoveStudents(c *gin.Context)
	ReserveStudent(c *gin.Context)
	ResumeStudent(c *gin.Context)
	TransferStudent(c *gin.Context)
	AssignTeacher(c *gin.Context)
	GetClassSchedules(c *gin.Context)
	CreateClassSchedule(c *gin.Context)
	DeleteClassSchedule(c *gin.Context)
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
	v1.GET("/:id/students", ctrl.GetClassRoster)

	// Admin-only operations
	v1.POST("", authMiddleware, adminRole, ctrl.CreateClass)
	v1.PUT("/:id", authMiddleware, adminRole, ctrl.UpdateClass)
	v1.DELETE("/:id", authMiddleware, adminRole, ctrl.DeleteClass)
	v1.POST("/:id/students", authMiddleware, adminRole, ctrl.EnrollStudents)
	v1.DELETE("/:id/students", authMiddleware, adminRole, ctrl.RemoveStudents)
	v1.POST("/:id/students/:studentId/reserve", authMiddleware, adminRole, ctrl.ReserveStudent)
	v1.POST("/:id/students/:studentId/resume", authMiddleware, adminRole, ctrl.ResumeStudent)
	v1.POST("/:id/students/:studentId/transfer", authMiddleware, adminRole, ctrl.TransferStudent)
	v1.PUT("/:id/teacher", authMiddleware, adminRole, ctrl.AssignTeacher)

	// Class Schedule operations
	v1.GET("/:id/schedules", authMiddleware, ctrl.GetClassSchedules)
	v1.POST("/:id/schedules", authMiddleware, adminRole, ctrl.CreateClassSchedule)
	v1.DELETE("/:id/schedules/:scheduleId", authMiddleware, adminRole, ctrl.DeleteClassSchedule)
}
