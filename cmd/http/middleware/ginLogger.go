package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

type GinLog struct {
	Level      string `json:"level"`
	Time       string `json:"time"`
	StatusCode int    `json:"status_code"`
	Latency    int64  `json:"latency"`
	ClientIP   string `json:"client_ip"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	UserAgent  string `json:"user_agent"`
	Error      string `json:"error"`
}

func JsonLogMiddleware(param gin.LogFormatterParams) string {
	record := &GinLog{
		Level:      "gin",
		Time:       param.TimeStamp.Format(time.RFC3339),
		StatusCode: param.StatusCode,
		Latency:    param.Latency.Microseconds(),
		ClientIP:   param.ClientIP,
		Method:     param.Method,
		Path:       param.Path,
		UserAgent:  param.Request.UserAgent(),
		Error:      param.ErrorMessage,
	}
	bytes, err := json.Marshal(record)
	if err != nil {
		panic(err)
	}
	return string(bytes) + "\n"
}
