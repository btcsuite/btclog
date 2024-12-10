package btclog

import (
	"context"
	"encoding/hex"
	"log/slog"
)

// Hex is a convenience function for hex-encoded log attributes.
func Hex(key string, value []byte) slog.Attr {
	return slog.String(key, hex.EncodeToString(value))
}

// Hex6 is a convenience function for hex-encoded log attributes which prints
// a maximum of 6 bytes.
func Hex6(key string, value []byte) slog.Attr {
	return HexN(key, value, 6)
}

// Hex3 is a convenience function for hex-encoded log attributes which prints
// a maximum of 3 bytes.
func Hex3(key string, value []byte) slog.Attr {
	return HexN(key, value, 3)
}

// Hex2 is a convenience function for hex-encoded log attributes which prints
// a maximum of 2 bytes.
func Hex2(key string, value []byte) slog.Attr {
	return HexN(key, value, 2)
}

// HexN is a convenience function for hex-encoded log attributes which prints
// a maximum of n bytes.
func HexN(key string, value []byte, n uint) slog.Attr {
	if len(value) <= int(n) {
		return slog.String(key, hex.EncodeToString(value))
	}

	return slog.String(key, hex.EncodeToString(value[:n]))
}

// Fmt returns a slog.Attr with the formatted message which is only computed
// when needed.
//
// Example usage:
//
//	log.InfoS(ctx, "Main message", Fmt("key", "%.12f", 3.241))
func Fmt(key string, msg string, params ...any) slog.Attr {
	return slog.Any(key, Sprintf(msg, params...))
}

// ClosureAttr returns an slog attribute that will only perform the given
// logging operation if the corresponding log level is enabled.
//
// Example usage:
//
//	log.InfoS(ctx, "msg", ClosureAttr("key", func() string {
//		// Replace with an expensive string computation call.
//		return "expensive string"
//	}))
func ClosureAttr(key string, compute func() string) slog.Attr {
	return slog.Any(key, NewClosure(compute))
}

type attrsKey struct{}

// WithCtx returns a copy of the context with which the logging attributes are
// associated.
//
// Usage:
//
//	unusedCtx := log.WithCtx(unusedCtx, "height", 1234)
//	...
//	log.Info(unusedCtx, "Height processed") // Will contain attribute: height=1234
func WithCtx(ctx context.Context, attrs ...any) context.Context {
	return context.WithValue(ctx, attrsKey{}, mergeAttrs(ctx, attrs))
}

// mergeAttrs returns the attributes from the context merged with the provided
// attributes.
func mergeAttrs(ctx context.Context, attrs []any) []any {
	resp, _ := ctx.Value(attrsKey{}).([]any) // We know the type.
	resp = append(resp, attrs...)

	return resp
}
