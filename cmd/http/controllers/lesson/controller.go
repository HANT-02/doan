package lesson

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	ListLessons(c *gin.Context)
	GetLesson(c *gin.Context)
}

// RegisterRoutesV1 registers lesson routes with the router
func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, configManager config.Manager) {
	v1 := router.Group("/v1/lessons")

	authMiddleware := middleware.AuthMiddleware(configManager)

	v1.GET("", authMiddleware, ctrl.ListLessons)
	v1.GET("/:id", authMiddleware, ctrl.GetLesson)
}
