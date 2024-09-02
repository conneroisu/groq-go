// Package test contains test helpers.
package test

import "errors"

var (
	// ErrTestErrorAccumulatorWriteFailed is the error returned by the failing
	// error buffer.
	ErrTestErrorAccumulatorWriteFailed = errors.New("test error accumulator failed")
)

// FailingErrorBuffer is a buffer that always fails to write.
type FailingErrorBuffer struct{}

// Write always fails.
//
// It is used to test the behavior of the error accumulator.
func (b *FailingErrorBuffer) Write(_ []byte) (n int, err error) {
	return 0, ErrTestErrorAccumulatorWriteFailed
}

// Len always returns 0.
func (b *FailingErrorBuffer) Len() int {
	return 0
}

// Bytes always
func (b *FailingErrorBuffer) Bytes() []byte {
	return []byte{}
}
