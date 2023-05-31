// Package logger provides a logger interface.
// It is mapped to logrus and can be replaced with any other logger.
package logger

import "github.com/sirupsen/logrus"

type (
	// Logger is a logger alias for logrus.Logger.
	Logger = logrus.Logger
	// Level is a logger alias for logrus.Level.
	Level = logrus.Level
)

// Log levels.
const (
	LevelDebug = logrus.DebugLevel
	LevelInfo  = logrus.InfoLevel
	LevelWarn  = logrus.WarnLevel
	LevelError = logrus.ErrorLevel
)

// NewLogger creates a new logger.
func NewLogger(level Level) *Logger {
	l := logrus.New()
	l.SetLevel(level)
	return l
}

// NewNoOpLogger creates a new logger
// that does not log anything except panics.
func NewNoOpLogger() *Logger {
	l := logrus.New()
	l.SetLevel(logrus.PanicLevel)
	return l
}
