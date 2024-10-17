package btclog

import (
	"log/slog"
	"strings"

	"github.com/btcsuite/btclog"
)

// Redefine the levels here so that any package importing the original btclog
// Level does not need to import both the old and new modules.
const (
	LevelTrace    = btclog.LevelTrace
	LevelDebug    = btclog.LevelDebug
	LevelInfo     = btclog.LevelInfo
	LevelWarn     = btclog.LevelWarn
	LevelError    = btclog.LevelError
	LevelCritical = btclog.LevelCritical
	LevelOff      = btclog.LevelOff
)

// LevelFromString returns a level based on the input string s.  If the input
// can't be interpreted as a valid log level, the info level and false is
// returned.
func LevelFromString(s string) (l btclog.Level, ok bool) {
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

// slog uses some pre-defined level integers. So we will need to sometimes map
// between the btclog.Level and the slog level. The slog library defines a few
// of the commonly used levels and allows us to add a few of our own too.
const (
	levelTrace    slog.Level = -5
	levelDebug               = slog.LevelDebug
	levelInfo                = slog.LevelInfo
	levelWarn                = slog.LevelWarn
	levelError               = slog.LevelError
	levelCritical slog.Level = 9
	levelOff      slog.Level = 10
)

// toSlogLevel converts a btclog.Level to the associated slog.Level type.
func toSlogLevel(l btclog.Level) slog.Level {
	switch l {
	case LevelTrace:
		return levelTrace
	case LevelDebug:
		return levelDebug
	case LevelInfo:
		return levelInfo
	case LevelWarn:
		return levelWarn
	case LevelError:
		return levelError
	case LevelCritical:
		return levelCritical
	default:
		return levelOff
	}
}

// fromSlogLevel converts an slog.Level type to the associated btclog.Level
// type.
func fromSlogLevel(l slog.Level) btclog.Level {
	switch l {
	case levelTrace:
		return LevelTrace
	case levelDebug:
		return LevelDebug
	case levelInfo:
		return LevelInfo
	case levelWarn:
		return LevelWarn
	case levelError:
		return LevelError
	case levelCritical:
		return LevelCritical
	default:
		return LevelOff
	}
}
