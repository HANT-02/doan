package scheduling

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Preview(c *gin.Context)
	Benchmark(c *gin.Context)
	GetPreview(c *gin.Context)
	GetLatestPreview(c *gin.Context)
	Commit(c *gin.Context)
	SuggestSubstitute(c *gin.Context)
	AssignSubstitute(c *gin.Context)
	FindMakeupSpots(c *gin.Context)
}

func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, configManager config.Manager) {
	v1 := router.Group("/v1/scheduling")

	authMiddleware := middleware.AuthMiddleware(configManager)
	adminRoutes := v1.Group("")
	adminRoutes.Use(authMiddleware, middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"))
	adminRoutes.POST("/preview", ctrl.Preview)
	adminRoutes.POST("/benchmark", ctrl.Benchmark)
	adminRoutes.GET("/preview/latest", ctrl.GetLatestPreview)
	adminRoutes.GET("/preview/:id", ctrl.GetPreview)
	adminRoutes.POST("/commit", ctrl.Commit)

	substituteRoutes := v1.Group("")
	substituteRoutes.Use(authMiddleware, middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN", "TEACHER"))
	substituteRoutes.GET("/lessons/:id/suggest-substitutes", ctrl.SuggestSubstitute)
	substituteRoutes.POST("/lessons/:id/assign-substitute", ctrl.AssignSubstitute)

	makeupRoutes := v1.Group("")
	makeupRoutes.Use(authMiddleware, middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN"))
	makeupRoutes.GET("/lessons/:id/find-makeup-spots", ctrl.FindMakeupSpots)
}
