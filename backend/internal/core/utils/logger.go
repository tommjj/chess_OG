package utils

import (
	"fmt"

	"github.com/tommjj/chess_OG/backend/internal/core/ports"
)

// LoggerFmtWrapper is a wrapper around LoggerPort that adds formatted logging methods.
type LoggerFmtWrapper struct {
	logger ports.LoggerPort
}

// NewLoggerFmtWrapper creates a new LoggerFmtWrapper.
func NewLoggerFmtWrapper(logger ports.LoggerPort) *LoggerFmtWrapper {
	return &LoggerFmtWrapper{logger: logger}
}

func (l *LoggerFmtWrapper) Level() ports.LogLevel {
	return l.logger.Level()
}

func (l *LoggerFmtWrapper) CanLog(messageLevel ports.LogLevel) bool {
	return ports.AllowLog(l.Level(), messageLevel)
}

// Infof logs an informational message with formatting.
func (l *LoggerFmtWrapper) Infof(format string, args ...any) {
	if !l.CanLog(ports.InfoLevel) {
		return
	}

	message := fmt.Sprintf(format, args...)
	l.logger.Info(message)
}

// Debugf logs a debug message with formatting.
func (l *LoggerFmtWrapper) Debugf(format string, args ...any) {
	if !l.CanLog(ports.DebugLevel) {
		return
	}

	message := fmt.Sprintf(format, args...)
	l.logger.Debug(message)
}

// Warnf logs a warning message with formatting.
func (l *LoggerFmtWrapper) Warnf(format string, args ...any) {
	if !l.CanLog(ports.DebugLevel) {
		return
	}

	message := fmt.Sprintf(format, args...)
	l.logger.Warn(message)
}

// Errorf logs an error message with formatting.
func (l *LoggerFmtWrapper) Errorf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.logger.Error(message)
}

// Fatalf logs a fatal message with formatting and terminates the application.
func (l *LoggerFmtWrapper) Fatalf(format string, args ...any) {
	if !l.CanLog(ports.DebugLevel) {
		return
	}
	message := fmt.Sprintf(format, args...)
	l.logger.Fatal(message)
}

// Logf logs a message at the specified log level with formatting.
func (l *LoggerFmtWrapper) Logf(level ports.LogLevel, format string, args ...any) {
	if !l.CanLog(level) {
		return
	}

	message := fmt.Sprintf(format, args...)
	l.logger.Log(level, message)
}

// Logger returns the underlying LoggerPort.
func (l *LoggerFmtWrapper) Logger() ports.LoggerPort {
	return l.logger
}
