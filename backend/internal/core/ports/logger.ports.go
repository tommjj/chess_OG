package ports

type Field struct {
	Name  string
	Value any
}

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// AllowLog determines if a message at messageLevel should be logged
func AllowLog(currentLevel, messageLevel LogLevel) bool {
	return messageLevel >= currentLevel
}

// LoggerPort defines the interface for logging functionalities
type LoggerPort interface {
	// Info logs an informational message with optional fields
	Info(message string, fields ...Field)
	// Debug logs a debug message with optional fields
	Debug(message string, fields ...Field)
	// Warn logs a warning message with optional fields
	Warn(message string, fields ...Field)
	// Error logs an error message with optional fields
	Error(message string, fields ...Field)
	// Fatal logs a fatal message with optional fields and terminates the application
	Fatal(message string, fields ...Field)
	// Log logs a message at the specified log level with optional fields
	// Defaults to ErrorLevel if an invalid level is provided
	Log(level LogLevel, message string, fields ...Field)
	// Level returns the current log level
	Level() LogLevel
}
