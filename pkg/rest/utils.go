package rest

import "github.com/gin-gonic/gin"

// ErrorResponse định nghĩa cấu trúc phản hồi lỗi
type ErrorResponse struct {
	Success     bool                   `json:"success"`
	Error       string                 `json:"error"`
	Code        string                 `json:"code,omitempty"`
	Description string                 `json:"description,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}
type StandardResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, res interface{}) {
	c.JSON(200, StandardResponse{
		Success: true,
		Data:    res,
	})
}
