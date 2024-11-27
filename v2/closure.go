package btclog

import (
	"fmt"
)

// Closure is used to provide a closure over expensive logging operations so
// that they don't have to be performed when the logging level doesn't warrant
// it.
type Closure func() string

// String invokes the underlying function and returns the result.
func (c Closure) String() string {
	return c()
}

// NewClosure returns a new closure over a function that returns a string
// which itself provides a Stringer interface so that it can be used with the
// logging system.
func NewClosure(compute func() string) Closure {
	return compute
}

// Sprintf returns a fmt.Stringer that will lazily compute the string when
// needed. This is useful when the string is expensive to compute and may not be
// needed due to the log level being used.
//
// Example usage:
//
//	log.InfoS(ctx, "msg", "key", Sprintf("%.12f", 3.241))
func Sprintf(msg string, params ...any) fmt.Stringer {
	return NewClosure(func() string {
		return fmt.Sprintf(msg, params...)
	})
}
