// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2009 The Go Authors.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
// * Redistributions of source code must retain the above copyright
//   notice, this list of conditions and the following disclaimer.
// * Redistributions in binary form must reproduce the above
//   copyright notice, this list of conditions and the following disclaimer
//   in the documentation and/or other materials provided with the
//   distribution.
// * Neither the name of Google LLC nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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
