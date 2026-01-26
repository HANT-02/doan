package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"doan/pkg/logger"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger implements GORM's logger interface using application logger
type GormLogger struct {
	logger                    logger.Logger
	level                     gormlogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

// NewGormLogger creates a new GORM logger instance
func NewGormLogger(appLogger logger.Logger, slowThreshold time.Duration, logLevel gormlogger.LogLevel, ignoreRecordNotFoundError bool) *GormLogger {
	return &GormLogger{
		logger:                    appLogger,
		level:                     logLevel,
		slowThreshold:             slowThreshold,
		ignoreRecordNotFoundError: ignoreRecordNotFoundError,
	}
}

// LogMode sets log level for GORM
func (g *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *g
	newLogger.level = level
	return &newLogger
}

// Info logs info level messages
func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.level >= gormlogger.Info {
		g.logger.Info(ctx, fmt.Sprintf(msg, data...))
	}
}

// Warn logs warning level messages
func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.level >= gormlogger.Warn {
		g.logger.Warn(ctx, fmt.Sprintf(msg, data...))
	}
}

// Error logs error level messages
func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.level >= gormlogger.Error {
		g.logger.Error(ctx, fmt.Sprintf(msg, data...))
	}
}

// Trace logs SQL queries with execution time
func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.level <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []interface{}{
		"elapsed_ms", float64(elapsed.Nanoseconds()) / 1e6,
		"rows", rows,
		"sql", sql,
	}

	switch {
	case err != nil && g.level >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !g.ignoreRecordNotFoundError):
		g.logger.Error(ctx, "gorm query error", append(fields, "error", err.Error())...)
	case elapsed > g.slowThreshold && g.slowThreshold != 0 && g.level >= gormlogger.Warn:
		g.logger.Warn(ctx, "slow sql query", append(fields, "threshold_ms", float64(g.slowThreshold.Nanoseconds())/1e6)...)
	case g.level == gormlogger.Info:
		g.logger.Debug(ctx, "gorm query", fields...)
	}
}
