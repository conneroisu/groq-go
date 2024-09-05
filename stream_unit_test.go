package groq_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/conneroisu/groq-go"
	"github.com/stretchr/testify/assert"
)

// TestCompletionsStreamWrongModel tests the completion stream returns an error when the model is not supported.
func TestCompletionsStreamWrongModel(t *testing.T) {
	a := assert.New(t)
	client, err := groq.NewClient(
		"whatever",
		groq.WithBaseURL("http://localhost/v1"),
	)
	a.NoError(err, "NewClient returned error")

	_, err = client.CreateCompletionStream(
		context.Background(),
		groq.CompletionRequest{
			MaxTokens: 5,
			Model:     groq.Whisper_Large_V3,
		},
	)
	if !errors.Is(
		err,
		groq.ErrCompletionUnsupportedModel{
			Model: groq.Whisper_Large_V3,
		},
	) {
		t.Fatalf(
			"CreateCompletion should return ErrCompletionUnsupportedModel, but returned: %v",
			err,
		)
	}
}

func TestCreateCompletionStream(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/completions",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")

			// Send test responses
			dataBytes := []byte{}
			dataBytes = append(dataBytes, []byte("event: message\n")...)
			data := `{"id":"1","object":"completion","created":1598069254,"model":"text-davinci-002","choices":[{"text":"response1","finish_reason":"max_tokens"}]}`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			dataBytes = append(dataBytes, []byte("event: message\n")...)
			data = `{"id":"2","object":"completion","created":1598069255,"model":"text-davinci-002","choices":[{"text":"response2","finish_reason":"max_tokens"}]}`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			dataBytes = append(dataBytes, []byte("event: done\n")...)
			dataBytes = append(dataBytes, []byte("data: [DONE]\n\n")...)

			_, err := w.Write(dataBytes)
			a.NoError(err, "Write error")
		},
	)

	stream, err := client.CreateCompletionStream(
		context.Background(),
		groq.CompletionRequest{
			Prompt:    "Ex falso quodlibet",
			Model:     "text-davinci-002",
			MaxTokens: 10,
			Stream:    true,
		},
	)
	a.NoError(err, "CreateCompletionStream returned error")
	defer stream.Close()

	expectedResponses := []groq.CompletionResponse{
		{
			ID:      "1",
			Object:  "completion",
			Created: 1598069254,
			Model:   "text-davinci-002",
			Choices: []groq.CompletionChoice{
				{Text: "response1", FinishReason: "max_tokens"},
			},
		},
		{
			ID:      "2",
			Object:  "completion",
			Created: 1598069255,
			Model:   "text-davinci-002",
			Choices: []groq.CompletionChoice{
				{Text: "response2", FinishReason: "max_tokens"},
			},
		},
	}

	for ix, expectedResponse := range expectedResponses {
		receivedResponse, streamErr := stream.Recv()
		if streamErr != nil {
			t.Errorf("stream.Recv() failed: %v", streamErr)
		}
		if !compareResponses(expectedResponse, receivedResponse) {
			t.Errorf(
				"Stream response %v is %v, expected %v",
				ix,
				receivedResponse,
				expectedResponse,
			)
		}
	}

	_, streamErr := stream.Recv()
	if !errors.Is(streamErr, io.EOF) {
		t.Errorf("stream.Recv() did not return EOF in the end: %v", streamErr)
	}

	_, streamErr = stream.Recv()
	if !errors.Is(streamErr, io.EOF) {
		t.Errorf(
			"stream.Recv() did not return EOF when the stream is finished: %v",
			streamErr,
		)
	}
}

// TestCreateCompletionStreamTooManyEmptyStreamMessagesError tests the completion stream returns an error when the stream has too many empty messages.
func TestCreateCompletionStreamTooManyEmptyStreamMessagesError(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/completions",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")

			// Send test responses
			dataBytes := []byte{}
			dataBytes = append(dataBytes, []byte("event: message\n")...)
			data := `{"id":"1","object":"completion","created":1598069254,"model":"text-davinci-002","choices":[{"text":"response1","finish_reason":"max_tokens"}]}`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			// Totally 301 empty messages (300 is the limit)
			for i := 0; i < 299; i++ {
				dataBytes = append(dataBytes, '\n')
			}

			dataBytes = append(dataBytes, []byte("event: message\n")...)
			data = `{"id":"2","object":"completion","created":1598069255,"model":"text-davinci-002","choices":[{"text":"response2","finish_reason":"max_tokens"}]}`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			dataBytes = append(dataBytes, []byte("event: done\n")...)
			dataBytes = append(dataBytes, []byte("data: [DONE]\n\n")...)

			_, err := w.Write(dataBytes)
			a.NoError(err, "Write error")
		},
	)

	stream, err := client.CreateCompletionStream(
		context.Background(),
		groq.CompletionRequest{
			Prompt:    "Ex falso quodlibet",
			Model:     "text-davinci-002",
			MaxTokens: 10,
			Stream:    true,
		},
	)
	a.NoError(err, "CreateCompletionStream returned error")
	defer stream.Close()

	_, _ = stream.Recv()
	_, streamErr := stream.Recv()
	if !errors.Is(streamErr, groq.ErrTooManyEmptyStreamMessages{}) {
		t.Errorf(
			"TestCreateCompletionStreamTooManyEmptyStreamMessagesError did not return ErrTooManyEmptyStreamMessages",
		)
	}
}

// TestCreateCompletionStreamUnexpectedTerminatedError tests the completion stream returns an error when the stream is terminated unexpectedly.
func TestCreateCompletionStreamUnexpectedTerminatedError(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/completions",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")

			// Send test responses
			dataBytes := []byte{}
			dataBytes = append(dataBytes, []byte("event: message\n")...)
			data := `{"id":"1","object":"completion","created":1598069254,"model":"text-davinci-002","choices":[{"text":"response1","finish_reason":"max_tokens"}]}`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			// Stream is terminated without sending "done" message

			_, err := w.Write(dataBytes)
			a.NoError(err, "Write error")
		},
	)

	stream, err := client.CreateCompletionStream(
		context.Background(),
		groq.CompletionRequest{
			Prompt:    "Ex falso quodlibet",
			Model:     "text-davinci-002",
			MaxTokens: 10,
			Stream:    true,
		},
	)
	a.NoError(err, "CreateCompletionStream returned error")
	defer stream.Close()

	_, _ = stream.Recv()
	_, streamErr := stream.Recv()
	if !errors.Is(streamErr, io.EOF) {
		t.Errorf(
			"TestCreateCompletionStreamUnexpectedTerminatedError did not return io.EOF",
		)
	}
}

// TestCreateCompletionStreamBrokenJSONError tests the completion stream returns an error when the stream is broken.
func TestCreateCompletionStreamBrokenJSONError(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/completions",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")

			// Send test responses
			dataBytes := []byte{}
			dataBytes = append(dataBytes, []byte("event: message\n")...)
			data := `{"id":"1","object":"completion","created":1598069254,"model":"text-davinci-002","choices":[{"text":"response1","finish_reason":"max_tokens"}]}`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			// Send broken json
			dataBytes = append(dataBytes, []byte("event: message\n")...)
			data = `{"id":"2","object":"completion","created":1598069255,"model":`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			dataBytes = append(dataBytes, []byte("event: done\n")...)
			dataBytes = append(dataBytes, []byte("data: [DONE]\n\n")...)

			_, err := w.Write(dataBytes)
			a.NoError(err, "Write error")
		},
	)

	stream, err := client.CreateCompletionStream(
		context.Background(),
		groq.CompletionRequest{
			Prompt:    "Ex falso quodlibet",
			Model:     "text-davinci-002",
			MaxTokens: 10,
			Stream:    true,
		},
	)
	a.NoError(err, "CreateCompletionStream returned error")
	defer stream.Close()

	_, _ = stream.Recv()
	_, streamErr := stream.Recv()
	var syntaxError *json.SyntaxError
	if !errors.As(streamErr, &syntaxError) {
		t.Errorf(
			"TestCreateCompletionStreamBrokenJSONError did not return json.SyntaxError",
		)
	}
}

func TestCreateCompletionStreamReturnTimeoutError(t *testing.T) {
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/completions",
		func(http.ResponseWriter, *http.Request) {
			time.Sleep(10 * time.Nanosecond)
		},
	)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Nanosecond)
	defer cancel()

	_, err := client.CreateCompletionStream(ctx, groq.CompletionRequest{
		Prompt:    "Ex falso quodlibet",
		Model:     "text-davinci-002",
		MaxTokens: 10,
		Stream:    true,
	})
	if err == nil {
		t.Fatal("Did not return error")
	}
	if !os.IsTimeout(err) {
		t.Fatal("Did not return timeout error")
	}
}

// Helper funcs.
func compareResponses(r1, r2 groq.CompletionResponse) bool {
	if r1.ID != r2.ID || r1.Object != r2.Object || r1.Created != r2.Created ||
		r1.Model != r2.Model {
		return false
	}
	if len(r1.Choices) != len(r2.Choices) {
		return false
	}
	for i := range r1.Choices {
		if !compareResponseChoices(r1.Choices[i], r2.Choices[i]) {
			return false
		}
	}
	return true
}

func compareResponseChoices(c1, c2 groq.CompletionChoice) bool {
	if c1.Text != c2.Text || c1.FinishReason != c2.FinishReason {
		return false
	}
	return true
}
