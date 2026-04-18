package teacherportal

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetTeacherLessons(c *gin.Context)
	GetTeacherLessonAttendance(c *gin.Context)
	SubmitTeacherLessonAttendance(c *gin.Context)
	UpdateTeacherLessonAttendance(c *gin.Context)
	GetTeacherAttendanceSummary(c *gin.Context)
	GetTeacherLessonSummary(c *gin.Context)
	UpsertTeacherLessonSummary(c *gin.Context)
	GetTeacherLessonAcademicRecords(c *gin.Context)
	UpsertTeacherLessonAcademicRecord(c *gin.Context)
	FinalizeTeacherLessonAcademicRecords(c *gin.Context)
	GetTeacherStudentAcademicRecords(c *gin.Context)
	ListTeacherLeaveRequests(c *gin.Context)
	ApproveTeacherLeaveRequest(c *gin.Context)
	RejectTeacherLeaveRequest(c *gin.Context)
}

func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, manager config.Manager) {
	v1 := router.Group("/v1/teacher")

	authMiddleware := middleware.AuthMiddleware(manager)
	teacherRole := middleware.RoleMiddleware("TEACHER")

	v1.GET("/lessons", authMiddleware, teacherRole, ctrl.GetTeacherLessons)
	v1.GET("/lessons/:lesson_id/attendance", authMiddleware, teacherRole, ctrl.GetTeacherLessonAttendance)
	v1.POST("/lessons/:lesson_id/attendance", authMiddleware, teacherRole, ctrl.SubmitTeacherLessonAttendance)
	v1.PUT("/lessons/:lesson_id/attendance/:student_id", authMiddleware, teacherRole, ctrl.UpdateTeacherLessonAttendance)
	v1.GET("/lessons/:lesson_id/summary", authMiddleware, teacherRole, ctrl.GetTeacherLessonSummary)
	v1.PUT("/lessons/:lesson_id/summary", authMiddleware, teacherRole, ctrl.UpsertTeacherLessonSummary)
	v1.GET("/lessons/:lesson_id/records", authMiddleware, teacherRole, ctrl.GetTeacherLessonAcademicRecords)
	v1.PUT("/lessons/:lesson_id/records/:student_id", authMiddleware, teacherRole, ctrl.UpsertTeacherLessonAcademicRecord)
	v1.POST("/lessons/:lesson_id/records/finalize", authMiddleware, teacherRole, ctrl.FinalizeTeacherLessonAcademicRecords)
	v1.GET("/classes/:class_id/attendance-summary", authMiddleware, teacherRole, ctrl.GetTeacherAttendanceSummary)
	v1.GET("/classes/:class_id/students/:student_id/records", authMiddleware, teacherRole, ctrl.GetTeacherStudentAcademicRecords)
	v1.GET("/leave-requests", authMiddleware, teacherRole, ctrl.ListTeacherLeaveRequests)
	v1.POST("/leave-requests/:id/approve", authMiddleware, teacherRole, ctrl.ApproveTeacherLeaveRequest)
	v1.POST("/leave-requests/:id/reject", authMiddleware, teacherRole, ctrl.RejectTeacherLeaveRequest)
}
