package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"os"
	"runtime"
	"strings"
	"time"
)

const TraceKey = "trace_id"

func GenerateTraceID(serviceName string) string {
	uUid := uuid.NewString()
	traceID := fmt.Sprintf("%s-%s", serviceName, strings.ReplaceAll(uUid, "-", ""))
	return traceID
}

func NewContextWithTraceID(ctx context.Context, serviceName string) context.Context {
	return context.WithValue(ctx, TraceKey, GenerateTraceID(serviceName))
}

func NewBackgroundContextWithTraceID(serviceName string) context.Context {
	return NewContextWithTraceID(context.Background(), serviceName)
}

type JSONLogger struct {
	Logger  log.Logger
	TraceID string
}

func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(3) // Adjust stack depth as needed
	if !ok {
		return "unknown"
	}
	projectRoot, _ := os.Getwd()
	relativePath := strings.TrimPrefix(file, fmt.Sprintf("%s/", projectRoot))
	return fmt.Sprintf("%s:%d", relativePath, line)
}

type logEntry struct {
	Time    string `json:"time"`
	Caller  string `json:"caller"`
	TraceID string `json:"trace_id,omitempty"`
	Msg     string `json:"msg"`
	Level   string `json:"level"`
}

func (l *JSONLogger) Log(level log.Level, keyvals ...interface{}) error {
	entry := logEntry{
		Time:    time.Now().Format(time.RFC3339),
		Level:   level.String(),
		TraceID: l.TraceID,
		Caller:  getCallerInfo(),
	}

	for i := 0; i < len(keyvals)-1; i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}

		switch key {
		case "caller":
			if val, ok := keyvals[i+1].(string); ok {
				entry.Caller = val
			}
		case "msg":
			if val, ok := keyvals[i+1].(string); ok {
				entry.Msg = val
			}
		default:
			// Handle other keys or ignore them
		}
	}

	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return l.Logger.Log(log.LevelInfo, "log", string(b))
}

func NewJSONLogger(traceID string) *JSONLogger {
	return &JSONLogger{
		Logger:  log.NewStdLogger(os.Stdout),
		TraceID: traceID,
	}
}

func NewLogger(ctx context.Context) *log.Helper {
	traceID, ok := ctx.Value(TraceKey).(string)
	if !ok {
		traceID = "unknown"
	}
	logger := NewJSONLogger(traceID)
	return log.NewHelper(logger)
}
