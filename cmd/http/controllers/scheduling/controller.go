package scheduling

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Preview(c *gin.Context)
	GetPreview(c *gin.Context)
	GetLatestPreview(c *gin.Context)
	Commit(c *gin.Context)
}

func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, configManager config.Manager) {
	v1 := router.Group("/v1/scheduling")

	authMiddleware := middleware.AuthMiddleware(configManager)
	adminRole := middleware.RoleMiddleware("ADMIN")

	v1.Use(authMiddleware, adminRole)
	v1.POST("/preview", ctrl.Preview)
	v1.GET("/preview/latest", ctrl.GetLatestPreview)
	v1.GET("/preview/:id", ctrl.GetPreview)
	v1.POST("/commit", ctrl.Commit)
}
