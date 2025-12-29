package log

import (
	"os"

	"github.com/tommjj/chess_OG/backend/internal/core/ports"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// zeroOrDefault returns the value if it's not the zero value of its type, otherwise returns the default value.
func zeroOrDefault[T comparable](value T, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}

// positiveOrDefault returns the value if it's positive, otherwise returns the default value.
func positiveOrDefault(value int, defaultValue int) int {
	if value <= 0 {
		return defaultValue
	}
	return value
}

const (
	FileMaxSizeDefault    = 100 // in MB
	FileMaxBackupsDefault = 7
	FileMaxAgeDefault     = 30 // in days
	EndCoderModeDefault   = "development"
)

func setDefaultLogConfig(conf *LogConfig) {
	conf.Level = zeroOrDefault(conf.Level, ports.InfoLevel)
	conf.FileMaxSize = positiveOrDefault(conf.FileMaxSize, FileMaxSizeDefault)          // 100 MB
	conf.FileMaxBackups = positiveOrDefault(conf.FileMaxBackups, FileMaxBackupsDefault) // 7 files
	conf.FileMaxAge = positiveOrDefault(conf.FileMaxAge, FileMaxAgeDefault)             // 30 days
	conf.EndCoderMode = zeroOrDefault(conf.EndCoderMode, EndCoderModeDefault)           // development mode
}

// mapping between ports.LogLevel and zapcore.Level
var logLevels = map[ports.LogLevel]zapcore.Level{
	ports.DebugLevel: zapcore.DebugLevel,
	ports.InfoLevel:  zapcore.InfoLevel,
	ports.WarnLevel:  zapcore.WarnLevel,
	ports.ErrorLevel: zapcore.ErrorLevel,
	ports.FatalLevel: zapcore.FatalLevel,
}

func newWriteSyncer(conf *LogConfig) zapcore.WriteSyncer {
	if conf.FilePath == "" {
		return zapcore.AddSync(os.Stdout)
	}

	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   conf.FilePath,
		MaxSize:    conf.FileMaxSize,
		MaxBackups: conf.FileMaxBackups,
		MaxAge:     conf.FileMaxAge,
	})

	// in production mode, only log to file
	if conf.EndCoderMode == "production" {
		return fileWriteSyncer
	}

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(fileWriteSyncer),
		zapcore.AddSync(os.Stdout),
	)
}

func newJSONEncoder(mode string) zapcore.EncoderConfig {
	if mode != "production" {
		return zap.NewDevelopmentEncoderConfig()
	}
	return zap.NewProductionEncoderConfig()
}
