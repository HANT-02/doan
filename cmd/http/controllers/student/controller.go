package student

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateStudent(c *gin.Context)
	GetStudent(c *gin.Context)
	UpdateStudent(c *gin.Context)
	DeleteStudent(c *gin.Context)
	ListStudents(c *gin.Context)
}

func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, manager config.Manager) {
	studentRoutes := router.Group("/v1/students")
	studentRoutes.Use(middleware.AuthMiddleware(manager))
	{
		studentRoutes.POST("", ctrl.CreateStudent)
		studentRoutes.GET("", ctrl.ListStudents)
		studentRoutes.GET("/:id", ctrl.GetStudent)
		studentRoutes.PUT("/:id", ctrl.UpdateStudent)
		studentRoutes.DELETE("/:id", ctrl.DeleteStudent)
	}
}
