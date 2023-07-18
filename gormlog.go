package db

import (
	"context"
	"time"

	"github.com/devil-dwj/wms/log"

	gormlogger "gorm.io/gorm/logger"
)

const (
	skip = 3
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}
func (l *Logger) Info(ctx context.Context, format string, args ...interface{}) {
	log.Infof(format, args...)
}
func (l *Logger) Warn(ctx context.Context, format string, args ...interface{}) {
	log.Warnf(format, args...)
}
func (l *Logger) Error(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(format, args...)
}
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

	elapsed := time.Since(begin).Milliseconds()
	sql, rows := fc()
	if err != nil {
		log.LogWithOptions(
			log.WithLevel(log.LevelInfo),
			log.WithSkip(skip),
			log.WithKeyVals(
				"err", err.Error(),
				"elapsed", elapsed,
				"rows", rows,
				"sql", sql,
			),
		)
	} else {
		log.LogWithOptions(
			log.WithLevel(log.LevelInfo),
			log.WithSkip(skip),
			log.WithKeyVals(
				"elapsed", elapsed,
				"rows", rows,
				"sql", sql,
			),
		)
	}
}
