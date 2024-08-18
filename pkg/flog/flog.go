package flog

import (
	"log"
	"os"

	"forum/pkg/color"
)

// LogLevel represents the severity level of a log message.
type LogLevel int

const (
	// LevelDebug represents debug-level log messages.
	LevelDebug LogLevel = iota
	// LevelInfo represents informational log messages.
	LevelInfo
	// LevelWarn represents warning log messages.
	LevelWarn
	// LevelError represents error log messages.
	LevelError
)

// Logger provides a structured logging mechanism with support for log levels.
type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	// level specifies the minimum log level to be logged.
	level LogLevel
}

// NewLogger creates and initializes a new Logger instance with the specified minimum log level.
func NewLogger(level int) *Logger {

	return &Logger{
		debugLogger: log.New(os.Stdout, color.Make("DEBUG\t", color.Blue), log.Ldate|log.Ltime),
		infoLogger:  log.New(os.Stdout, color.Make("INFO\t", color.Green), log.Ldate|log.Ltime),
		warnLogger:  log.New(os.Stdout, color.Make("WARN\t", color.Yellow), log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, color.Make("ERROR\t", color.Red), log.Ldate|log.Ltime),
		level:       LogLevel(level),
	}
}

// Debug logs a debug-level message.
func (l *Logger) Debug(format string, v ...any) {
	if l.level <= LevelDebug {
		l.debugLogger.Printf(format, v...)
	}
}

// Info logs an informational message.
func (l *Logger) Info(format string, v ...any) {
	if l.level <= LevelInfo {
		l.infoLogger.Printf(format, v...)
	}
}

// Warn logs a warning message.
func (l *Logger) Warn(format string, v ...any) {
	if l.level <= LevelWarn {
		l.warnLogger.Printf(format, v...)
	}
}

// Error logs an error message.
func (l *Logger) Error(format string, v ...any) {
	if l.level <= LevelError {
		l.errorLogger.Printf(format, v...)
	}
}
