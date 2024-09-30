// Package toolhouse provides a Toolhouse extension for groq-go.
package toolhouse

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go"
)

const (
	getToolsEndpoint = "/get_tools"
	defaultBaseURL   = "https://api.toolhouse.ai/v1"
)

// Tool is a Toolhouse tool.
type Tool struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Extension is a Toolhouse extension.
type Extension struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// Options is a function that sets options for a Toolhouse extension.
type Options func(*Extension)

// WithBaseURL sets the base URL for the Toolhouse extension.
func WithBaseURL(baseURL string) Options {
	return func(e *Extension) {
		e.baseURL = baseURL
	}
}

// WithClient sets the client for the Toolhouse extension.
func WithClient(client *http.Client) Options {
	return func(e *Extension) {
		e.client = client
	}
}

// NewExtension creates a new Toolhouse extension.
func NewExtension(apiKey string, opts ...Options) (e *Extension) {
	e.apiKey = apiKey
	for _, opt := range opts {
		opt(e)
	}
	if e.baseURL == "" {
		e.baseURL = defaultBaseURL
	}
	return e
}

type RunOptions func(*ToolhouseRequest)

type StreamOptions func(*ToolhouseRequest)

type ToolhouseRequest struct {
	History []groq.ChatCompletionMessage
	Tools   []Tool
}

type ToolhouseResponse struct {
	Name string // Name is the name of the tool used.
}

type ToolhouseStreamResponse struct {
	History []groq.ChatCompletionStreamResponse
	Tools   []Tool
}

// Run runs the extension on the given history.
func (e *Extension) Run(
	ctx context.Context,
	resp groq.ChatCompletionResponse,
	opts ...RunOptions,
) ([]*groq.ChatCompletionMessage, error) {
	if resp.Choices[0].FinishReason != groq.FinishReasonFunctionCall {
		return resp.History, nil
	}
	// replace the existance of the function call with the tool call
	resp.History[0].FunctionCall = nil
	// resp.History[0].ToolCalls = []groq.ToolCall{
	//         {
	//                 Name: resp.History[len(resp.History)-1].FunctionCall.Name,
	//         },
	// }
	return resp.History, nil
}

// // Stream runs the extension on a stream from the groq api.
// func (e *Extension) Stream(
//         ctx context.Context,
//         resp groq.ChatCompletionStreamResponse,
// ) (groq.ChatCompletionStream, error) {
//         return []groq.ChatCompletionStreamResponse{}, nil
// }

// GetTools returns a list of tools that the extension can use.
func (e *Extension) GetTools(
	ctx context.Context,
	opts ...GetToolsOptions,
) ([]groq.Tool, error) {
	return []groq.Tool{}, nil
}

// GetToolsOptions represents the options for the GetTools method.
type GetToolsOptions struct {
	Provider string `json:"provider,omitempty"`
	Metadata string `json:"metadata,omitempty"`
	Bundle   string `json:"bundle,omitempty"`
}

type getToolsResponse struct {
}

func (e *Extension) convertTools(
	ctx context.Context,
	tools []Tool,
) (*[]groq.Tool, error) {
	return nil, nil
}
