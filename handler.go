package btclog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"unicode"
	"unicode/utf8"
)

// HandlerOption is the signature of a functional option that can be used to
// modify the behaviour of the DefaultHandler.
type HandlerOption func(*handlerOpts)

// handlerOpts holds options that can be modified by a HandlerOption.
type handlerOpts struct {
	flag          uint32
	withTimestamp bool
	styled        bool
	timeSource    func() time.Time
}

// defaultHandlerOpts constructs a handlerOpts with default settings.
func defaultHandlerOpts() *handlerOpts {
	return &handlerOpts{
		flag:          defaultFlags,
		withTimestamp: true,
		styled:        false,
	}
}

// WithCallerFlags can be used to overwrite the default caller flag option.
func WithCallerFlags(flags uint32) HandlerOption {
	return func(opts *handlerOpts) {
		opts.flag = flags
	}
}

// WithTimeSource can be used to overwrite the time sourced from the slog
// Record.
func WithTimeSource(fn func() time.Time) HandlerOption {
	return func(opts *handlerOpts) {
		opts.timeSource = fn
	}
}

// WithNoTimestamp is an option that can be used to omit timestamps from the log
// lines.
func WithNoTimestamp() HandlerOption {
	return func(opts *handlerOpts) {
		opts.withTimestamp = false
	}
}

// DefaultHandler is a Handler that can be used along with NewSLogger to
// instantiate a structured logger.
type DefaultHandler struct {
	opts *handlerOpts

	level           int64
	tag             string
	fields          []slog.Attr
	callstackOffset bool

	flag uint32
	buf  *buffer
	mu   *sync.Mutex
	w    io.Writer
}

// A compile-time check to ensure that DefaultHandler implements Handler.
var _ Handler = (*DefaultHandler)(nil)

// Level returns the current logging level of the Handler.
//
// NOTE: This is part of the Handler interface.
func (d *DefaultHandler) Level() Level {
	return Level(atomic.LoadInt64(&d.level))
}

// SetLevel changes the logging level of the Handler to the passed
// level.
//
// NOTE: This is part of the Handler interface.
func (d *DefaultHandler) SetLevel(level Level) {
	atomic.StoreInt64(&d.level, int64(level))
}

// NewDefaultHandler creates a new Handler that can be used along with
// NewSLogger to instantiate a structured logger.
func NewDefaultHandler(w io.Writer, options ...HandlerOption) *DefaultHandler {
	opts := defaultHandlerOpts()
	for _, o := range options {
		o(opts)
	}

	return &DefaultHandler{
		w:     w,
		level: int64(LevelInfo),
		opts:  opts,
		buf:   newBuffer(),
		mu:    &sync.Mutex{},
	}
}

// Enabled reports whether the handler handles records at the given level.
//
// NOTE: this is part of the slog.Handler interface.
func (d *DefaultHandler) Enabled(_ context.Context, level slog.Level) bool {
	return atomic.LoadInt64(&d.level) <= int64(level)
}

// Handle handles the Record. It will only be called if Enabled returns true.
//
// NOTE: this is part of the slog.Handler interface.
func (d *DefaultHandler) Handle(_ context.Context, r slog.Record) error {
	buf := newBuffer()
	defer buf.free()

	// Timestamp.
	if d.opts.withTimestamp {
		// First check if the options provided specified a different
		// time source to use. Otherwise, use the provided record time.
		if d.opts.timeSource != nil {
			writeTimestamp(buf, d.opts.timeSource())
		} else if !r.Time.IsZero() {
			writeTimestamp(buf, r.Time)
		}
	}

	// Level.
	d.writeLevel(buf, Level(r.Level))

	// Sub-system tag.
	if d.tag != "" {
		buf.writeString(" " + d.tag)
	}

	// The call-site.
	if d.opts.flag&(Lshortfile|Llongfile) != 0 {
		skip := 6
		if d.callstackOffset {
			skip = 4
		}
		file, line := callsite(d.opts.flag, skip)
		d.writeCallSite(buf, file, line)
	}

	// Finish off the header.
	buf.writeString(": ")

	// Write the log message itself.
	if r.Message != "" {
		buf.writeString(r.Message)
	}

	// Append logger fields.
	for _, attr := range d.fields {
		d.appendAttr(buf, attr)
	}

	// Append slog attributes
	r.Attrs(func(a slog.Attr) bool {
		d.appendAttr(buf, a)
		return true
	})
	buf.writeByte('\n')

	d.mu.Lock()
	defer d.mu.Unlock()
	_, err := d.w.Write(*buf)

	return err
}

// WithAttrs returns a new Handler with the given attributes added.
//
// NOTE: this is part of the slog.Handler interface.
func (d *DefaultHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return d.with(d.tag, true, attrs...)
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups. All this implementation does is add to the
// existing tag used for the logger.
//
// NOTE: this is part of the slog.Handler interface.
func (d *DefaultHandler) WithGroup(name string) slog.Handler {
	if d.tag != "" {
		name = d.tag + "." + name
	}
	return d.with(name, true)
}

// SubSystem returns a copy of the given handler but with the new tag. All
// attributes added with WithAttrs will be kept but all groups added with
// WithGroup are lost.
//
// NOTE: this is part of the Handler interface.
func (d *DefaultHandler) SubSystem(tag string) Handler {
	return d.with(tag, false)
}

// with returns a new logger with the given attributes added.
// withCallstackOffset should be false if the caller returns a concrete
// DefaultHandler and true if the caller returns the Handler interface.
func (d *DefaultHandler) with(tag string, withCallstackOffset bool,
	attrs ...slog.Attr) *DefaultHandler {

	d.mu.Lock()
	sl := *d
	d.mu.Unlock()
	sl.buf = newBuffer()

	sl.mu = &sync.Mutex{}
	sl.fields = append(
		make([]slog.Attr, 0, len(d.fields)+len(attrs)), d.fields...,
	)
	sl.fields = append(sl.fields, attrs...)
	sl.callstackOffset = withCallstackOffset
	sl.tag = tag

	return &sl
}

func (d *DefaultHandler) appendAttr(buf *buffer, a slog.Attr) {
	// Resolve the Attr's value before doing anything else.
	a.Value = a.Value.Resolve()

	// Ignore empty Attrs.
	if a.Equal(slog.Attr{}) {
		return
	}

	d.appendKey(buf, a.Key)
	appendValue(buf, a.Value)
}

func (d *DefaultHandler) writeLevel(buf *buffer, level Level) {
	str := fmt.Sprintf("[%s]", level)

	buf.writeString(str)
}

func (d *DefaultHandler) writeCallSite(buf *buffer, file string, line int) {
	if file == "" {
		return
	}

	buf.writeString(fmt.Sprintf(" %s:%d", file, line))
}

func appendString(buf *buffer, str string) {
	if needsQuoting(str) {
		*buf = strconv.AppendQuote(*buf, str)
	} else {
		buf.writeString(str)
	}
}

func (d *DefaultHandler) appendKey(buf *buffer, key string) {
	buf.writeString(" ")
	if needsQuoting(key) {
		key = strconv.Quote(key)
	}
	key += "="

	buf.writeString(key)
}

func appendValue(buf *buffer, v slog.Value) {
	defer func() {
		// Recovery in case of nil pointer dereferences.
		if r := recover(); r != nil {
			// Catch any panics that are most likely due to nil
			// pointers.
			appendString(buf, fmt.Sprintf("!PANIC: %v", r))
		}
	}()

	appendTextValue(buf, v)
}

func appendTextValue(buf *buffer, v slog.Value) {
	switch v.Kind() {
	case slog.KindString:
		appendString(buf, v.String())
	case slog.KindAny:
		appendString(buf, fmt.Sprintf("%+v", v.Any()))
	default:
		appendString(buf, fmt.Sprintf("%s", v))
	}
}

// Copied from log/slog/text_handler.go.
//
// needsQuoting returns true if the given strings should be wrapped in quotes.
func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			// Quote anything except a backslash that would need
			// quoting in a JSON string, as well as space and '='.
			if b != '\\' && (b == ' ' || b == '=' || !safeSet[b]) {
				return true
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError || unicode.IsSpace(r) ||
			!unicode.IsPrint(r) {

			return true
		}
		i += size
	}
	return false
}

// Copied from encoding/json/tables.go.
//
// safeSet holds the value true if the ASCII character with the given array
// position can be represented inside a JSON string without any further
// escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), and the backslash character ("\").
var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}
