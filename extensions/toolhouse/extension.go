// Package toolhouse provides a Toolhouse extension for groq-go.
package toolhouse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/conneroisu/groq-go"
)

const (
	defaultBaseURL   = "https://api.toolhouse.ai/v1"
	getToolsEndpoint = "/get_tools"
	runToolEndpoint  = "/run_tools"
	applicationJSON  = "application/json"
)

// Extension is a Toolhouse extension.
type Extension struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// LocalTool is a Toolhouse tool.
type LocalTool struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
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
func NewExtension(apiKey string, opts ...Options) (e *Extension, err error) {
	e = &Extension{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		client:  http.DefaultClient,
	}
	for _, opt := range opts {
		opt(e)
	}
	if e.apiKey == "" {
		err = fmt.Errorf("api key is required")
		return
	}
	return e, nil
}

// RunOptions are options for running a tool.
type RunOptions func(*Request)

// StreamOptions are options for streaming a tool.
type StreamOptions func(*Request)

// Request is the request to the Tool
type Request struct {
	History []groq.ChatCompletionMessage
	Tools   []Tool
}

// Response is the response from the Toolhouse API when running a tool.
type Response struct {
	Name string // Name is the name of the tool used.
}

// StreamResponse is the response from the Toolhouse API when streaming.
type StreamResponse struct {
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

// WithBundle sets the bundle for the get tools request.
func WithBundle(bundle string) GetToolsOptions {
	return func(r *getToolsRequest) {
		r.Bundle = bundle
	}
}

// WithProvider sets the provider for the get tools request.
func WithProvider(provider string) GetToolsOptions {
	return func(r *getToolsRequest) {
		r.Provider = provider
	}
}

// WithMetadata sets the metadata for the get tools request.
func WithMetadata(metadata map[string]any) GetToolsOptions {
	return func(r *getToolsRequest) {
		r.Metadata = metadata
	}
}

// GetTools returns a list of tools that the extension can use.
func (e *Extension) GetTools(
	ctx context.Context,
	opts ...GetToolsOptions,
) ([]groq.Tool, error) {
	params := getToolsRequest{
		Bundle: "default",
	}
	for _, opt := range opts {
		opt(&params)
	}
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(jsonBytes)
	url := e.baseURL + getToolsEndpoint
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		body,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Toolhouse/1.2.1 Python/3.11.0")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.apiKey))
	req.Header.Set("Content-Type", applicationJSON)
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bdy, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ts tools
	err = json.Unmarshal(bdy, &ts)
	if err != nil {
		return nil, err
	}
	return convertTools(ts)
}

type getToolsRequest struct {
	Provider string         `json:"provider"`
	Metadata map[string]any `json:"metadata"`
	Bundle   string         `json:"bundle"`
}

// GetToolsOptions represents the options for the GetTools method.
type GetToolsOptions func(*getToolsRequest)

type tools []Tool

// Tool is a Toolhouse tool.
type Tool struct {
	Type     string   `json:"type"` // Type is the type of the tool.
	Required []string `json:"required"`
	Function function `json:"function"`
}

type function struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  functionParams `json:"parameters"`
}

type functionParams struct {
	Type       string `json:"type"`
	Properties struct {
		CodeStr struct {
			Type        string `json:"type"`
			Description string `json:"description"`
		} `json:"code_str"`
	} `json:"properties"`
}

func convertTools(
	tools []Tool,
) ([]groq.Tool, error) {
	resTools := make([]groq.Tool, len(tools))
	for _, tool := range tools {
		sch, err := groq.ReflectionFromType(tool.Function.Parameters)
		if err != nil {
			return nil, err
		}
		resTools = append(resTools, groq.Tool{
			Type: groq.ToolTypeFunction,
			Function: &groq.FunctionDefinition{
				Name:        tool.Function.Name,
				Description: tool.Function.Description,
				Strict:      true,
				Parameters:  *sch,
			},
		})
	}
	return resTools, nil
}
