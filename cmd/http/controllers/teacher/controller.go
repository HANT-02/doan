package teacher

import "github.com/gin-gonic/gin"

// Controller defines the interface for teacher HTTP handlers
type Controller interface {
	CreateTeacher(ctx *gin.Context)
	GetTeacher(ctx *gin.Context)
	UpdateTeacher(ctx *gin.Context)
	DeleteTeacher(ctx *gin.Context)
	ListTeachers(ctx *gin.Context)
	GetTeacherTimetable(ctx *gin.Context)
	GetTeachingHoursStats(ctx *gin.Context)
}
