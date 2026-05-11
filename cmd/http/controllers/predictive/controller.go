package predictive

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	ListAtRiskStudents(c *gin.Context)
	GetModelMetadata(c *gin.Context)
	TrainAtRiskFromDB(c *gin.Context)
}

func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, manager config.Manager) {
	predictiveRoutes := router.Group("/v1/predictive")
	predictiveRoutes.Use(middleware.AuthMiddleware(manager))
	{
		predictiveRoutes.GET("/at-risk/students", ctrl.ListAtRiskStudents)
		predictiveRoutes.GET("/at-risk/model-metadata", ctrl.GetModelMetadata)
		predictiveRoutes.POST("/at-risk/train-from-db", ctrl.TrainAtRiskFromDB)
	}
}
