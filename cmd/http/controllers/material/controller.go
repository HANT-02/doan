package material

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Upload(c *gin.Context)
	List(c *gin.Context)
	ListFlagged(c *gin.Context)
	Get(c *gin.Context)
	Download(c *gin.Context)
	Review(c *gin.Context)
}

func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, configManager config.Manager) {
	v1 := router.Group("/v1/materials")

	authMiddleware := middleware.AuthMiddleware(configManager)
	teacherRole := middleware.RoleMiddleware("TEACHER")

	v1.GET("", authMiddleware, ctrl.List)
	v1.GET("/flagged", authMiddleware, ctrl.ListFlagged)
	v1.GET("/:id", authMiddleware, ctrl.Get)
	v1.GET("/:id/download", authMiddleware, ctrl.Download)
	v1.POST("/upload", authMiddleware, teacherRole, ctrl.Upload)
	v1.POST("/:id/review", authMiddleware, ctrl.Review)
}
