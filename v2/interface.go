package btclog

import (
	"context"
	"github.com/btcsuite/btclog"
)

// Logger is an interface which describes a level-based logger. A default
// implementation of Logger is implemented by this package and can be created
// by calling (*Backend).Logger.
type Logger interface {
	// Tracef creates a formatted message from the to format specifier
	// along with any parameters then writes it to the logger with
	// LevelTrace.
	Tracef(format string, params ...any)

	// Debugf creates a formatted message from the to format specifier
	// along with any parameters then writes it to the logger with
	// LevelDebug.
	Debugf(format string, params ...any)

	// Infof creates a formatted message from the to format specifier
	// along with any parameters then writes it to the logger with
	// LevelInfo.
	Infof(format string, params ...any)

	// Warnf creates a formatted message from the to format specifier
	// along with any parameters then writes it to the logger with
	// LevelWarn.
	Warnf(format string, params ...any)

	// Errorf creates a formatted message from the to format specifier
	// along with any parameters then writes it to the logger with
	// LevelError.
	Errorf(format string, params ...any)

	// Criticalf creates a formatted message from the to format specifier
	// along with any parameters then writes it to the logger with
	// LevelCritical.
	Criticalf(format string, params ...any)

	// Trace formats a message using the default formats for its operands
	// and writes to log with LevelTrace.
	Trace(v ...any)

	// Debug formats a message using the default formats for its operands
	// and writes to log with LevelDebug.
	Debug(v ...any)

	// Info formats a message using the default formats for its operands
	// and writes to log with LevelInfo.
	Info(v ...any)

	// Warn formats a message using the default formats for its operands
	// and writes to log with LevelWarn.
	Warn(v ...any)

	// Error formats a message using the default formats for its operands
	// and writes to log with LevelError.
	Error(v ...any)

	// Critical formats a message using the default formats for its operands
	// and writes to log with LevelCritical.
	Critical(v ...any)

	// TraceS writes a structured log with the given message and key-value
	// pair attributes with LevelTrace to the log.
	TraceS(ctx context.Context, msg string, attrs ...any)

	// DebugS writes a structured log with the given message and key-value
	// pair attributes with LevelDebug to the log.
	DebugS(ctx context.Context, msg string, attrs ...any)

	// InfoS writes a structured log with the given message and key-value
	// pair attributes with LevelInfo to the log.
	InfoS(ctx context.Context, msg string, attrs ...any)

	// WarnS writes a structured log with the given message and key-value
	// pair attributes with LevelWarn to the log.
	WarnS(ctx context.Context, msg string, err error, attrs ...any)

	// ErrorS writes a structured log with the given message and key-value
	// pair attributes with LevelError to the log.
	ErrorS(ctx context.Context, msg string, err error, attrs ...any)

	// CriticalS writes a structured log with the given message and
	// key-value pair attributes with LevelCritical to the log.
	CriticalS(ctx context.Context, msg string, err error, attrs ...any)

	// Level returns the current logging level.
	Level() btclog.Level

	// SetLevel changes the logging level to the passed level.
	SetLevel(level btclog.Level)

	// SubSystem returns a copy of the logger but with the new subsystem
	// tag.
	SubSystem(tag string) Logger
}

// Ensure that the Logger implements the btclog.Logger interface so that an
// implementation of the new and expanded interface can still be used by older
// code depending on the older interface.
var _ btclog.Logger = (Logger)(nil)
