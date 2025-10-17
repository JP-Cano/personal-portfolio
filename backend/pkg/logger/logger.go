package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// Logger is the main logger structure
type Logger struct {
	level      LogLevel
	logger     *log.Logger
	useColor   bool
	includePos bool // Include file position in logs
}

var defaultLogger *Logger

// Config holds logger configuration
type Config struct {
	Level      LogLevel
	Output     io.Writer
	UseColor   bool
	IncludePos bool // Include file:line in log output
}

// Init initializes the default logger with the given configuration
func Init(config Config) {
	if config.Output == nil {
		config.Output = os.Stdout
	}

	defaultLogger = &Logger{
		level:      config.Level,
		logger:     log.New(config.Output, "", 0), // We'll format ourselves
		useColor:   config.UseColor,
		includePos: config.IncludePos,
	}
}

// GetLogger returns the default logger instance
func GetLogger() *Logger {
	if defaultLogger == nil {
		// Initialize with default config if not already initialized
		Init(Config{
			Level:      INFO,
			Output:     os.Stdout,
			UseColor:   true,
			IncludePos: true,
		})
	}
	return defaultLogger
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// getColor returns the color code for a log level
func (l *Logger) getColor(level LogLevel) string {
	if !l.useColor {
		return ""
	}
	switch level {
	case DEBUG:
		return colorCyan
	case INFO:
		return colorGreen
	case WARN:
		return colorYellow
	case ERROR:
		return colorRed
	case FATAL:
		return colorPurple
	default:
		return colorWhite
	}
}

// log is the internal logging function
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)

	// Get caller information
	var position string
	if l.includePos {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			position = fmt.Sprintf(" [%s:%d]", filepath.Base(file), line)
		}
	}

	color := l.getColor(level)
	reset := ""
	if l.useColor {
		reset = colorReset
	}

	logMessage := fmt.Sprintf("%s[%s]%s [%s]%s %s",
		color,
		timestamp,
		reset,
		level.String(),
		position,
		message,
	)

	l.logger.Println(logMessage)

	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// WithFields returns a new logger with context fields
func (l *Logger) WithFields(fields map[string]interface{}) *ContextLogger {
	return &ContextLogger{
		logger: l,
		fields: fields,
	}
}

// ContextLogger is a logger with context fields
type ContextLogger struct {
	logger *Logger
	fields map[string]interface{}
}

// formatFields formats the context fields
func (cl *ContextLogger) formatFields() string {
	if len(cl.fields) == 0 {
		return ""
	}

	result := " |"
	for key, value := range cl.fields {
		result += fmt.Sprintf(" %s=%v", key, value)
	}
	return result
}

// Debug logs a debug message with context
func (cl *ContextLogger) Debug(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...) + cl.formatFields()
	cl.logger.Debug(message)
}

// Info logs an info message with context
func (cl *ContextLogger) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...) + cl.formatFields()
	cl.logger.Info(message)
}

// Warn logs a warning message with context
func (cl *ContextLogger) Warn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...) + cl.formatFields()
	cl.logger.Warn(message)
}

// Error logs an error message with context
func (cl *ContextLogger) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...) + cl.formatFields()
	cl.logger.Error(message)
}

// Fatal logs a fatal message with context and exits
func (cl *ContextLogger) Fatal(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...) + cl.formatFields()
	cl.logger.Fatal(message)
}

// Package-level convenience functions that use the default logger
func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Fatal(format, args...)
}

func WithFields(fields map[string]interface{}) *ContextLogger {
	return GetLogger().WithFields(fields)
}
