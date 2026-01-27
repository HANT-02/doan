package rest

import (
	xerror "doan/pkg/error"
	"doan/pkg/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseResponse struct {
	Success   bool        `json:"success"`
	Message   *string     `json:"message"`
	ErrorCode *string     `json:"error_code"`
	Data      interface{} `json:"data"`
}

func HandleError(c *gin.Context, err error) {
	var xErr *xerror.AppError
	if errors.As(err, &xErr) {
		fmt.Printf("HandleError, %s, %s\n", xErr.Error(), xErr.WithUserMessage(""))
		c.AbortWithStatusJSON(xErr.HTTPStatusCode(), &BaseResponse{
			Success:   false,
			Message:   utils.NewStringPtr(xErr.Message),
			ErrorCode: utils.NewStringPtr(xErr.Code),
			Data:      nil,
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, &BaseResponse{
		Success:   false,
		Message:   utils.NewStringPtr(xErr.Message),
		ErrorCode: utils.NewStringPtr(xErr.Code),
		Data:      nil,
	})
}

// ResponseSuccess sends a successful response
func ResponseSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, &BaseResponse{
		Success: true,
		Message: utils.NewStringPtr(message),
		Data:    data,
	})
}

// ResponseError sends an error response
func ResponseError(c *gin.Context, statusCode int, message string, err error) {
	errorCode := ""
	if err != nil {
		errorCode = err.Error()
	}
	c.AbortWithStatusJSON(statusCode, &BaseResponse{
		Success:   false,
		Message:   utils.NewStringPtr(message),
		ErrorCode: utils.NewStringPtr(errorCode),
		Data:      nil,
	})
}
