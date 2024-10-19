package groq

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/stretchr/testify/assert"
)

// TestStreamReaderReturnsUnmarshalerErrors tests the stream reader returns an unmarshaler error.
func TestStreamReaderReturnsUnmarshalerErrors(t *testing.T) {
	stream := &streamReader[ChatCompletionStreamResponse]{
		errAccumulator: newErrorAccumulator(),
	}

	respErr := stream.unmarshalError()
	if respErr != nil {
		t.Fatalf("Did not return nil with empty buffer: %v", respErr)
	}

	err := stream.errAccumulator.Write([]byte("{"))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	respErr = stream.unmarshalError()
	if respErr != nil {
		t.Fatalf("Did not return nil when unmarshaler failed: %v", respErr)
	}
}

// TestStreamReaderReturnsErrTooManyEmptyStreamMessages tests the stream reader returns an error when the stream has too many empty messages.
func TestStreamReaderReturnsErrTooManyEmptyStreamMessages(t *testing.T) {
	a := assert.New(t)
	stream := &streamReader[ChatCompletionStreamResponse]{
		emptyMessagesLimit: 3,
		reader: bufio.NewReader(
			bytes.NewReader([]byte("\n\n\n\n")),
		),
		errAccumulator: newErrorAccumulator(),
	}
	_, err := stream.Recv()
	a.ErrorIs(
		err,
		ErrTooManyEmptyStreamMessages{},
		"Did not return error when recv failed",
		err.Error(),
	)
}

// TestStreamReaderReturnsErrTestErrorAccumulatorWriteFailed tests the stream reader returns an error when the error accumulator fails to write.
func TestStreamReaderReturnsErrTestErrorAccumulatorWriteFailed(t *testing.T) {
	a := assert.New(t)
	stream := &streamReader[ChatCompletionStreamResponse]{
		reader: bufio.NewReader(bytes.NewReader([]byte("\n"))),
		errAccumulator: &DefaultErrorAccumulator{
			Buffer: &test.FailingErrorBuffer{},
		},
	}
	_, err := stream.Recv()
	a.ErrorIs(
		err,
		test.ErrTestErrorAccumulatorWriteFailed{},
		"Did not return error when write failed",
		err.Error(),
	)
}

// Helper function to create a new `streamReader` for testing
func newStreamReader[T streamer](data string) *streamReader[T] {
	resp := &http.Response{
		Body: io.NopCloser(bytes.NewBufferString(data)),
	}
	reader := bufio.NewReader(resp.Body)

	return &streamReader[T]{
		emptyMessagesLimit: 5,
		isFinished:         false,
		reader:             reader,
		response:           resp,
		errAccumulator:     newErrorAccumulator(),
		Header:             resp.Header,
	}
}

// Test the `Recv` method with multiple empty messages triggering an error
func TestStreamReader_TooManyEmptyMessages(t *testing.T) {
	data := "\n\n\n\n\n\n"
	stream := newStreamReader[ChatCompletionStreamResponse](data)

	_, err := stream.Recv()
	assert.ErrorIs(t, err, ErrTooManyEmptyStreamMessages{})
}

// Test the `Close` method
func TestStreamReader_Close(t *testing.T) {
	stream := newStreamReader[ChatCompletionStreamResponse]("")

	err := stream.Close()
	assert.NoError(t, err)
}
