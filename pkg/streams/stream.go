package streams

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/groqerr"
)

type (
	// Streamer is an interface for a Streamer.
	Streamer[T any] interface {
		*T
	}
	// DefaultErrorAccumulator is a default implementation of ErrorAccumulator
	DefaultErrorAccumulator struct {
		Buffer groqerr.ErrorBuffer
	}
	// StreamReader is a stream reader.
	StreamReader[T any] struct {
		emptyMessagesLimit uint
		isFinished         bool
		Reader             *bufio.Reader
		response           *http.Response
		ErrAccumulator     ErrorAccumulator
		Header             http.Header // Header is the header of the response.
	}
	// ErrorAccumulator is an interface for a unit that accumulates errors.
	ErrorAccumulator interface {
		// Write method writes bytes to the error accumulator
		//
		// It implements the io.Writer interface.
		Write(p []byte) error
		// Bytes method returns the bytes of the error accumulator.
		Bytes() []byte
	}
)

// Recv receives a response from the stream.
func (stream *StreamReader[T]) Recv() (response T, err error) {
	if stream.isFinished {
		err = io.EOF
		return response, err
	}
	return stream.processLines()
}

// processLines processes the lines of the current response in the stream.
func (stream *StreamReader[T]) processLines() (T, error) {
	var (
		headerData         = []byte("data: ")
		errorPrefix        = []byte(`data: {"error":`)
		emptyMessagesCount uint
		hasErrorPrefix     bool
	)
	for {
		rawLine, err := stream.Reader.ReadBytes('\n')
		if err != nil || hasErrorPrefix {
			respErr := stream.unmarshalError()
			if respErr != nil {
				return *new(T),
					fmt.Errorf("error, %w", respErr.Error)
			}
			return *new(T), err
		}
		noSpaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasPrefix(noSpaceLine, errorPrefix) {
			hasErrorPrefix = true
		}
		if !bytes.HasPrefix(noSpaceLine, headerData) || hasErrorPrefix {
			if hasErrorPrefix {
				noSpaceLine = bytes.TrimPrefix(noSpaceLine, headerData)
			}
			err := stream.ErrAccumulator.Write(noSpaceLine)
			if err != nil {
				return *new(T), err
			}
			emptyMessagesCount++
			if emptyMessagesCount > stream.emptyMessagesLimit {
				return *new(T), groqerr.ErrTooManyEmptyStreamMessages{}
			}
			continue
		}
		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)
		if string(noPrefixLine) == "[DONE]" {
			stream.isFinished = true
			return *new(T), io.EOF
		}
		var response T
		unmarshalErr := json.Unmarshal(noPrefixLine, &response)
		if unmarshalErr != nil {
			return *new(T), unmarshalErr
		}
		return response, nil
	}
}
func (stream *StreamReader[T]) unmarshalError() (errResp *groqerr.ErrorResponse) {
	errBytes := stream.ErrAccumulator.Bytes()
	if len(errBytes) == 0 {
		return
	}
	err := json.Unmarshal(errBytes, &errResp)
	if err != nil {
		errResp = nil
	}
	return
}

// Close closes the stream.
func (stream *StreamReader[T]) Close() error {
	return stream.response.Body.Close()
}

// NewErrorAccumulator creates a new error accumulator
func NewErrorAccumulator() ErrorAccumulator {
	return &DefaultErrorAccumulator{
		Buffer: &bytes.Buffer{},
	}
}

// Write method writes bytes to the error accumulator.
func (e *DefaultErrorAccumulator) Write(p []byte) error {
	_, err := e.Buffer.Write(p)
	if err != nil {
		return fmt.Errorf("error accumulator write error, %w", err)
	}
	return nil
}

// Bytes method returns the bytes of the error accumulator.
func (e *DefaultErrorAccumulator) Bytes() (errBytes []byte) {
	if e.Buffer.Len() == 0 {
		return
	}
	errBytes = e.Buffer.Bytes()
	return
}

// NewStreamReader creates a new stream reader.
func NewStreamReader[Q any, T Streamer[Q]](
	response *http.Response,
	emptyMessagesLimit uint,
) *StreamReader[T] {
	stream := &StreamReader[T]{
		emptyMessagesLimit: emptyMessagesLimit,
		isFinished:         false,
		Header:             response.Header,
		Reader:             bufio.NewReader(response.Body),
		response:           response,
		ErrAccumulator:     NewErrorAccumulator(),
	}
	stream.Header = response.Header
	return stream
}
