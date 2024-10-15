package btclog

import "github.com/btcsuite/btclog"

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

	// Level returns the current logging level.
	Level() btclog.Level

	// SetLevel changes the logging level to the passed level.
	SetLevel(level btclog.Level)
}

// Ensure that the Logger implements the btclog.Logger interface so that an
// implementation of the new and expanded interface can still be used by older
// code depending on the older interface.
var _ btclog.Logger = (Logger)(nil)
