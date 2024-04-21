package gogroq

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/conneroisu/go-groq/internal"
)

// ChatRequest is a request to the chat endpoint
type ChatRequest struct {
	Messages  []Message `json:"messages"`
	Model     string    `json:"model"`
	TopP      float64   `json:"top_p"`
	MaxTokens int       `json:"max_tokens"`
	Stop      []string  `json:"stop omitempty"`
	Seed      int       `json:"seed omitempty"`
	Stream    bool      `json:"stream omitempty"`
}

// Message is a message in a chat request
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse is a response from the chat endpoint
type ChatResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int     `json:"prompt_tokens"`
		PromptTime       float64 `json:"prompt_time"`
		CompletionTokens int     `json:"completion_tokens"`
		CompletionTime   float64 `json:"completion_time"`
		ToTalTokens      int     `json:"to tal_tokens"`
		TotalTime        float64 `json:"total_time"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
	XGroq             struct {
		Id string `json:"id"`
	} `json:"x_groq"`
}

// Chat sends a request to the chat endpoint
func (c *Client) Chat(req ChatRequest) (*ChatResponse, error) {
	request, err := c.newChatReq(req)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	done, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer done.Body.Close()
	resp, err := c.parseChatResp(done)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}
	return resp, nil
}

type ChatStreamResponse struct {
}

func (c *Client) ChatStream(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	req.Stream = true
	request, err := c.newChatReq(req)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	done, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer done.Body.Close()
	resp, err := c.parseChatResp(done)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}
	return resp, nil
}

var (
	EMPTY_MESSAGES_LIMIT = uint(10)
)

// doStreamRequest sends a request to the chat endpoint and streams the response back
func doStreamRequest[T streamable](client *Client, req *http.Request) (*streamReader[T], error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.client.Do(req) //nolint:bodyclose // body is closed in stream.Close()
	if err != nil {
		return new(streamReader[T]), err
	}
	if isFailureStatusCode(resp) {
		return new(streamReader[T]), client.handleErrorResp(resp)
	}
	return &streamReader[T]{
		emptyMessagesLimit: EMPTY_MESSAGES_LIMIT,
		reader:             bufio.NewReader(resp.Body),
		response:           resp,
		errAccumulator:     internal.NewErrorAccumulator(),
		unmarshaler:        &internal.JSONUnmarshaler{},
		httpHeader:         httpHeader(resp.Header),
	}, nil
}

// isFailureStatusCode checks if the status code is a failure status code
func isFailureStatusCode(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

// handleErrorResp handles error responses from the API
func (c *Client) handleErrorResp(resp *http.Response) error {
	var errRes ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errRes)
	if err != nil || errRes.Error == nil {
		reqErr := &RequestError{
			HTTPStatusCode: resp.StatusCode,
			Err:            err,
		}
		if errRes.Error != nil {
			reqErr.Err = errRes.Error
		}
		return reqErr
	}

	errRes.Error.HTTPStatusCode = resp.StatusCode
	return errRes.Error
}
