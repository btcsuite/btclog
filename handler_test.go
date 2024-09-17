package btclog

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"
)

// TestDefaultHandler tests that the DefaultHandler's output looks as expected.
func TestDefaultHandler(t *testing.T) {
	t.Parallel()

	timeSource := func() time.Time {
		return time.Unix(100, 100)
	}

	tests := []struct {
		name               string
		handlerConstructor func(w io.Writer) Handler
		logFunc            func(log Logger)
		expectedLog        string
	}{
		{
			name: "Basic calls and levels",
			handlerConstructor: func(w io.Writer) Handler {
				h := NewDefaultHandler(
					w, WithTimeSource(timeSource),
				)
				h.SetLevel(LevelDebug)
				return h
			},
			logFunc: func(log Logger) {
				log.Info("Test Basic Log")
				log.Debugf("Test basic log with %s", "format")
				log.Trace("Log should not appear due to level")
			},
			expectedLog: `1970-01-01 02:01:40.000 [INF]: Test Basic Log
1970-01-01 02:01:40.000 [DBG]: Test basic log with format
`,
		},
		{
			name: "Sub-system tag",
			handlerConstructor: func(w io.Writer) Handler {
				h := NewDefaultHandler(w, WithNoTimestamp())
				return h.SubSystem("SUBS")
			},
			logFunc: func(log Logger) {
				log.Info("Test Basic Log")
			},
			expectedLog: `[INF] SUBS: Test Basic Log
`,
		},
		{
			name: "Test all levels",
			handlerConstructor: func(w io.Writer) Handler {
				h := NewDefaultHandler(w, WithNoTimestamp())
				h.SetLevel(LevelTrace)
				return h
			},
			logFunc: func(log Logger) {
				log.Trace("Trace")
				log.Debug("Debug")
				log.Info("Info")
				log.Warn("Warn")
				log.Error("Error")
				log.Critical("Critical")
			},
			expectedLog: `[TRC]: Trace
[DBG]: Debug
[INF]: Info
[WRN]: Warn
[ERR]: Error
[CRT]: Critical
`,
		},
		{
			name: "Structured Logs",
			handlerConstructor: func(w io.Writer) Handler {
				return NewDefaultHandler(w, WithNoTimestamp())
			},
			logFunc: func(log Logger) {
				ctx := context.Background()
				log.InfoS(ctx, "No attributes")
				log.InfoS(ctx, "Single word attribute", "key", "value")
				log.InfoS(ctx, "Multi word string value", "key with spaces", "value")
				log.InfoS(ctx, "Number attribute", "key", 5)
				log.InfoS(ctx, "Bad key", "key")

				type b struct {
					name    string
					age     int
					address *string
				}

				var c *b
				log.InfoS(ctx, "Nil pointer value", "key", c)

				c = &b{name: "Bob", age: 5}
				log.InfoS(ctx, "Struct values", "key", c)

				ctx = WithCtx(ctx, "request_id", 5, "user_name", "alice")
				log.InfoS(ctx, "Test context attributes", "key", "value")
			},
			expectedLog: `[INF]: No attributes
[INF]: Single word attribute key=value
[INF]: Multi word string value "key with spaces"=value
[INF]: Number attribute key=5
[INF]: Bad key !BADKEY=key
[INF]: Nil pointer value key=<nil>
[INF]: Struct values key="&{name:Bob age:5 address:<nil>}"
[INF]: Test context attributes request_id=5 user_name=alice key=value
`,
		},
		{
			name: "Error logs",
			handlerConstructor: func(w io.Writer) Handler {
				return NewDefaultHandler(w, WithNoTimestamp())
			},
			logFunc: func(log Logger) {
				log.Error("Error string")
				log.Errorf("Error formatted string")

				ctx := context.Background()
				log.ErrorS(ctx, "Structured error log with nil error", nil)
				log.ErrorS(ctx, "Structured error with non-nil error", errors.New("oh no"))
				log.ErrorS(ctx, "Structured error with attributes", errors.New("oh no"), "key", "value")
			},
			expectedLog: `[ERR]: Error string
[ERR]: Error formatted string
[ERR]: Structured error log with nil error
[ERR]: Structured error with non-nil error err="oh no"
[ERR]: Structured error with attributes err="oh no" key=value
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := test.handlerConstructor(&buf)

			test.logFunc(NewSLogger(handler))

			if string(buf.Bytes()) != test.expectedLog {
				t.Fatalf("Log result mismatch. Expected \n\"%s\", got \n\"%s\"", test.expectedLog, buf.Bytes())
			}
		})
	}
}
