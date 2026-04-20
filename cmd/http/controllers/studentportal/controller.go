package studentportal

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetStudentTimetable(c *gin.Context)
	GetMyAttendance(c *gin.Context)
	GetMyAcademicRecords(c *gin.Context)
}

func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, manager config.Manager) {
	v1 := router.Group("/v1/student")

	authMiddleware := middleware.AuthMiddleware(manager)
	studentRole := middleware.RoleMiddleware("STUDENT")

	v1.GET("/timetable", authMiddleware, studentRole, ctrl.GetStudentTimetable)
	v1.GET("/attendance", authMiddleware, studentRole, ctrl.GetMyAttendance)
	v1.GET("/academic-records", authMiddleware, studentRole, ctrl.GetMyAcademicRecords)
}
