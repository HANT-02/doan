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
	var xErr *xerror.Error
	if errors.As(err, &xErr) {
		fmt.Printf("HandleError, %s, %s\n", xErr.ErrCode(), xErr.Message())
		c.AbortWithStatusJSON(xErr.SttCode(), &BaseResponse{
			Success:   false,
			Message:   utils.NewString(xErr.Message()),
			ErrorCode: utils.NewString(xErr.ErrCode()),
			Data:      nil,
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, &BaseResponse{
		Success:   false,
		Message:   utils.NewString(err.Error()),
		ErrorCode: nil,
		Data:      nil,
	})
}
