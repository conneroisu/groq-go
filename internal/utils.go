package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// ErrorAccumulator is an interface for accumulating errors
type ErrorAccumulator interface {
	Write(p []byte) error
	Bytes() []byte
}

var (
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty messages")
)

// errorBuffer is an interface for error buffer
type errorBuffer interface {
	io.Writer
	Len() int
	Bytes() []byte
}

// DefaultErrorAccumulator is the default error accumulator
type DefaultErrorAccumulator struct {
	Buffer errorBuffer
}

// NewErrorAccumulator creates a new error accumulator
func NewErrorAccumulator() ErrorAccumulator {
	return &DefaultErrorAccumulator{
		Buffer: &bytes.Buffer{},
	}
}

// Write writes bytes to the error accumulator
func (e *DefaultErrorAccumulator) Write(p []byte) error {
	_, err := e.Buffer.Write(p)
	if err != nil {
		return fmt.Errorf("error accumulator write error, %w", err)
	}
	return nil
}

// Bytes returns the accumulated error bytes
func (e *DefaultErrorAccumulator) Bytes() (errBytes []byte) {
	if e.Buffer.Len() == 0 {
		return
	}
	errBytes = e.Buffer.Bytes()
	return
}

// Marshaller is an interface for marshalling data
type Marshaller interface {
	Marshal(value any) ([]byte, error)
}

// JSONMarshaller is a marshaller for JSON data
type JSONMarshaller struct{}

// Marshal marshals data
func (jm *JSONMarshaller) Marshal(value any) ([]byte, error) {
	return json.Marshal(value)
}

// Unmarshaler is an interface for unmarshalling data
type Unmarshaler interface {
	Unmarshal(data []byte, v any) error
}

// JSONUnmarshaler is an unmarshaller for JSON data
type JSONUnmarshaler struct{}

// Unmarshal unmarshals data
func (jm *JSONUnmarshaler) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
