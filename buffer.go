// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted from go/src/log/slog/internal/buffer.go

package btclog

import "sync"

type buffer []byte

// bufferPool defines a concurrent safe free list of byte slices used to provide
// temporary buffers for formatting log messages prior to outputting them.
var bufferPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return (*buffer)(&b)
	},
}

// newBuffer returns a byte slice from the free list. A new buffer is allocated
// if there are not any available on the free list. The returned byte slice
// should be returned to the fee list by using the recycleBuffer function when
// the caller is done with it.
func newBuffer() *buffer {
	return bufferPool.Get().(*buffer)
}

// free puts the provided byte slice, which should have been obtained via the
// newBuffer function, back on the free list.
func (b *buffer) free() {
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufferPool.Put(b)
	}
}

func (b *buffer) writeByte(p byte) {
	*b = append(*b, p)
}

func (b *buffer) writeBytes(p []byte) {
	*b = append(*b, p...)
}

func (b *buffer) writeString(s string) {
	*b = append(*b, s...)
}
