package btclog

import (
	"context"
	"fmt"
	"log/slog"
)

// Handler wraps the slog.Handler interface with a few more methods that we
// need in order to satisfy the Logger interface.
type Handler interface {
	slog.Handler

	// Level returns the current logging level of the Handler.
	Level() Level

	// SetLevel changes the logging level of the Handler to the passed
	// level.
	SetLevel(level Level)

	// SubSystem returns a copy of the given handler but with the new tag.
	SubSystem(tag string) Handler
}

// sLogger is an implementation of Logger backed by a structured sLogger.
type sLogger struct {
	Handler
	logger *slog.Logger
}

// NewSLogger constructs a new structured logger from the given Handler.
func NewSLogger(handler Handler) Logger {
	return &sLogger{
		Handler: handler,
		logger:  slog.New(handler),
	}
}

// Tracef formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelTrace.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Tracef(format string, params ...any) {
	l.toSlogf(LevelTrace, format, params...)
}

// Debugf formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelDebug.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Debugf(format string, params ...any) {
	l.toSlogf(LevelDebug, format, params...)
}

// Infof formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelInfo.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Infof(format string, params ...any) {
	l.toSlogf(LevelInfo, format, params...)
}

// Warnf formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelWarn.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Warnf(format string, params ...any) {
	l.toSlogf(LevelWarn, format, params...)
}

// Errorf formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelError.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Errorf(format string, params ...any) {
	l.toSlogf(LevelError, format, params...)
}

// Criticalf formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelCritical.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Criticalf(format string, params ...any) {
	l.toSlogf(LevelCritical, format, params...)
}

// Trace formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelTrace.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Trace(v ...any) {
	l.toSlog(LevelTrace, v...)
}

// Debug formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelDebug.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Debug(v ...any) {
	l.toSlog(LevelDebug, v...)
}

// Info formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelInfo.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Info(v ...any) {
	l.toSlog(LevelInfo, v...)
}

// Warn formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelWarn.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Warn(v ...any) {
	l.toSlog(LevelWarn, v...)
}

// Error formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelError.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Error(v ...any) {
	l.toSlog(LevelError, v...)
}

// Critical formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelCritical.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Critical(v ...any) {
	l.toSlog(LevelCritical, v...)
}

// TraceS writes a structured log with the given message and key-value pair
// attributes with LevelTrace to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) TraceS(ctx context.Context, msg string, attrs ...any) {
	l.toSlogS(ctx, LevelTrace, msg, attrs...)
}

// DebugS writes a structured log with the given message and key-value pair
// attributes with LevelDebug to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) DebugS(ctx context.Context, msg string, attrs ...any) {
	l.toSlogS(ctx, LevelDebug, msg, attrs...)
}

// InfoS writes a structured log with the given message and key-value pair
// attributes with LevelInfo to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) InfoS(ctx context.Context, msg string, attrs ...any) {
	l.toSlogS(ctx, LevelInfo, msg, attrs...)
}

// WarnS writes a structured log with the given message and key-value pair
// attributes with LevelWarn to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) WarnS(ctx context.Context, msg string, err error,
	attrs ...any) {

	if err != nil {
		attrs = append([]any{slog.String("err", err.Error())}, attrs...)
	}

	l.toSlogS(ctx, LevelWarn, msg, attrs...)
}

// ErrorS writes a structured log with the given message and key-value pair
// attributes with LevelError to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) ErrorS(ctx context.Context, msg string, err error,
	attrs ...any) {

	if err != nil {
		attrs = append([]any{slog.String("err", err.Error())}, attrs...)
	}

	l.toSlogS(ctx, LevelError, msg, attrs...)
}

// CriticalS writes a structured log with the given message and key-value pair
// attributes with LevelCritical to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) CriticalS(ctx context.Context, msg string, err error,
	attrs ...any) {
	if err != nil {
		attrs = append([]any{slog.String("err", err.Error())}, attrs...)
	}

	l.toSlogS(ctx, LevelCritical, msg, attrs...)
}

// toSlogf is a helper method that converts an unstructured log call that
// contains a format string and parameters for the string into the appropriate
// form expected by the structured logger.
func (l *sLogger) toSlogf(level Level, format string, params ...any) {
	l.logger.Log(context.Background(), slog.Level(level),
		fmt.Sprintf(format, params...))
}

// toSlog is a helper method that converts an unstructured log call that
// contains a number of parameters into the appropriate form expected by the
// structured logger.
func (l *sLogger) toSlog(level Level, v ...any) {
	l.logger.Log(context.Background(), slog.Level(level), fmt.Sprint(v...))
}

// toSlogS is a helper method that can be used by all the structured log calls
// to access the underlying logger.
func (l *sLogger) toSlogS(ctx context.Context, level Level, msg string,
	attrs ...any) {

	l.logger.Log(ctx, slog.Level(level), msg, mergeAttrs(ctx, attrs)...)
}

var _ Logger = (*sLogger)(nil)
