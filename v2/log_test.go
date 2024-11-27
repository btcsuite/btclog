package btclog

import (
	"context"
	"io"
	"testing"

	"github.com/btcsuite/btclog"
)

type complexType struct {
	m map[string]string
	s []string
}

var testType = complexType{
	m: map[string]string{
		"key1": "value1",
		"key2": "value2",
	},
	s: []string{"a", "b", "c"},
}

// BenchmarkLogger benchmarks the performance of the default v2 logger for each
// logging level. This helps evaluate the effect of any change to the v2 logger.
func BenchmarkLogger(b *testing.B) {
	ctx := context.Background()
	log := NewSLogger(NewDefaultHandler(io.Discard))

	// Set the level to Info so that Debug logs are skipped.
	log.SetLevel(LevelInfo)

	tests := []struct {
		name    string
		logFunc func()
	}{
		{
			name: "Skipped Simple `f` Log",
			logFunc: func() {
				log.Debugf("msg")
			},
		},
		{
			name: "Skipped Simple `S` Log",
			logFunc: func() {
				log.DebugS(ctx, "msg")
			},
		},
		{
			name: "Simple `f` Log",
			logFunc: func() {
				log.Infof("msg")
			},
		},
		{
			name: "Simple `S` Log",
			logFunc: func() {
				log.InfoS(ctx, "msg")
			},
		},
		{
			name: "Skipped Complex `f` Log",
			logFunc: func() {
				log.Debugf("Debugf "+
					"request_id=%d, "+
					"complex_type=%v,  "+
					"type=%T, "+
					"floating_point=%.12f, "+
					"hex_value=%x, "+
					"fmt_string=%s",
					5, testType, testType,
					3.141592653589793, []byte{0x01, 0x02},
					Sprintf("%s, %v, %T, %.12f",
						"string", testType, testType,
						3.141592653589793))
			},
		},
		{
			name: "Skipped Complex `S` Log",
			logFunc: func() {
				log.DebugS(ctx, "InfoS",
					"request_id", 5,
					"complex_type", testType,
					Fmt("type", "%T", testType),
					Fmt("floating_point", "%.12f", 3.141592653589793),
					Hex("hex_value", []byte{0x01, 0x02}),
					Fmt("fmt_string", "%s, %v, %T, %.12f",
						"string", testType, testType,
						3.141592653589793))
			},
		},
		{
			name: "Complex `f` Log",
			logFunc: func() {
				log.Infof("Infof "+
					"request_id=%d, "+
					"complex_type=%v,  "+
					"type=%T, "+
					"floating_point=%.12f, "+
					"hex_value=%x, "+
					"fmt_string=%s",
					5, testType, testType,
					3.141592653589793, []byte{0x01, 0x02},
					Sprintf("%s, %v, %T, %.12f",
						"string", testType, testType,
						3.141592653589793))
			},
		},
		{
			name: "Complex `S` Log",
			logFunc: func() {
				log.InfoS(ctx, "InfoS",
					"request_id", 5,
					"complex_type", testType,
					Fmt("type", "%T", testType),
					Fmt("floating_point", "%.12f", 3.141592653589793),
					Hex("hex_value", []byte{0x01, 0x02}),
					Fmt("fmt_string", "%s, %v, %T, %.12f",
						"string", testType, testType,
						3.141592653589793))
			},
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				test.logFunc()
			}
		})
	}
}

// BenchmarkV1vsV2Loggers compares the performance of the btclog V1 logger to
// the btclog V2 logger in various logs that can be used with both the legacy
// Logger interface along with the new Logger interface. This, in other words,
// therefore compares the performance change when the V1 logger is swapped out
// for the V2 logger.
func BenchmarkV1vsV2Loggers(b *testing.B) {
	loggers := []struct {
		name string
		btclog.Logger
	}{
		{
			name:   "btclog v1",
			Logger: btclog.NewBackend(io.Discard).Logger("V1"),
		},
		{
			name:   "btclog v2",
			Logger: NewSLogger(NewDefaultHandler(io.Discard)),
		},
	}

	for _, logger := range loggers {
		// Test a basic message log with no formatted strings.
		b.Run(logger.name+" simple", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				logger.Infof("Basic")
			}
		})

		// Test a basic message log with no formatted strings that gets
		// skipped due to the current log level.
		b.Run(logger.name+" skipped simple", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				logger.Debugf("Basic")
			}
		})

		// Test a log line with a variety of different types and
		// formats.
		b.Run(logger.name+" complex", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				logger.Infof("Infof "+
					"request_id=%d, "+
					"complex_type=%v,  "+
					"type=%T, "+
					"floating_point=%.12f, "+
					"hex_value=%x, "+
					"fmt_string=%s",
					5, testType, testType,
					3.141592653589793, []byte{0x01, 0x02},
					Sprintf("%s, %v, %T, %.12f",
						"string", testType, testType,
						3.141592653589793))
			}
		})

		// Test a log line with a variety of different types and formats
		// that then gets skipped due to the current log level.
		b.Run(logger.name+" skipped complex", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				logger.Debugf("Infof "+
					"request_id=%d, "+
					"complex_type=%v,  "+
					"type=%T, "+
					"floating_point=%.12f, "+
					"hex_value=%x, "+
					"fmt_string=%s",
					5, testType, testType,
					3.141592653589793, []byte{0x01, 0x02},
					Sprintf("%s, %v, %T, %.12f",
						"string", testType, testType,
						3.141592653589793))
			}
		})
	}
}
