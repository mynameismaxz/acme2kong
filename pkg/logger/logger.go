package logger

import (
	"log/slog"
	"os"
)

// Logger struct that holder the slog.Logger instance
type Logger struct {
	logger *slog.Logger
}

// New creates a new Logger instance
func New() *Logger {
	// new logger handler
	handler := NewJSONHandler(os.Stdout)

	l := slog.New(handler)

	return &Logger{
		logger: l,
	}
}

// Info logs a message with INFO level
func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

// Debug logs a message with DEBUG level
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

// Error logs an error level message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

// Warn logs a warning level message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}
