package logtool

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger is a GORM logger that uses the local logtool
type GormLogger struct {
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
}

// NewGormLogger creates a new GORM logger with logtool
func NewGormLogger() *GormLogger {
	return &GormLogger{
		LogLevel:                  gormlogger.Info,
		SlowThreshold:             200 * time.Millisecond, // Hardcoded to 200ms
		Colorful:                  false,
		IgnoreRecordNotFoundError: true,
	}
}

// NewGormLoggerWithConfig creates a new GORM logger with custom configuration
func NewGormLoggerWithConfig(level gormlogger.LogLevel, ignoreRecordNotFound bool) *GormLogger {
	return &GormLogger{
		LogLevel:                  level,
		SlowThreshold:             200 * time.Millisecond, // Hardcoded to 200ms
		Colorful:                  false,
		IgnoreRecordNotFoundError: ignoreRecordNotFound,
	}
}

// LogMode sets the log level
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs info level messages
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		getLogger().Infof(msg, data...)
	}
}

// Warn logs warn level messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		getLogger().Warnf(msg, data...)
	}
}

// Error logs error level messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		getLogger().Errorf(msg, data...)
	}
}

// Trace logs SQL queries - always logs at trace level regardless of configured level
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// Create a logger with caller skip 3 for GORM to show the actual database call location
	logger := getLogger().Desugar().WithOptions(zap.AddCallerSkip(3)).Sugar()
	// logger := GetLogger()

	// Log slow queries as warnings
	if elapsed > l.SlowThreshold {
		logger.Warnw("Slow SQL query detected",
			"sql", sql,
			"rows", rows,
			"elapsed", elapsed,
			"threshold", l.SlowThreshold,
		)
		return
	}

	// Log based on error
	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
		logger.Errorw("SQL query error",
			"sql", sql,
			"rows", rows,
			"elapsed", elapsed,
			"error", err,
		)
		return
	}

	// Always log successful queries at trace level
	// logger.Debugw("SQL query executed",
	// 	"sql", sql,
	// 	"rows", rows,
	// 	"elapsed", elapsed,
	// )
}
