package btclog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync/atomic"

	"github.com/btcsuite/btclog"
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
	//
	// NOTE: this creates a new logger with an independent log level. This
	// means that SetLevel needs to be called on the new logger to change
	// the level as any changes to the parent logger's level after creation
	// will not be inherited by the new logger.
	SubSystem(tag string) Handler

	// WithPrefix returns a copy of the Handler but with the given string
	// prefixed to each log message. Note that the subsystem of the original
	// logger is kept but any existing prefix is overridden.
	WithPrefix(prefix string) Handler
}

// sLogger is an implementation of Logger backed by a structured sLogger.
type sLogger struct {
	level atomic.Int64

	handler Handler
	logger  *slog.Logger

	// unusedCtx is a context that will be passed to the non-structured
	// logging calls for backwards compatibility with the old v1 Logger
	// interface. Transporting a context in a struct is an anti-pattern but
	// this is purely used for backwards compatibility and to prevent
	// needing to create a fresh context for each call to the old interface
	// methods. This is ok to do since the slog package does not use this
	// context for cancellation or deadlines. It purely uses it to extract
	// any slog attributes that have been added as values to the context.
	unusedCtx context.Context
}

// NewSLogger constructs a new structured logger from the given Handler.
func NewSLogger(handler Handler) Logger {
	l := &sLogger{
		handler:   handler,
		logger:    slog.New(handler),
		unusedCtx: context.Background(),
	}
	l.level.Store(int64(toSlogLevel(handler.Level())))

	return l
}

// Tracef creates a formatted message from the to format specifier along with
// any parameters then writes it to the logger with LevelTrace.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Tracef(format string, params ...any) {
	if l.level.Load() > int64(levelTrace) {
		return
	}

	l.logger.Log(l.unusedCtx, levelTrace, fmt.Sprintf(format, params...))
}

// Debugf creates a formatted message from the to format specifier along with
// any parameters then writes it to the logger with LevelDebug.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Debugf(format string, params ...any) {
	if l.level.Load() > int64(levelDebug) {
		return
	}

	l.logger.Log(l.unusedCtx, levelDebug, fmt.Sprintf(format, params...))
}

// Infof creates a formatted message from the to format specifier along with
// any parameters then writes it to the logger with LevelInfo.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Infof(format string, params ...any) {
	if l.level.Load() > int64(levelInfo) {
		return
	}

	l.logger.Log(l.unusedCtx, levelInfo, fmt.Sprintf(format, params...))
}

// Warnf creates a formatted message from the to format specifier along with
// any parameters then writes it to the logger with LevelWarn.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Warnf(format string, params ...any) {
	if l.level.Load() > int64(levelWarn) {
		return
	}

	l.logger.Log(l.unusedCtx, levelWarn, fmt.Sprintf(format, params...))
}

// Errorf creates a formatted message from the to format specifier along with
// any parameters then writes it to the logger with LevelError.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Errorf(format string, params ...any) {
	if l.level.Load() > int64(levelError) {
		return
	}

	l.logger.Log(l.unusedCtx, levelError, fmt.Sprintf(format, params...))
}

// Criticalf creates a formatted message from the to format specifier along
// with any parameters then writes it to the logger with LevelCritical.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Criticalf(format string, params ...any) {
	if l.level.Load() > int64(levelCritical) {
		return
	}

	l.logger.Log(l.unusedCtx, levelCritical, fmt.Sprintf(format, params...))
}

// Trace formats a message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelTrace.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Trace(v ...any) {
	if l.level.Load() > int64(levelTrace) {
		return
	}

	l.logger.Log(l.unusedCtx, levelTrace, fmt.Sprint(v...))
}

// Debug formats a message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelDebug.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Debug(v ...any) {
	if l.level.Load() > int64(levelDebug) {
		return
	}

	l.logger.Log(l.unusedCtx, levelDebug, fmt.Sprint(v...))
}

// Info formats a message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelInfo.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Info(v ...any) {
	if l.level.Load() > int64(levelInfo) {
		return
	}

	l.logger.Log(l.unusedCtx, levelInfo, fmt.Sprint(v...))
}

// Warn formats a message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelWarn.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Warn(v ...any) {
	if l.level.Load() > int64(levelWarn) {
		return
	}

	l.logger.Log(l.unusedCtx, levelWarn, fmt.Sprint(v...))
}

// Error formats a message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelError.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Error(v ...any) {
	if l.level.Load() > int64(levelError) {
		return
	}

	l.logger.Log(l.unusedCtx, levelError, fmt.Sprint(v...))
}

// Critical formats a message using the default formats for its operands,
// prepends the prefix as necessary, and writes to log with LevelCritical.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Critical(v ...any) {
	if l.level.Load() > int64(levelCritical) {
		return
	}

	l.logger.Log(l.unusedCtx, levelCritical, fmt.Sprint(v...))
}

// TraceS writes a structured log with the given message and key-value pair
// attributes with LevelTrace to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) TraceS(ctx context.Context, msg string, attrs ...any) {
	if l.level.Load() > int64(levelTrace) {
		return
	}

	l.logger.Log(ctx, levelTrace, msg, mergeAttrs(ctx, attrs)...)
}

// DebugS writes a structured log with the given message and key-value pair
// attributes with LevelDebug to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) DebugS(ctx context.Context, msg string, attrs ...any) {
	if l.level.Load() > int64(levelDebug) {
		return
	}

	l.logger.Log(ctx, levelDebug, msg, mergeAttrs(ctx, attrs)...)
}

// InfoS writes a structured log with the given message and key-value pair
// attributes with LevelInfo to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) InfoS(ctx context.Context, msg string, attrs ...any) {
	if l.level.Load() > int64(levelInfo) {
		return
	}

	l.logger.Log(ctx, levelInfo, msg, mergeAttrs(ctx, attrs)...)
}

// WarnS writes a structured log with the given message and key-value pair
// attributes with LevelWarn to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) WarnS(ctx context.Context, msg string, err error,
	attrs ...any) {

	if l.level.Load() > int64(levelWarn) {
		return
	}

	if err != nil {
		attrs = append([]any{slog.String("err", err.Error())}, attrs...)
	}

	l.logger.Log(ctx, levelWarn, msg, mergeAttrs(ctx, attrs)...)
}

// ErrorS writes a structured log with the given message and key-value pair
// attributes with LevelError to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) ErrorS(ctx context.Context, msg string, err error,
	attrs ...any) {

	if l.level.Load() > int64(levelError) {
		return
	}

	if err != nil {
		attrs = append([]any{slog.String("err", err.Error())}, attrs...)
	}

	l.logger.Log(ctx, levelError, msg, mergeAttrs(ctx, attrs)...)
}

// CriticalS writes a structured log with the given message and key-value pair
// attributes with LevelCritical to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) CriticalS(ctx context.Context, msg string, err error,
	attrs ...any) {

	if l.level.Load() > int64(levelCritical) {
		return
	}

	if err != nil {
		attrs = append([]any{slog.String("err", err.Error())}, attrs...)
	}

	l.logger.Log(ctx, levelCritical, msg, mergeAttrs(ctx, attrs)...)
}

// Level returns the current logging level of the Handler.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Level() btclog.Level {
	return fromSlogLevel(slog.Level(l.level.Load()))
}

// SetLevel changes the logging level of the Handler to the passed level.
//
// This is part of the Logger interface implementation.
func (l *sLogger) SetLevel(level btclog.Level) {
	l.level.Store(int64(toSlogLevel(level)))
	l.handler.SetLevel(level)
}

// SubSystem returns a copy of the logger but with the new subsystem tag.
//
// This is part of the Logger interface implementation.
func (l *sLogger) SubSystem(tag string) Logger {
	return NewSLogger(l.handler.SubSystem(tag))
}

// WithPrefix returns a copy of the logger but with the given string prefixed to
// each log message. Note that the subsystem of the original logger is kept but
// any existing prefix is overridden.
//
// This is part of the Logger interface implementation.
func (l *sLogger) WithPrefix(prefix string) Logger {
	return NewSLogger(l.handler.WithPrefix(prefix))
}

var _ Logger = (*sLogger)(nil)

func init() {
	// Initialise the Disabled logger.
	Disabled = NewSLogger(NewDefaultHandler(io.Discard))
	Disabled.SetLevel(btclog.LevelOff)
}
