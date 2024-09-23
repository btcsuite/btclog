// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btclog

import "context"

// Logger is an interface which describes a level-based logger.  A default
// implementation of Logger is implemented by this package and can be created
// by calling (*Backend).Logger.
type Logger interface {
	// Tracef formats message according to format specifier and writes to
	// to log with LevelTrace.
	Tracef(format string, params ...any)

	// Debugf formats message according to format specifier and writes to
	// log with LevelDebug.
	Debugf(format string, params ...any)

	// Infof formats message according to format specifier and writes to
	// log with LevelInfo.
	Infof(format string, params ...any)

	// Warnf formats message according to format specifier and writes to
	// to log with LevelWarn.
	Warnf(format string, params ...any)

	// Errorf formats message according to format specifier and writes to
	// to log with LevelError.
	Errorf(format string, params ...any)

	// Criticalf formats message according to format specifier and writes to
	// log with LevelCritical.
	Criticalf(format string, params ...any)

	// Trace formats message using the default formats for its operands
	// and writes to log with LevelTrace.
	Trace(v ...any)

	// Debug formats message using the default formats for its operands
	// and writes to log with LevelDebug.
	Debug(v ...any)

	// Info formats message using the default formats for its operands
	// and writes to log with LevelInfo.
	Info(v ...any)

	// Warn formats message using the default formats for its operands
	// and writes to log with LevelWarn.
	Warn(v ...any)

	// Error formats message using the default formats for its operands
	// and writes to log with LevelError.
	Error(v ...any)

	// Critical formats message using the default formats for its operands
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
	Level() Level

	// SetLevel changes the logging level to the passed level.
	SetLevel(level Level)
}
