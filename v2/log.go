package btclog

import (
	"context"
	"fmt"
	"github.com/btcsuite/btclog"
	"io"
	"log/slog"
)

// Disabled is a Logger that will never output anything.
var Disabled Logger

// Handler wraps the slog.Handler interface with a few more methods that we
// need in order to satisfy the Logger interface.
type Handler interface {
	slog.Handler

	// Level returns the current logging level of the Handler.
	Level() btclog.Level

	// SetLevel changes the logging level of the Handler to the passed
	// level.
	SetLevel(level btclog.Level)

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

// Tracef creates a formatted message from the to format specifier along with
// any parameters then writes it to the logger with LevelTrace.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Tracef(format string, params ...any) {
	l.toSlogf(levelTrace, format, params...)
}

// Debugf creates a formatted message from the to format specifier along with
// any parameters then writes it to the logger with LevelDebug.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Debugf(format string, params ...any) {
	l.toSlogf(levelDebug, format, params...)
}

// Infof creates a formatted message from the to format specifier along with
// any parameters then writes it to the logger with LevelInfo.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Infof(format string, params ...any) {
	l.toSlogf(levelInfo, format, params...)
}

// Warnf creates a formatted message from the to format specifier along with
// any parameters then writes it to the logger with LevelWarn.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Warnf(format string, params ...any) {
	l.toSlogf(levelWarn, format, params...)
}

// Errorf creates a formatted message from the to format specifier along with
// any parameters then writes it to the logger with LevelError.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Errorf(format string, params ...any) {
	l.toSlogf(levelError, format, params...)
}

// Criticalf creates a formatted message from the to format specifier along
// with any parameters then writes it to the logger with LevelCritical.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Criticalf(format string, params ...any) {
	l.toSlogf(levelCritical, format, params...)
}

// Trace formats a message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelTrace.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Trace(v ...any) {
	l.toSlog(levelTrace, v...)
}

// Debug formats a message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelDebug.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Debug(v ...any) {
	l.toSlog(levelDebug, v...)
}

// Info formats a message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelInfo.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Info(v ...any) {
	l.toSlog(levelInfo, v...)
}

// Warn formats a message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelWarn.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Warn(v ...any) {
	l.toSlog(levelWarn, v...)
}

// Error formats a message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelError.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Error(v ...any) {
	l.toSlog(levelError, v...)
}

// Critical formats a message using the default formats for its operands,
// prepends the prefix as necessary, and writes to log with LevelCritical.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Critical(v ...any) {
	l.toSlog(levelCritical, v...)
}

// TraceS writes a structured log with the given message and key-value pair
// attributes with LevelTrace to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) TraceS(ctx context.Context, msg string, attrs ...any) {
	l.toSlogS(ctx, levelTrace, msg, attrs...)
}

// DebugS writes a structured log with the given message and key-value pair
// attributes with LevelDebug to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) DebugS(ctx context.Context, msg string, attrs ...any) {
	l.toSlogS(ctx, levelDebug, msg, attrs...)
}

// InfoS writes a structured log with the given message and key-value pair
// attributes with LevelInfo to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) InfoS(ctx context.Context, msg string, attrs ...any) {
	l.toSlogS(ctx, levelInfo, msg, attrs...)
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

	l.toSlogS(ctx, levelWarn, msg, attrs...)
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

	l.toSlogS(ctx, levelError, msg, attrs...)
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

	l.toSlogS(ctx, levelCritical, msg, attrs...)
}

// toSlogf is a helper method that converts an unstructured log call that
// contains a format string and parameters for the string into the appropriate
// form expected by the structured logger.
func (l *sLogger) toSlogf(level slog.Level, format string, params ...any) {
	l.logger.Log(context.Background(), level,
		fmt.Sprintf(format, params...))
}

// toSlog is a helper method that converts an unstructured log call that
// contains a number of parameters into the appropriate form expected by the
// structured logger.
func (l *sLogger) toSlog(level slog.Level, v ...any) {
	l.logger.Log(context.Background(), level, fmt.Sprint(v...))
}

// toSlogS is a helper method that can be used by all the structured log calls
// to access the underlying logger.
func (l *sLogger) toSlogS(ctx context.Context, level slog.Level, msg string,
	attrs ...any) {

	l.logger.Log(ctx, level, msg, mergeAttrs(ctx, attrs)...)
}

var _ Logger = (*sLogger)(nil)

func init() {
	// Initialise the Disabled logger.
	Disabled = NewSLogger(NewDefaultHandler(io.Discard))
	Disabled.SetLevel(btclog.LevelOff)
}
