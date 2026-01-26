package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for logging
type Logger interface {
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Warn(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
	Fatal(ctx context.Context, msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
	ErrorWithStack(ctx context.Context, msg string, err error, stack string, keysAndValues ...interface{})
}

// TraceIDKey là key cho trace ID trong context
type TraceIDKey struct{}

// Config represents the configuration for the logger
type Config struct {
	Level       string `json:"level" yaml:"level"`
	Format      string `json:"format" yaml:"format"`
	Output      string `json:"output" yaml:"output"`
	TimeFormat  string `json:"time_format" yaml:"time_format"`
	ServiceName string `json:"service_name" yaml:"service_name"`
	Environment string `json:"environment" yaml:"environment"`
}

// ZapLogger implements the Logger interface using Zap
type ZapLogger struct {
	logger      *zap.SugaredLogger
	serviceName string
	environment string
}

// NewLogger creates a new logger instance
func NewZapLogger(config Config) Logger {
	var level zapcore.Level
	switch strings.ToLower(config.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	var output io.Writer = os.Stdout
	if config.Output != "" && config.Output != "stdout" {
		file, err := os.OpenFile(config.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		} else {
			output = file
		}
	}

	timeFormat := config.TimeFormat
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(timeFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if strings.ToLower(config.Format) == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(output), zap.NewAtomicLevelAt(level))
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// Thêm thông tin cơ bản
	fields := []interface{}{
		"service", config.ServiceName,
		"env", config.Environment,
	}

	return &ZapLogger{
		logger:      zapLogger.Sugar().With(fields...),
		serviceName: config.ServiceName,
		environment: config.Environment,
	}
}

// extractTraceID lấy traceID từ context
func extractTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// Xem có traceID trong context theo key mới
	if traceID, ok := ctx.Value(TraceIDKey{}).(string); ok && traceID != "" {
		return traceID
	}

	// Kiểm tra các key thông dụng cho backwards compatibility
	for _, key := range []string{"trace_id", "traceID", "X-Request-ID", "requestID", "request_id"} {
		if traceID, ok := ctx.Value(key).(string); ok && traceID != "" {
			return traceID
		}
	}

	return ""
}

// generateLoggerFields tạo các fields phổ biến cho log
func (l *ZapLogger) generateLoggerFields(ctx context.Context, additionalFields ...interface{}) []interface{} {
	// Thêm traceID nếu có
	fields := additionalFields[:]
	traceID := extractTraceID(ctx)
	if traceID != "" {
		fields = append(fields, "trace_id", traceID)
	}

	return fields
}

// Debug logs a debug message
func (l *ZapLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	fields := l.generateLoggerFields(ctx, keysAndValues...)
	l.logger.Debugw(msg, fields...)
}

// Info logs an informational message
func (l *ZapLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	fields := l.generateLoggerFields(ctx, keysAndValues...)
	l.logger.Infow(msg, fields...)
}

// Warn logs a warning message
func (l *ZapLogger) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	fields := l.generateLoggerFields(ctx, keysAndValues...)
	l.logger.Warnw(msg, fields...)
}

// Error logs an error message
func (l *ZapLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	fields := l.generateLoggerFields(ctx, keysAndValues...)
	l.logger.Errorw(msg, fields...)
}

// ErrorWithStack logs an error message with full stack trace
func (l *ZapLogger) ErrorWithStack(ctx context.Context, msg string, err error, stack string, keysAndValues ...interface{}) {
	baseFields := []interface{}{
		"error", err.Error(),
	}

	// Thêm stack trace nếu được cung cấp
	if stack != "" {
		baseFields = append(baseFields, "stack_trace", stack)
	} else {
		// Tự tạo stack trace nếu không có
		baseFields = append(baseFields, "stack_trace", captureStackTrace())
	}

	// Kết hợp với các fields khác
	fields := append(baseFields, keysAndValues...)
	finalFields := l.generateLoggerFields(ctx, fields...)

	l.logger.Errorw(msg, finalFields...)
}

// Fatal logs a fatal message and exits
func (l *ZapLogger) Fatal(ctx context.Context, msg string, keysAndValues ...interface{}) {
	fields := l.generateLoggerFields(ctx, keysAndValues...)
	l.logger.Fatalw(msg, fields...)
}

// With returns a logger with the given key-value pairs
func (l *ZapLogger) With(keysAndValues ...interface{}) Logger {
	return &ZapLogger{
		logger:      l.logger.With(keysAndValues...),
		serviceName: l.serviceName,
		environment: l.environment,
	}
}

// captureStackTrace tạo stack trace cho vị trí gọi hiện tại
func captureStackTrace() string {
	const depth = 32
	var pcs [depth]uintptr

	// skip 3 frames: captureStackTrace, ErrorWithStack/Error, và log caller
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var builder strings.Builder
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&builder, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return builder.String()
}
