package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func IgnorePath(ginCtx *gin.Context) bool {
	mapIgnorePaths := map[string]string{
		"/api/v1/ping":         "GET",
		"/api/v1/health-check": "GET",
	}

	requestPath := ginCtx.FullPath()
	if method, ok := mapIgnorePaths[requestPath]; ok && strings.ToLower(method) == strings.ToLower(ginCtx.Request.Method) {
		return true
	}
	if strings.Contains(requestPath, "/api/swagger") {
		return true
	}
	return false
}

//func GetUserID(ctx context.Context) string {
//	userID, _ := ctx.Value(constants.ContextUserID).(string)
//	return userID
//}
