package utils

import "go.uber.org/zap"

// Logger is a simple wrapper around zap.Logger for compatibility
type Logger struct {
	zapLogger *zap.Logger
}

// NewLogger creates a new Logger instance
func NewLogger(zapLogger *zap.Logger) *Logger {
	return &Logger{
		zapLogger: zapLogger,
	}
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.zapLogger.Info(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.zapLogger.Error(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.zapLogger.Warn(msg)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.zapLogger.Debug(msg)
}
