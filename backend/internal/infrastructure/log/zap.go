// logger adapter using zap

package log

import (
	"fmt"

	"github.com/tommjj/chess_OG/backend/internal/core/ports"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogConfig holds the configuration for the ZapLogger
type LogConfig struct {
	// Level defines the logging level
	//
	// Possible values are:
	// - ports.DebugLevel
	// - ports.InfoLevel
	// - ports.WarnLevel
	// - ports.ErrorLevel
	// - ports.FatalLevel
	Level ports.LogLevel

	// FilePath specifies the path to the log file if zero log is not used
	FilePath string
	// FileMaxSize specifies the maximum size in megabytes of the log file before it gets rotated
	FileMaxSize int
	// FileMaxBackups specifies the maximum number of old log files to retain
	FileMaxBackups int
	// FileMaxAge specifies the maximum number of days to retain old log files
	FileMaxAge int
	// EndCoderMode mode specifies the encoder mode: "production" or "development"
	// if value is non-standard, "development" mode is used
	EndCoderMode string
}

// ZapLogger is a logger implementation using zap
type ZapLogger struct {
	logger *zap.Logger
	level  ports.LogLevel
}

// Sync flushes any buffered log entries
func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}

func NewZapLogger(conf LogConfig) (*ZapLogger, error) {
	level, ok := logLevels[conf.Level]
	if !ok {
		return nil, fmt.Errorf("invalid log level: %v", conf.Level)
	}

	setDefaultLogConfig(&conf)

	writeSyncer := newWriteSyncer(&conf)
	encoderConfig := newJSONEncoder(conf.EndCoderMode)
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, writeSyncer, level)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &ZapLogger{logger: logger, level: conf.Level}, nil
}

// convertFields converts ports.Field to zap.Field
func convertFields(fields ...ports.Field) []zap.Field {
	if len(fields) == 0 {
		return nil
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Name, field.Value))
	}
	return zapFields
}

func (z *ZapLogger) Info(message string, fields ...ports.Field) {
	z.logger.Info(message, convertFields(fields...)...)
}

func (z *ZapLogger) Debug(message string, fields ...ports.Field) {
	z.logger.Debug(message, convertFields(fields...)...)
}

func (z *ZapLogger) Warn(message string, fields ...ports.Field) {
	z.logger.Warn(message, convertFields(fields...)...)
}

func (z *ZapLogger) Error(message string, fields ...ports.Field) {
	z.logger.Error(message, convertFields(fields...)...)
}

func (z *ZapLogger) Fatal(message string, fields ...ports.Field) {
	z.logger.Fatal(message, convertFields(fields...)...)
}

func (z *ZapLogger) Log(level ports.LogLevel, message string, fields ...ports.Field) {
	zapLevel, ok := logLevels[level]
	if !ok { // default to ErrorLevel if invalid level is provided
		zapLevel = zapcore.ErrorLevel
	}
	z.logger.Log(zapLevel, message, convertFields(fields...)...)
}

func (z *ZapLogger) Level() ports.LogLevel {
	return z.level
}
