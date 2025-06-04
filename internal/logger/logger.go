package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level represents the severity of a log message
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String returns the string representation of a log level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel parses a string into a log level
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "WARN", "WARNING":
		return WarnLevel
	case "ERROR":
		return ErrorLevel
	case "FATAL":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Logger is the main logger interface
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	WithFields(fields ...Field) Logger
}

// StructuredLogger implements structured logging
type StructuredLogger struct {
	mu       sync.Mutex
	out      io.Writer
	level    Level
	fields   []Field
	jsonMode bool
}

// New creates a new structured logger
func New(out io.Writer, level Level, jsonMode bool) *StructuredLogger {
	return &StructuredLogger{
		out:      out,
		level:    level,
		jsonMode: jsonMode,
		fields:   []Field{},
	}
}

// NewDefault creates a logger with default settings
func NewDefault() *StructuredLogger {
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = "INFO"
	}
	
	jsonMode := os.Getenv("LOG_FORMAT") == "json"
	
	return New(os.Stdout, ParseLevel(levelStr), jsonMode)
}

// WithFields creates a new logger with additional fields
func (l *StructuredLogger) WithFields(fields ...Field) Logger {
	newFields := make([]Field, len(l.fields)+len(fields))
	copy(newFields, l.fields)
	copy(newFields[len(l.fields):], fields)
	
	return &StructuredLogger{
		out:      l.out,
		level:    l.level,
		fields:   newFields,
		jsonMode: l.jsonMode,
	}
}

// log writes a log entry
func (l *StructuredLogger) log(level Level, msg string, fields ...Field) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	entry := make(map[string]interface{})
	entry["timestamp"] = time.Now().Format(time.RFC3339)
	entry["level"] = level.String()
	entry["message"] = msg

	// Add default fields
	for _, field := range l.fields {
		entry[field.Key] = field.Value
	}

	// Add message-specific fields
	for _, field := range fields {
		entry[field.Key] = field.Value
	}

	// Add caller information for errors
	if level >= ErrorLevel {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			entry["caller"] = fmt.Sprintf("%s:%d", file, line)
		}
	}

	if l.jsonMode {
		// JSON output
		encoder := json.NewEncoder(l.out)
		encoder.Encode(entry)
	} else {
		// Human-readable output
		fmt.Fprintf(l.out, "[%s] %s %s",
			entry["timestamp"],
			entry["level"],
			entry["message"])
		
		// Add fields
		if len(l.fields)+len(fields) > 0 {
			fmt.Fprint(l.out, " ")
			first := true
			for _, field := range append(l.fields, fields...) {
				if !first {
					fmt.Fprint(l.out, " ")
				}
				fmt.Fprintf(l.out, "%s=%v", field.Key, field.Value)
				first = false
			}
		}
		
		// Add caller for errors
		if caller, ok := entry["caller"]; ok {
			fmt.Fprintf(l.out, " caller=%s", caller)
		}
		
		fmt.Fprintln(l.out)
	}
}

// Debug logs a debug message
func (l *StructuredLogger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields...)
}

// Info logs an info message
func (l *StructuredLogger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields...)
}

// Warn logs a warning message
func (l *StructuredLogger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields...)
}

// Error logs an error message
func (l *StructuredLogger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *StructuredLogger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields...)
	os.Exit(1)
}

// Helper functions for creating fields

// String creates a string field
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int creates an int field
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 creates an int64 field
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Bool creates a bool field
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Duration creates a duration field
func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value.String()}
}

// Err creates an error field
func Err(err error) Field {
	if err == nil {
		return Field{Key: "error", Value: nil}
	}
	return Field{Key: "error", Value: err.Error()}
}

// Any creates a field with any value
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}