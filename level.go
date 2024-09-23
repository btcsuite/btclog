package btclog

import (
	"fmt"
	"log/slog"
	"strings"
)

// Level is the level at which a logger is configured.  All messages sent
// to a level which is below the current level are filtered.
type Level slog.Level

// Level constants.
// Names for common levels.
const (
	LevelTrace    Level = -5
	LevelDebug          = Level(slog.LevelDebug)
	LevelInfo           = Level(slog.LevelInfo)
	LevelWarn           = Level(slog.LevelWarn)
	LevelError          = Level(slog.LevelError)
	LevelCritical Level = 9
	LevelOff      Level = 10
)

// LevelFromString returns a level based on the input string s.  If the input
// can't be interpreted as a valid log level, the info level and false is
// returned.
func LevelFromString(s string) (l Level, ok bool) {
	switch strings.ToLower(s) {
	case "trace", "trc":
		return LevelTrace, true
	case "debug", "dbg":
		return LevelDebug, true
	case "info", "inf":
		return LevelInfo, true
	case "warn", "wrn":
		return LevelWarn, true
	case "error", "err":
		return LevelError, true
	case "critical", "crt":
		return LevelCritical, true
	case "off":
		return LevelOff, true
	default:
		return LevelInfo, false
	}
}

// String returns a name for the level. If the level has a name, then that name
// in uppercase is returned. If the level is between named values, then an
// integer is appended to the uppercase name.
// Examples:
//
//	LevelWarn.String() => "WARN"
//	(LevelInfo+2).String() => "INFO+2"
func (l Level) String() string {
	if l >= LevelOff {
		return "OFF"
	}

	str := func(base string, val Level) string {
		if val == 0 {
			return base
		}
		return fmt.Sprintf("%s%+d", base, val)
	}

	switch {
	case l < LevelDebug:
		return str("TRC", l-LevelTrace)
	case l < LevelInfo:
		return str("DBG", l-LevelDebug)
	case l < LevelWarn:
		return str("INF", l-LevelInfo)
	case l < LevelError:
		return str("WRN", l-LevelWarn)
	case l < LevelCritical:
		return str("ERR", l-LevelError)
	default:
		return str("CRT", l-LevelCritical)
	}
}

type ansiColorSeq string

const (
	ansiColorSeqDarkTeal  ansiColorSeq = "38;5;30"
	ansiColorSeqDarkBlue  ansiColorSeq = "38;5;63"
	ansiColorSeqLightBlue ansiColorSeq = "38;5;86"
	ansiColorSeqYellow    ansiColorSeq = "38;5;192"
	ansiColorSeqRed       ansiColorSeq = "38;5;204"
	ansiColorSeqPink      ansiColorSeq = "38;5;134"
)

func (l Level) ansiColoSeq() ansiColorSeq {
	switch l {
	case LevelTrace:
		return ansiColorSeqDarkTeal
	case LevelDebug:
		return ansiColorSeqDarkBlue
	case LevelWarn:
		return ansiColorSeqYellow
	case LevelError:
		return ansiColorSeqRed
	case LevelCritical:
		return ansiColorSeqPink
	default:
		return ansiColorSeqLightBlue
	}
}
