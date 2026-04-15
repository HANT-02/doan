package leave

import (
	"doan/cmd/http/middleware"
	"doan/pkg/config"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	ListLeaveRequests(c *gin.Context)
	CreateLeaveRequest(c *gin.Context)
	ApproveLeaveRequest(c *gin.Context)
	RejectLeaveRequest(c *gin.Context)
	CancelLeaveRequest(c *gin.Context)
}

func RegisterRoutesV1(router *gin.RouterGroup, ctrl Controller, manager config.Manager) {
	v1 := router.Group("/v1/leave-requests")
	authMiddleware := middleware.AuthMiddleware(manager)
	staffRole := middleware.RoleMiddleware("ADMIN", "SUPER_ADMIN", "TEACHER")
	studentRole := middleware.RoleMiddleware("STUDENT", "PARENT", "ADMIN", "SUPER_ADMIN")

	v1.GET("", authMiddleware, ctrl.ListLeaveRequests)
	v1.POST("", authMiddleware, studentRole, ctrl.CreateLeaveRequest)
	v1.POST("/:id/approve", authMiddleware, staffRole, ctrl.ApproveLeaveRequest)
	v1.POST("/:id/reject", authMiddleware, staffRole, ctrl.RejectLeaveRequest)
	v1.DELETE("/:id", authMiddleware, middleware.RoleMiddleware("STUDENT", "PARENT"), ctrl.CancelLeaveRequest)
}
