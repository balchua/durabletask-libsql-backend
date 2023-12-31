package logging

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

// SlogWrapper implements the Logger interface using slog
type SlogWrapper struct {
	logger *slog.Logger
	opts   *slog.HandlerOptions
}

// NewSlogWrapper creates a new SlogWrapper instance
func NewSlogWrapper(options ...Option) SlogWrapper {
	defaultOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	w := SlogWrapper{
		opts: defaultOpts,
	}
	for _, opt := range options {
		opt(&w)
	}

	w.logger = slog.New(slog.NewJSONHandler(os.Stdout, w.opts))

	slog.SetDefault(w.logger)
	return w
}

// Options for configuring the SlogWrapper
type Option func(*SlogWrapper)

// WithLogLevel sets the log level for the wrapper
func WithLogLevel(level slog.Level) Option {
	return func(w *SlogWrapper) {
		w.opts.Level = level
	}
}

// Debug logs a message at level Debug
func (w SlogWrapper) Debug(v ...any) {
	w.logger.Debug(fmt.Sprint(v...))
}

// Debugf logs a message at level Debug
func (w SlogWrapper) Debugf(format string, v ...any) {
	w.logger.Debug(fmt.Sprintf(format, v...))
}

// Debugf logs a structured message at level Debug
func (w SlogWrapper) DebugS(format string, v ...any) {
	w.logger.Debug(format, v...)
}

// Info logs a message at level Info
func (w SlogWrapper) Info(v ...any) {
	w.logger.Info(fmt.Sprint(v...))
}

// Infof logs a message at level Info
func (w SlogWrapper) Infof(format string, v ...any) {
	w.logger.Info(fmt.Sprintf(format, v...))
}

// InfoS logs a structured message at level Info
func (w SlogWrapper) InfoS(format string, v ...any) {
	w.logger.Info(format, v...)
}

// Warn logs a message at level Warn
func (w SlogWrapper) Warn(v ...any) {
	w.logger.Warn(fmt.Sprint(v...))
}

// Warnf logs a message at level Warn
func (w SlogWrapper) Warnf(format string, v ...any) {
	w.logger.Warn(fmt.Sprintf(format, v...))
}

// WarnS logs a message at level Warn
func (w SlogWrapper) WarnS(format string, v ...any) {
	w.logger.Warn(format, v...)
}

// Error logs a message at level Error
func (w SlogWrapper) Error(v ...any) {
	w.logger.Error(fmt.Sprint(v...))
}

// Errorf logs a message at level Error
func (w SlogWrapper) Errorf(format string, v ...any) {
	w.logger.Error(fmt.Sprintf(format, v...))
}

// ErrorS logs a message at level Error
func (w SlogWrapper) ErrorS(msg string, v ...any) {
	_, file, line, _ := runtime.Caller(1)
	errorMsg := fmt.Sprintf("%s at %s:%d", msg, file, line)
	w.logger.Error(errorMsg, v...)
}
