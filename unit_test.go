package groq_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestTestServer(t *testing.T) {
	num := rand.Intn(100)
	a := assert.New(t)
	ctx := context.Background()
	client, err := groq.NewClient(os.Getenv("GROQ_KEY"))
	a.NoError(err, "NewClient error")
	strm, err := client.CreateChatCompletionStream(
		ctx,
		groq.ChatCompletionRequest{
			Model: groq.Llama38B8192,
			Messages: []groq.ChatCompletionMessage{
				{
					Role: groq.ChatMessageRoleUser,
					Content: fmt.Sprintf(`
problem: %d
You have a six-sided die that you roll once. Let $R{i}$ denote the event that the roll is $i$. Let $G{j}$ denote the event that the roll is greater than $j$. Let $E$ denote the event that the roll of the die is even-numbered.
(a) What is $P\left[R{3} \mid G{1}\right]$, the conditional probability that 3 is rolled given that the roll is greater than 1 ?
(b) What is the conditional probability that 6 is rolled given that the roll is greater than 3 ?
(c) What is the $P\left[G_{3} \mid E\right]$, the conditional probability that the roll is greater than 3 given that the roll is even?
(d) Given that the roll is greater than 3, what is the conditional probability that the roll is even?
					`, num,
					),
				},
			},
			MaxTokens: 2000,
			Stream:    true,
		},
	)
	a.NoError(err, "CreateCompletionStream error")

	i := 0
	for {
		i++
		val, err := strm.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		// t.Logf("%d %s\n", i, val.Choices[0].Delta.Content)
		print(val.Choices[0].Delta.Content)
	}
}

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
			Model:     groq.WhisperLargeV3,
		},
	)
	if !errors.Is(
		err,
		groq.ErrCompletionStreamNotSupported{},
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

// TestModerate tests the Moderate method of the client.
func TestModerate(t *testing.T) {
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/chat/completions",
		handleModerationEndpoint,
	)
	mod, err := client.Moderate(context.Background(), groq.ModerationRequest{
		Model: groq.ModerationTextStable,
		Messages: []groq.ChatCompletionMessage{
			{
				Role:    groq.ChatMessageRoleUser,
				Content: "I want to kill them.",
			},
		},
	})
	a := assert.New(t)
	a.NoError(err, "Moderation error")
	a.Equal(true, mod.Flagged)
	a.Equal(
		mod.Categories,
		[]groq.HarmfulCategory{
			groq.CategoryViolentCrimes,
			groq.CategoryNonviolentCrimes,
		},
	)
}

// handleModerationEndpoint handles the moderation endpoint.
func handleModerationEndpoint(w http.ResponseWriter, r *http.Request) {
	response := groq.ChatCompletionResponse{
		ID:      "chatcmpl-123",
		Object:  "chat.completion",
		Created: 1693721698,
		Model:   groq.ModerationTextStable,
		Choices: []groq.ChatCompletionChoice{
			{
				Message: groq.ChatCompletionMessage{
					Role:    groq.ChatMessageRoleAssistant,
					Content: "unsafe\nS1,S2",
				},
				FinishReason: "stop",
			},
		},
	}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(response)
	if err != nil {
		http.Error(
			w,
			"could not encode response",
			http.StatusInternalServerError,
		)
		return
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(
			w,
			"could not write response",
			http.StatusInternalServerError,
		)
		return
	}
}

// TestCompletionWithStream tests the CreateCompletion method with a stream.
func TestCompletionWithStream(t *testing.T) {
	a := assert.New(t)
	client, err := groq.NewClient(
		"whatever",
		groq.WithBaseURL("http://localhost/v1"),
	)
	a.NoError(err, "NewClient error")

	ctx := context.Background()
	req := groq.CompletionRequest{Stream: true}
	_, err = client.CreateCompletion(ctx, req)
	if !errors.Is(err, groq.ErrCompletionStreamNotSupported{}) {
		t.Fatalf(
			"CreateCompletion didn't return ErrCompletionStreamNotSupported",
		)
	}
}

// TestCompletions Tests the completions endpoint of the API using the mocked server.
func TestCompletions(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler("/v1/completions", handleCompletionEndpoint)
	req := groq.CompletionRequest{
		MaxTokens: 5,
		Model:     "ada",
		Prompt:    "Lorem ipsum",
	}
	_, err := client.CreateCompletion(context.Background(), req)
	a.NoError(err, "CreateCompletion error")
}

// handleCompletionEndpoint Handles the completion endpoint by the test server.
func handleCompletionEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error
	var resBytes []byte

	// completions only accepts POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	var completionReq groq.CompletionRequest
	if completionReq, err = getCompletionBody(r); err != nil {
		http.Error(w, "could not read request", http.StatusInternalServerError)
		return
	}
	res := groq.CompletionResponse{
		ID:      strconv.Itoa(int(time.Now().Unix())),
		Object:  "test-object",
		Created: time.Now().Unix(),
		// would be nice to validate Model during testing, but
		// this may not be possible with how much upkeep
		// would be required / wouldn't make much sense
		Model: completionReq.Model,
	}
	// create completions
	n := completionReq.N
	if n == 0 {
		n = 1
	}
	for i := 0; i < n; i++ {
		// generate a random string of length completionReq.Length
		completionStr := strings.Repeat("a", completionReq.MaxTokens)
		if completionReq.Echo {
			completionStr = completionReq.Prompt.(string) + completionStr
		}
		res.Choices = append(res.Choices, groq.CompletionChoice{
			Text:  completionStr,
			Index: i,
		})
	}
	inputTokens := numTokens(completionReq.Prompt.(string)) * n
	completionTokens := completionReq.MaxTokens * n
	res.Usage = groq.Usage{
		PromptTokens:     inputTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      inputTokens + completionTokens,
	}
	resBytes, _ = json.Marshal(res)
	fmt.Fprintln(w, string(resBytes))
}

// numTokens Returns the number of GPT-3 encoded tokens in the given text.
// This function approximates based on the rule of thumb stated by OpenAI:
// https://beta.openai.com/tokenizer
//
// TODO: implement an actual tokenizer for each model available and use that
// instead.
func numTokens(s string) int {
	return int(float32(len(s)) / 4)
}

// getCompletionBody Returns the body of the request to create a completion.
func getCompletionBody(r *http.Request) (groq.CompletionRequest, error) {
	completion := groq.CompletionRequest{}
	// read the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return groq.CompletionRequest{}, err
	}
	err = json.Unmarshal(reqBody, &completion)
	if err != nil {
		return groq.CompletionRequest{}, err
	}
	return completion, nil
}

func setupGroqTestServer() (
	client *groq.Client,
	server *test.ServerTest,
	teardown func(),
) {
	server = test.NewTestServer()
	ts := server.GroqTestServer()
	ts.Start()
	teardown = ts.Close
	client, err := groq.NewClient(
		test.GetTestToken(),
		groq.WithBaseURL(ts.URL+"/v1"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func TestEmptyKeyClientCreation(t *testing.T) {
	client, err := groq.NewClient("")
	a := assert.New(t)
	a.Error(err, "NewClient should return error")
	a.Nil(client, "NewClient should return nil")
}

// TestCreateChatCompletionStream tests the CreateChatCompletionStream method.
func TestCreateChatCompletionStream(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/chat/completions",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")

			// Send test responses
			dataBytes := []byte{}
			dataBytes = append(dataBytes, []byte("event: message\n")...)
			data := `{"id":"1","object":"completion","created":1598069254,"model":"llama3-groq-70b-8192-tool-use-preview","system_fingerprint": "fp_d9767fc5b9","choices":[{"index":0,"delta":{"content":"response1"},"finish_reason":"max_tokens"}]}`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			dataBytes = append(dataBytes, []byte("event: message\n")...)
			data = `{"id":"2","object":"completion","created":1598069255,"model":"llama3-groq-70b-8192-tool-use-preview","system_fingerprint": "fp_d9767fc5b9","choices":[{"index":0,"delta":{"content":"response2"},"finish_reason":"max_tokens"}]}`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			dataBytes = append(dataBytes, []byte("event: done\n")...)
			dataBytes = append(dataBytes, []byte("data: [DONE]\n\n")...)

			_, err := w.Write(dataBytes)
			a.NoError(err, "Write error")
		},
	)

	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		groq.ChatCompletionRequest{
			MaxTokens: 5,
			Model:     groq.Llama38B8192,
			Messages: []groq.ChatCompletionMessage{
				{
					Role:    groq.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
			Stream: true,
		},
	)
	a.NoError(err, "CreateCompletionStream returned error")
	defer stream.Close()

	expectedResponses := []groq.ChatCompletionStreamResponse{
		{
			ID:                "1",
			Object:            "completion",
			Created:           1598069254,
			Model:             groq.Llama38B8192,
			SystemFingerprint: "fp_d9767fc5b9",
			Choices: []groq.ChatCompletionStreamChoice{
				{
					Delta: groq.ChatCompletionStreamChoiceDelta{
						Content: "response1",
					},
					FinishReason: "max_tokens",
				},
			},
		},
		{
			ID:                "2",
			Object:            "completion",
			Created:           1598069255,
			Model:             groq.Llama38B8192,
			SystemFingerprint: "fp_d9767fc5b9",
			Choices: []groq.ChatCompletionStreamChoice{
				{
					Delta: groq.ChatCompletionStreamChoiceDelta{
						Content: "response2",
					},
					FinishReason: "max_tokens",
				},
			},
		},
	}

	for ix, expectedResponse := range expectedResponses {
		b, _ := json.Marshal(expectedResponse)
		t.Logf("%d: %s", ix, string(b))

		receivedResponse, streamErr := stream.Recv()
		a.NoError(streamErr, "stream.Recv() failed")
		if !compareChatResponses(t, expectedResponse, receivedResponse) {
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
		t.Errorf("stream.Recv() did not return EOF in the end: %v", streamErr)
	}

	_, streamErr = stream.Recv()

	a.ErrorIs(
		streamErr,
		io.EOF,
		"stream.Recv() did not return EOF when the stream is finished",
	)
	if !errors.Is(streamErr, io.EOF) {
		t.Errorf(
			"stream.Recv() did not return EOF when the stream is finished: %v",
			streamErr,
		)
	}
}

// TestCreateChatCompletionStreamError tests the CreateChatCompletionStream function with an error
// in the response.
func TestCreateChatCompletionStreamError(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/chat/completions",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")

			// Send test responses
			dataBytes := []byte{}
			dataStr := []string{
				`{`,
				`"error": {`,
				`"message": "Incorrect API key provided: gsk-***************************************",`,
				`"type": "invalid_request_error",`,
				`"param": null,`,
				`"code": "invalid_api_key"`,
				`}`,
				`}`,
			}
			for _, str := range dataStr {
				dataBytes = append(dataBytes, []byte(str+"\n")...)
			}

			_, err := w.Write(dataBytes)
			a.NoError(err, "Write error")
		},
	)

	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		groq.ChatCompletionRequest{
			MaxTokens: 5,
			Model:     groq.Llama38B8192,
			Messages: []groq.ChatCompletionMessage{
				{
					Role:    groq.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
			Stream: true,
		},
	)
	a.NoError(err, "CreateCompletionStream returned error")
	defer stream.Close()

	_, streamErr := stream.Recv()
	a.Error(streamErr, "stream.Recv() did not return error")

	var apiErr *groq.APIError
	if !errors.As(streamErr, &apiErr) {
		t.Errorf("stream.Recv() did not return APIError")
	}
	t.Logf("%+v\n", apiErr)
}

func TestCreateChatCompletionStreamWithHeaders(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	xCustomHeader := "x-custom-header"
	xCustomHeaderValue := "x-custom-header-value"
	server.RegisterHandler(
		"/v1/chat/completions",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set(xCustomHeader, xCustomHeaderValue)

			// Send test responses
			dataBytes := []byte(
				`data: {"error":{"message":"The server had an error while processing your request. Sorry about that!", "type":"server_ error", "param":null,"code":null}}`,
			)
			dataBytes = append(dataBytes, []byte("\n\ndata: [DONE]\n\n")...)

			_, err := w.Write(dataBytes)
			a.NoError(err, "Write error")
		},
	)

	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		groq.ChatCompletionRequest{
			MaxTokens: 5,
			Model:     groq.Llama38B8192,
			Messages: []groq.ChatCompletionMessage{
				{
					Role:    groq.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
			Stream: true,
		},
	)
	a.NoError(err, "CreateCompletionStream returned error")
	defer stream.Close()

	value := stream.Header.Get(xCustomHeader)
	if value != xCustomHeaderValue {
		t.Errorf("expected %s to be %s", xCustomHeaderValue, value)
	}
}

func TestCreateChatCompletionStreamWithRatelimitHeaders(t *testing.T) {
	client, server, teardown := setupGroqTestServer()
	a := assert.New(t)
	rateLimitHeaders := map[string]interface{}{
		"x-ratelimit-limit-requests":     100,
		"x-ratelimit-limit-tokens":       1000,
		"x-ratelimit-remaining-requests": 99,
		"x-ratelimit-remaining-tokens":   999,
		"x-ratelimit-reset-requests":     "1s",
		"x-ratelimit-reset-tokens":       "1m",
	}
	defer teardown()
	server.RegisterHandler(
		"/v1/chat/completions",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			for k, v := range rateLimitHeaders {
				switch val := v.(type) {
				case int:
					w.Header().Set(k, strconv.Itoa(val))
				default:
					w.Header().Set(k, fmt.Sprintf("%s", v))
				}
			}

			// Send test responses
			dataBytes := []byte(
				`data: {"error":{"message":"The server had an error while processing your request. Sorry about that!", "type":"server_ error", "param":null,"code":null}}`,
			)
			dataBytes = append(dataBytes, []byte("\n\ndata: [DONE]\n\n")...)

			_, err := w.Write(dataBytes)
			a.NoError(err, "Write error")
		},
	)

	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		groq.ChatCompletionRequest{
			MaxTokens: 5,
			Model:     groq.Llama38B8192,
			Messages: []groq.ChatCompletionMessage{
				{
					Role:    groq.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
			Stream: true,
		},
	)
	a.NoError(err, "CreateCompletionStream returned error")
	defer stream.Close()

	headers := newRateLimitHeaders(stream.Header)
	bs1, _ := json.Marshal(headers)
	bs2, _ := json.Marshal(rateLimitHeaders)
	if string(bs1) != string(bs2) {
		t.Errorf("expected rate limit header %s to be %s", bs2, bs1)
	}
}

// newRateLimitHeaders creates a new RateLimitHeaders from an http.Header.
func newRateLimitHeaders(h http.Header) groq.RateLimitHeaders {
	limitReq, _ := strconv.Atoi(h.Get("x-ratelimit-limit-requests"))
	limitTokens, _ := strconv.Atoi(h.Get("x-ratelimit-limit-tokens"))
	remainingReq, _ := strconv.Atoi(h.Get("x-ratelimit-remaining-requests"))
	remainingTokens, _ := strconv.Atoi(h.Get("x-ratelimit-remaining-tokens"))
	return groq.RateLimitHeaders{
		LimitRequests:     limitReq,
		LimitTokens:       limitTokens,
		RemainingRequests: remainingReq,
		RemainingTokens:   remainingTokens,
		ResetRequests:     groq.ResetTime(h.Get("x-ratelimit-reset-requests")),
		ResetTokens:       groq.ResetTime(h.Get("x-ratelimit-reset-tokens")),
	}
}

func TestCreateChatCompletionStreamErrorWithDataPrefix(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/chat/completions",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")

			// Send test responses
			dataBytes := []byte(
				`data: {"error":{"message":"The server had an error while processing your request. Sorry about that!", "type":"server_ error", "param":null,"code":null}}`,
			)
			dataBytes = append(dataBytes, []byte("\n\ndata: [DONE]\n\n")...)

			_, err := w.Write(dataBytes)
			a.NoError(err, "Write error")
		},
	)

	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		groq.ChatCompletionRequest{
			MaxTokens: 5,
			Model:     groq.Llama38B8192,
			Messages: []groq.ChatCompletionMessage{
				{
					Role:    groq.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
			Stream: true,
		},
	)
	a.NoError(err, "CreateCompletionStream returned error")
	defer stream.Close()

	_, streamErr := stream.Recv()
	a.Error(streamErr, "stream.Recv() did not return error")

	var apiErr *groq.APIError
	if !errors.As(streamErr, &apiErr) {
		t.Errorf("stream.Recv() did not return APIError")
	}
	t.Logf("%+v\n", apiErr)
}

func TestCreateChatCompletionStreamRateLimitError(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/chat/completions",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(429)

			// Send test responses
			dataBytes := []byte(`{"error":{` +
				`"message": "You are sending requests too quickly.",` +
				`"type":"rate_limit_reached",` +
				`"param":null,` +
				`"code":"rate_limit_reached"}}`)

			_, err := w.Write(dataBytes)
			a.NoError(err, "Write error")
		},
	)
	_, err := client.CreateChatCompletionStream(
		context.Background(),
		groq.ChatCompletionRequest{
			MaxTokens: 5,
			Model:     groq.Llama38B8192,
			Messages: []groq.ChatCompletionMessage{
				{
					Role:    groq.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
			Stream: true,
		},
	)
	var apiErr *groq.APIError
	if !errors.As(err, &apiErr) {
		t.Errorf(
			"TestCreateChatCompletionStreamRateLimitError did not return APIError",
		)
	}
	t.Logf("%+v\n", apiErr)
}

func TestCreateChatCompletionStreamStreamOptions(t *testing.T) {
	a := assert.New(t)
	client, server, teardown := setupGroqTestServer()
	defer teardown()

	server.RegisterHandler(
		"/v1/chat/completions",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")

			// Send test responses
			var dataBytes []byte
			data := `{"id":"1","object":"completion","created":1598069254,"model":"llama3-groq-70b-8192-tool-use-preview","system_fingerprint": "fp_d9767fc5b9","choices":[{"index":0,"delta":{"content":"response1"},"finish_reason":"max_tokens"}],"usage":null}`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			data = `{"id":"2","object":"completion","created":1598069255,"model":"llama3-groq-70b-8192-tool-use-preview","system_fingerprint": "fp_d9767fc5b9","choices":[{"index":0,"delta":{"content":"response2"},"finish_reason":"max_tokens"}],"usage":null}`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			data = `{"id":"3","object":"completion","created":1598069256,"model":"llama3-groq-70b-8192-tool-use-preview","system_fingerprint": "fp_d9767fc5b9","choices":[],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`
			dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

			dataBytes = append(dataBytes, []byte("data: [DONE]\n\n")...)

			_, err := w.Write(dataBytes)
			a.NoError(err, "Write error")
		},
	)

	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		groq.ChatCompletionRequest{
			MaxTokens: 5,
			Model:     groq.Llama38B8192,
			Messages: []groq.ChatCompletionMessage{
				{
					Role:    groq.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
			Stream: true,
			StreamOptions: &groq.StreamOptions{
				IncludeUsage: true,
			},
		},
	)
	a.NoError(err, "CreateCompletionStream returned error")
	defer stream.Close()
	expectedResponses := []groq.ChatCompletionStreamResponse{
		{
			ID:                "1",
			Object:            "completion",
			Created:           1598069254,
			Model:             groq.Llama38B8192,
			SystemFingerprint: "fp_d9767fc5b9",
			Choices: []groq.ChatCompletionStreamChoice{
				{
					Delta: groq.ChatCompletionStreamChoiceDelta{
						Content: "response1",
					},
					FinishReason: "max_tokens",
				},
			},
		},
		{
			ID:                "2",
			Object:            "completion",
			Created:           1598069255,
			Model:             groq.Llama38B8192,
			SystemFingerprint: "fp_d9767fc5b9",
			Choices: []groq.ChatCompletionStreamChoice{
				{
					Delta: groq.ChatCompletionStreamChoiceDelta{
						Content: "response2",
					},
					FinishReason: "max_tokens",
				},
			},
		},
		{
			ID:                "3",
			Object:            "completion",
			Created:           1598069256,
			Model:             groq.Llama38B8192,
			SystemFingerprint: "fp_d9767fc5b9",
			Choices:           []groq.ChatCompletionStreamChoice{},
			Usage: &groq.Usage{
				PromptTokens:     1,
				CompletionTokens: 1,
				TotalTokens:      2,
			},
		},
	}

	for ix, expectedResponse := range expectedResponses {
		ix++
		b, _ := json.Marshal(expectedResponse)
		t.Logf("%d: %s", ix, string(b))

		receivedResponse, streamErr := stream.Recv()
		if !errors.Is(streamErr, io.EOF) {
			a.NoError(streamErr, "stream.Recv() failed")
		}
		if !compareChatResponses(t, expectedResponse, receivedResponse) {
			t.Errorf(
				"Stream response %v: %v,BUT expected %v",
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

	a.ErrorIs(
		streamErr,
		io.EOF,
		"stream.Recv() did not return EOF when the stream is finished",
	)
	if !errors.Is(streamErr, io.EOF) {
		t.Errorf(
			"stream.Recv() did not return EOF when the stream is finished: %v",
			streamErr,
		)
	}
}

// Helper funcs.
func compareChatResponses(
	t *testing.T,
	r1, r2 groq.ChatCompletionStreamResponse,
) bool {
	if r1.ID != r2.ID {
		t.Logf("Not Equal ID: %v", r1.ID)
		return false
	}
	if r1.Object != r2.Object {
		t.Logf("Not Equal Object: %v", r1.Object)
		return false
	}
	if r1.Created != r2.Created {
		t.Logf("Not Equal Created: %v", r1.Created)
		return false
	}
	if len(r1.Choices) != len(r2.Choices) {
		t.Logf("Not Equal Choices: %v", r1.Choices)
		return false
	}
	for i := range r1.Choices {
		if !compareChatStreamResponseChoices(r1.Choices[i], r2.Choices[i]) {
			t.Logf("Not Equal Choices: %v", r1.Choices[i])
			return false
		}
	}
	if r1.Usage != nil || r2.Usage != nil {
		if r1.Usage == nil || r2.Usage == nil {
			return false
		}
		if r1.Usage.PromptTokens != r2.Usage.PromptTokens ||
			r1.Usage.CompletionTokens != r2.Usage.CompletionTokens ||
			r1.Usage.TotalTokens != r2.Usage.TotalTokens {
			return false
		}
	}
	return true
}

func compareChatStreamResponseChoices(
	c1, c2 groq.ChatCompletionStreamChoice,
) bool {
	if c1.Index != c2.Index {
		return false
	}
	if c1.Delta.Content != c2.Delta.Content {
		return false
	}
	if c1.FinishReason != c2.FinishReason {
		return false
	}
	return true
}

// TestAudio Tests the transcription and translation endpoints of the API using the mocked server.
func TestAudio(t *testing.T) {
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler("/v1/audio/transcriptions", handleAudioEndpoint)
	server.RegisterHandler("/v1/audio/translations", handleAudioEndpoint)

	testcases := []struct {
		name     string
		createFn func(context.Context, groq.AudioRequest) (groq.AudioResponse, error)
	}{
		{
			"transcribe",
			client.CreateTranscription,
		},
		{
			"translate",
			client.CreateTranslation,
		},
	}

	ctx := context.Background()

	dir, cleanup := test.CreateTestDirectory(t)
	defer cleanup()

	a := assert.New(t)
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(dir, "fake.mp3")
			test.CreateTestFile(t, path)

			req := groq.AudioRequest{
				FilePath: path,
				Model:    "distil-whisper-large-v3-en",
			}
			_, err := tc.createFn(ctx, req)
			a.NoError(err, "audio API error")
		})

		t.Run(tc.name+" (with reader)", func(t *testing.T) {
			req := groq.AudioRequest{
				FilePath: "fake.webm",
				Reader:   bytes.NewBuffer([]byte(`some webm binary data`)),
				Model:    "whisper-3",
			}
			_, err := tc.createFn(ctx, req)
			a.NoError(err, "audio API error")
		})
	}
}

func TestAudioWithOptionalArgs(t *testing.T) {
	client, server, teardown := setupGroqTestServer()
	defer teardown()
	server.RegisterHandler("/v1/audio/transcriptions", handleAudioEndpoint)
	server.RegisterHandler("/v1/audio/translations", handleAudioEndpoint)

	testcases := []struct {
		name     string
		createFn func(context.Context, groq.AudioRequest) (groq.AudioResponse, error)
	}{
		{
			"transcribe",
			client.CreateTranscription,
		},
		{
			"translate",
			client.CreateTranslation,
		},
	}

	ctx := context.Background()

	dir, cleanup := test.CreateTestDirectory(t)
	defer cleanup()

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			a := assert.New(t)
			path := filepath.Join(dir, "fake.mp3")
			test.CreateTestFile(t, path)

			req := groq.AudioRequest{
				FilePath:    path,
				Model:       "whisper-3",
				Prompt:      "用简体中文",
				Temperature: 0.5,
				Language:    "zh",
				Format:      groq.AudioResponseFormatSRT,
			}
			_, err := tc.createFn(ctx, req)
			a.NoError(err, "audio API error")
		})
	}
}

// handleAudioEndpoint Handles the completion endpoint by the test server.
func handleAudioEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error

	// audio endpoints only accept POST requests
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, "failed to parse media type", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(mediaType, "multipart") {
		http.Error(w, "request is not multipart", http.StatusBadRequest)
	}

	boundary, ok := params["boundary"]
	if !ok {
		http.Error(w, "no boundary in params", http.StatusBadRequest)
		return
	}

	fileData := &bytes.Buffer{}
	mr := multipart.NewReader(r.Body, boundary)
	part, err := mr.NextPart()
	if err != nil && errors.Is(err, io.EOF) {
		http.Error(w, "error accessing file", http.StatusBadRequest)
		return
	}
	if _, err = io.Copy(fileData, part); err != nil {
		http.Error(w, "failed to copy file", http.StatusInternalServerError)
		return
	}

	if len(fileData.Bytes()) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "received empty file data", http.StatusBadRequest)
		return
	}

	if _, err = w.Write([]byte(`{"body": "hello"}`)); err != nil {
		http.Error(w, "failed to write body", http.StatusInternalServerError)
		return
	}
}
