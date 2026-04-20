package services

import (
	"fmt"
	"os"

	"med-predict-backend/internal/config"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus logger with convenience methods
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a configured logger instance
func NewLogger(cfg *config.Config) *Logger {
	log := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// Create logs directory if needed
	if cfg.Env == "production" {
		os.MkdirAll(cfg.LogDir, 0755)

		// File logging for production
		errFile, err := os.OpenFile(
			fmt.Sprintf("%s/error.log", cfg.LogDir),
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0666,
		)
		if err == nil {
			log.SetOutput(errFile)
		}
	}

	// Set JSON formatter for production, text for development
	if cfg.Env == "production" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	return &Logger{log}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	entry := l.Logger.WithFields(toLogrusFields(fields))
	entry.Info(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	entry := l.Logger.WithFields(toLogrusFields(fields))
	entry.Warn(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...interface{}) {
	entry := l.Logger.WithFields(toLogrusFields(fields))
	entry.Error(msg)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	entry := l.Logger.WithFields(toLogrusFields(fields))
	entry.Debug(msg)
}

// toLogrusFields converts variadic key-value pairs to logrus Fields
func toLogrusFields(fields []interface{}) logrus.Fields {
	logFields := logrus.Fields{}
	for i := 0; i < len(fields)-1; i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		val := fields[i+1]
		logFields[key] = val
	}
	return logFields
}
