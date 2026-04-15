package academic

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	ListMyAcademicRecords(c *gin.Context)
}

func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, manager config.Manager) {
	v1 := router.Group("/v1/academic-records")
	authMiddleware := middleware.AuthMiddleware(manager)
	studentRole := middleware.RoleMiddleware("STUDENT", "PARENT")

	v1.GET("/my", authMiddleware, studentRole, ctrl.ListMyAcademicRecords)
}
