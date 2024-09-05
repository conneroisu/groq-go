package groq_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	groq "github.com/conneroisu/groq-go"
	"github.com/stretchr/testify/assert"
)

// // TestChatCompletionsStreamInvalidModel tests the CreateChatCompletionStream method with an invalid model.
//
//	func TestChatCompletionsStreamInvalidModel(t *testing.T) {
//	        client, server, teardown := setupGroqTestServer()
//	        defer teardown()
//	        server.RegisterHandler(
//	                "/v1/chat/completions",
//	                func(w http.ResponseWriter, _ *http.Request) {
//	                },
//	        )
//	        ctx := context.Background()
//	        req := groq.ChatCompletionRequest{
//	                MaxTokens: 5,
//	                Model:     "adaadsf",
//	                Messages: []groq.ChatCompletionMessage{
//	                        {
//	                                Role:    groq.ChatMessageRoleUser,
//	                                Content: "Hello!",
//	                        },
//	                },
//	        }
//	        _, err := client.CreateChatCompletionStream(ctx, req)
//	        if !errors.Is(err, groq.ErrChatCompletionInvalidModel{}) {
//	                t.Fatalf(
//	                        "CreateChatCompletionStream should return ErrChatCompletionInvalidModel, but returned: %v",
//	                        err,
//	                )
//	        }
//	}

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
			Model:     groq.Llama3_8b_8192,
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
			Model:             groq.Llama3_8b_8192,
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
			Model:             groq.Llama3_8b_8192,
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
			Model:     groq.Llama3_8b_8192,
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
			Model:     groq.Llama3_8b_8192,
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
			Model:     groq.Llama3_8b_8192,
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
			Model:     groq.Llama3_8b_8192,
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
			Model:     groq.Llama3_8b_8192,
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
			Model:     groq.Llama3_8b_8192,
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
			Model:             groq.Llama3_8b_8192,
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
			Model:             groq.Llama3_8b_8192,
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
			Model:             groq.Llama3_8b_8192,
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
