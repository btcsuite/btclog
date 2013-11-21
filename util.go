// Copyright (c) 2013 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btclog

import (
	"github.com/conformal/seelog"
	"io"
	"strings"
)

// Disabled is a default logger that can be used to disable all logging output.
// The level must not be changed since it's not backed by a real logger.
var Disabled Logger = &SubsystemLogger{level: Off}

// LogLevelFromString returns a LogLevel given a string representation of the
// level along with a boolean that indicates if the provided string could be
// converted.
func LogLevelFromString(level string) (LogLevel, bool) {
	level = strings.ToLower(level)
	for lvl, str := range logLevelStrings {
		if level == str {
			return lvl, true
		}
	}
	return Off, false
}

// NewLoggerFromWriter creates a logger for use with non-btclog based systems.
func NewLoggerFromWriter(w io.Writer, minLevel LogLevel) (Logger, error) {
	l, err := seelog.LoggerFromWriterWithMinLevel(w, seelog.LogLevel(minLevel))
	if err != nil {
		return nil, err
	}

	return NewSubsystemLogger(l, ""), nil
}
