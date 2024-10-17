package btclog

import (
	"os"
	"runtime"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// defaultFlags specifies changes to the default logger behavior.  It is set
// during package init and configured using the LOGFLAGS environment variable.
// New logger backends can override these default flags using WithFlags.
var defaultFlags uint32

// Flags to modify Backend's behavior.
const (
	// Llongfile modifies the logger output to include full path and line number
	// of the logging callsite, e.g. /a/b/c/main.go:123.
	Llongfile uint32 = 1 << iota

	// Lshortfile modifies the logger output to include filename and line number
	// of the logging callsite, e.g. main.go:123.  Overrides Llongfile.
	Lshortfile
)

// From stdlib log package.
// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid
// zero-padding.
func itoa(buf *buffer, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	buf.writeBytes(b[bp:])
}

// writeTimestamp writes the date in the format 'YYYY-MM-DD hh:mm:ss.sss' to the
// buffer.
func writeTimestamp(buf *buffer, t time.Time) {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	ms := t.Nanosecond() / 1e6

	itoa(buf, year, 4)
	buf.writeByte('-')
	itoa(buf, int(month), 2)
	buf.writeByte('-')
	itoa(buf, day, 2)
	buf.writeByte(' ')
	itoa(buf, hour, 2)
	buf.writeByte(':')
	itoa(buf, min, 2)
	buf.writeByte(':')
	itoa(buf, sec, 2)
	buf.writeByte('.')
	itoa(buf, ms, 3)
	buf.writeByte(' ')
}

// callsite returns the file name and line number of the callsite to the
// subsystem logger.
func callsite(flag uint32, skipDepth int) (string, int) {
	_, file, line, ok := runtime.Caller(skipDepth)
	if !ok {
		return "???", 0
	}
	if flag&Lshortfile != 0 {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if os.IsPathSeparator(file[i]) {
				short = file[i+1:]
				break
			}
		}
		file = short
	}
	return file, line
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

func init() {
	// Read logger flags from the LOGFLAGS environment variable. Multiple
	// flags can be set at once, separated by commas.
	for _, f := range strings.Split(os.Getenv("LOGFLAGS"), ",") {
		switch f {
		case "longfile":
			defaultFlags |= Llongfile
		case "shortfile":
			defaultFlags |= Lshortfile
		}
	}
}
