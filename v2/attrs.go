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
	if len(value) <= 6 {
		return slog.String(key, hex.EncodeToString(value))
	}

	return slog.String(key, hex.EncodeToString(value[:6]))
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
//	ctx := log.WithCtx(ctx, "height", 1234)
//	...
//	log.Info(ctx, "Height processed") // Will contain attribute: height=1234
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
