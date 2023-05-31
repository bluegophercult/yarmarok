// Package logger provides a logger interface.
// It is mapped to logrus and can be replaced with any other logger.
package logger

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type (
	// Logger is a logger alias for logrus.Logger.
	Logger = logrus.Logger
	// Level is a logger alias for logrus.Level.
	Level = logrus.Level
	// Fields is a logger alias for logrus.Fields.
	Fields = logrus.Fields
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

// ResponseMetric is a struct that captures
// response status code and size.
type ResponseMetric struct {
	Status int
	Size   int
}

// LoggingResponseWriter is a wrapper around http.ResponseWriter
// that captures response status code and size.
type LoggingResponseWriter struct {
	http.ResponseWriter
	responseMetric *ResponseMetric
}

// NewLoggingResponseWriter creates a new LoggingResponseWriter.
func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{
		ResponseWriter: w,
		responseMetric: &ResponseMetric{
			Status: 0,
			Size:   0,
		},
	}
}

// WriteHeader captures response status code.
func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseMetric.Size += size
	return size, err
}

// WriteHeader captures response status code.
func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseMetric.Status = statusCode
}

// ResponseMetric returns response data.
func (r *LoggingResponseWriter) ResponseMetric() *ResponseMetric {
	return r.responseMetric
}
