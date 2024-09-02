package groq //nolint:testpackage // testing private field

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/conneroisu/groq-go/internal/test"
	"github.com/conneroisu/groq-go/internal/test/checks"
)

var errTestUnmarshalerFailed = errors.New("test unmarshaler failed")

type failingUnMarshaller struct{}

func (*failingUnMarshaller) Unmarshal(_ []byte, _ any) error {
	return errTestUnmarshalerFailed
}

func TestStreamReaderReturnsUnmarshalerErrors(t *testing.T) {
	stream := &streamReader[ChatCompletionStreamResponse]{
		errAccumulator: NewErrorAccumulator(),
		unmarshaler:    &failingUnMarshaller{},
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

func TestStreamReaderReturnsErrTooManyEmptyStreamMessages(t *testing.T) {
	stream := &streamReader[ChatCompletionStreamResponse]{
		emptyMessagesLimit: 3,
		reader:             bufio.NewReader(bytes.NewReader([]byte("\n\n\n\n"))),
		errAccumulator:     NewErrorAccumulator(),
		unmarshaler:        json.NewDecoder(bytes.NewReader([]byte("{}"))),
	}
	_, err := stream.Recv()
	checks.ErrorIs(t, err, ErrTooManyEmptyStreamMessages, "Did not return error when recv failed", err.Error())
}

func TestStreamReaderReturnsErrTestErrorAccumulatorWriteFailed(t *testing.T) {
	stream := &streamReader[ChatCompletionStreamResponse]{
		reader: bufio.NewReader(bytes.NewReader([]byte("\n"))),
		errAccumulator: &utils.DefaultErrorAccumulator{
			Buffer: &test.FailingErrorBuffer{},
		},
		unmarshaler: &utils.JSONUnmarshaler{},
	}
	_, err := stream.Recv()
	checks.ErrorIs(t, err, test.ErrTestErrorAccumulatorWriteFailed, "Did not return error when write failed", err.Error())
}
