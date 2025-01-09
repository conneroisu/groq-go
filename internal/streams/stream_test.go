package streams_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/internal/streams"
	"github.com/conneroisu/groq-go/pkg/groqerr"
	"github.com/conneroisu/groq-go/internal/test"
	"github.com/stretchr/testify/assert"
)

// TestStreamReaderReturnsUnmarshalerErrors tests the stream reader returns an unmarshaler error.
func TestStreamReaderReturnsUnmarshalerErrors(t *testing.T) {
	stream := &streams.StreamReader[groq.ChatCompletionStreamResponse]{
		ErrAccumulator: streams.NewErrorAccumulator(),
	}

	respErr := stream.UnmarshalError()
	if respErr != nil {
		t.Fatalf("Did not return nil with empty buffer: %v", respErr)
	}

	err := stream.ErrAccumulator.Write([]byte("{"))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	respErr = stream.UnmarshalError()
	if respErr != nil {
		t.Fatalf("Did not return nil when unmarshaler failed: %v", respErr)
	}
}

// TestStreamReaderReturnsErrTooManyEmptyStreamMessages tests the stream reader returns an error when the stream has too many empty messages.
func TestStreamReaderReturnsErrTooManyEmptyStreamMessages(t *testing.T) {
	a := assert.New(t)
	reader := &http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte("\n\n\n\n"))),
	}
	stream := streams.NewStreamReader[groq.ChatCompletionStreamResponse](
		reader.Body,
		map[string][]string{
			"Content-Type": {"text/event-stream"},
		},
		3,
	)
	_, err := stream.Recv()
	a.ErrorIs(
		err,
		groqerr.ErrTooManyEmptyStreamMessages{},
		"Did not return error when recv failed",
		err.Error(),
	)
}

// TestStreamReaderReturnsErrTestErrorAccumulatorWriteFailed tests the stream reader returns an error when the error accumulator fails to write.
func TestStreamReaderReturnsErrTestErrorAccumulatorWriteFailed(t *testing.T) {
	a := assert.New(t)
	reader := &http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte("\n"))),
	}
	stream := streams.NewStreamReader[groq.ChatCompletionStreamResponse](
		reader.Body,
		map[string][]string{
			"Content-Type": {"text/event-stream"},
		},
		0,
	)
	_, err := stream.Recv()
	a.ErrorIs(
		err,
		groqerr.ErrTooManyEmptyStreamMessages{},
		"Did not return error when write failed",
		err.Error(),
	)
}

// Test the `Recv` method with multiple empty messages triggering an error
func TestStreamReader_TooManyEmptyMessages(t *testing.T) {
	data := "\n\n\n\n\n\n"
	resp := &http.Response{
		Body: io.NopCloser(bytes.NewBufferString(data)),
	}
	stream := streams.NewStreamReader[*groq.ChatCompletionStreamResponse](
		resp.Body,
		map[string][]string{
			"Content-Type": {"text/event-stream"},
		},
		5,
	)

	_, err := stream.Recv()
	assert.ErrorIs(t, err, groqerr.ErrTooManyEmptyStreamMessages{})
}

// Test the `Close` method
func TestStreamReader_Close(t *testing.T) {
	resp := &http.Response{
		Body: io.NopCloser(bytes.NewBufferString("")),
	}
	stream := streams.NewStreamReader[groq.ChatCompletionStreamResponse](
		resp.Body,
		map[string][]string{
			"Content-Type": {"text/event-stream"},
		},
		5,
	)

	err := stream.Close()
	assert.NoError(t, err)
}

func TestErrorAccumulatorBytes(t *testing.T) {
	accumulator := &streams.DefaultErrorAccumulator{
		Buffer: &bytes.Buffer{},
	}

	errBytes := accumulator.Bytes()
	if len(errBytes) != 0 {
		t.Fatalf("Did not return nil with empty bytes: %s", string(errBytes))
	}

	err := accumulator.Write([]byte("{}"))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	errBytes = accumulator.Bytes()
	if len(errBytes) == 0 {
		t.Fatalf(
			"Did not return error bytes when has error: %s",
			string(errBytes),
		)
	}
}

func TestErrorByteWriteErrors(t *testing.T) {
	accumulator := &streams.DefaultErrorAccumulator{
		Buffer: &test.FailingErrorBuffer{},
	}
	err := accumulator.Write([]byte("{"))
	if !errors.Is(err, test.ErrTestErrorAccumulatorWriteFailed{}) {
		t.Fatalf("Did not return error when write failed: %v", err)
	}
}
