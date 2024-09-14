//go:build !test
// +build !test

package test

import (
	"errors"
	"testing"
)

// TestErrTestErrorAccumulatorWriteFailed_Error tests the Error method of ErrTestErrorAccumulatorWriteFailed.
func TestErrTestErrorAccumulatorWriteFailed_Error(t *testing.T) {
	err := ErrTestErrorAccumulatorWriteFailed{}
	expected := "test error accumulator failed"

	if err.Error() != expected {
		t.Errorf("Error() returned %q, expected %q", err.Error(), expected)
	}
}

// TestFailingErrorBuffer_Write tests the Write method of FailingErrorBuffer with various inputs.
func TestFailingErrorBuffer_Write(t *testing.T) {
	buf := &FailingErrorBuffer{}

	testCases := []struct {
		name  string
		input []byte
	}{
		{"nil slice", nil},
		{"empty slice", []byte{}},
		{"non-empty slice", []byte("test data")},
		{"large slice", make([]byte, 1000)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, err := buf.Write(tc.input)
			if n != 0 {
				t.Errorf("Write(%q) returned n=%d, expected n=0", tc.input, n)
			}
			if !errors.Is(err, ErrTestErrorAccumulatorWriteFailed{}) {
				t.Errorf("Write(%q) returned err=%v, expected ErrTestErrorAccumulatorWriteFailed{}", tc.input, err)
			}
		})
	}
}

// TestFailingErrorBuffer_Len tests the Len method of FailingErrorBuffer.
func TestFailingErrorBuffer_Len(t *testing.T) {
	buf := &FailingErrorBuffer{}

	length := buf.Len()
	if length != 0 {
		t.Errorf("Len() returned %d, expected 0", length)
	}

	// After Write calls
	_, _ = buf.Write([]byte("test"))
	length = buf.Len()
	if length != 0 {
		t.Errorf("Len() after Write returned %d, expected 0", length)
	}
}

// TestFailingErrorBuffer_Bytes tests the Bytes method of FailingErrorBuffer.
func TestFailingErrorBuffer_Bytes(t *testing.T) {
	buf := &FailingErrorBuffer{}

	bytes := buf.Bytes()
	if len(bytes) != 0 {
		t.Errorf("Bytes() returned %v (len=%d), expected empty byte slice", bytes, len(bytes))
	}

	// After Write calls
	_, _ = buf.Write([]byte("test"))
	bytes = buf.Bytes()
	if len(bytes) != 0 {
		t.Errorf("Bytes() after Write returned %v (len=%d), expected empty byte slice", bytes, len(bytes))
	}
}

// TestFailingErrorBuffer_MultipleWrites tests multiple Write calls to FailingErrorBuffer.
func TestFailingErrorBuffer_MultipleWrites(t *testing.T) {
	buf := &FailingErrorBuffer{}

	for i := 0; i < 5; i++ {
		n, err := buf.Write([]byte("data"))
		if n != 0 {
			t.Errorf("Write call %d returned n=%d, expected n=0", i+1, n)
		}
		if !errors.Is(err, ErrTestErrorAccumulatorWriteFailed{}) {
			t.Errorf("Write call %d returned err=%v, expected ErrTestErrorAccumulatorWriteFailed{}", i+1, err)
		}
	}

	if buf.Len() != 0 {
		t.Errorf("Len() after multiple Writes returned %d, expected 0", buf.Len())
	}

	if len(buf.Bytes()) != 0 {
		t.Errorf("Bytes() after multiple Writes returned len=%d, expected 0", len(buf.Bytes()))
	}
}

// TestFailingErrorBuffer_InterfaceCompliance checks if FailingErrorBuffer implements necessary interfaces.
func TestFailingErrorBuffer_InterfaceCompliance(t *testing.T) {
	var _ error = ErrTestErrorAccumulatorWriteFailed{}
	var _ interface{ Write([]byte) (int, error) } = &FailingErrorBuffer{}
	var _ interface{ Len() int } = &FailingErrorBuffer{}
	var _ interface{ Bytes() []byte } = &FailingErrorBuffer{}
}
