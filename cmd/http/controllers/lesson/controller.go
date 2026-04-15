package lesson

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	ListLessons(c *gin.Context)
	GetLesson(c *gin.Context)
	GetLessonAttendance(c *gin.Context)
	UpsertLessonAttendance(c *gin.Context)
	GetLessonSummary(c *gin.Context)
	UpsertLessonSummary(c *gin.Context)
	GetLessonAcademicRecords(c *gin.Context)
	UpsertLessonAcademicRecords(c *gin.Context)
	FinalizeLessonAcademicRecords(c *gin.Context)
}

// RegisterRoutesV1 registers lesson routes with the router
func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, configManager config.Manager) {
	v1 := router.Group("/v1/lessons")

	authMiddleware := middleware.AuthMiddleware(configManager)
	lessonOpsRole := middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN", "TEACHER")

	v1.GET("", authMiddleware, ctrl.ListLessons)
	v1.GET("/:id", authMiddleware, ctrl.GetLesson)
	v1.GET("/:id/attendance", authMiddleware, lessonOpsRole, ctrl.GetLessonAttendance)
	v1.PUT("/:id/attendance", authMiddleware, lessonOpsRole, ctrl.UpsertLessonAttendance)
	v1.GET("/:id/summary", authMiddleware, lessonOpsRole, ctrl.GetLessonSummary)
	v1.PUT("/:id/summary", authMiddleware, lessonOpsRole, ctrl.UpsertLessonSummary)
	v1.GET("/:id/records", authMiddleware, lessonOpsRole, ctrl.GetLessonAcademicRecords)
	v1.PUT("/:id/records", authMiddleware, lessonOpsRole, ctrl.UpsertLessonAcademicRecords)
	v1.POST("/:id/records/finalize", authMiddleware, lessonOpsRole, ctrl.FinalizeLessonAcademicRecords)
}
