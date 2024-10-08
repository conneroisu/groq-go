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

type (
	// Extension is a Toolhouse extension.
	Extension struct {
		apiKey   string
		baseURL  string
		client   *http.Client
		provider string
		metadata map[string]any
		bundle   string
		tools    []groq.Tool
	}

	// Options is a function that sets options for a Toolhouse extension.
	Options     func(*Extension)
	runResponse struct {
		Provider string `json:"provider"`
		Content  struct {
			Role       string `json:"role"`
			ToolCallID string `json:"tool_call_id"`
			Name       string `json:"name"`
			Content    string `json:"content"`
		} `json:"content"`
	}
)

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

// WithMetadata sets the metadata for the get tools request.
func WithMetadata(metadata map[string]any) Options {
	return func(r *Extension) {
		r.metadata = metadata
	}
}

// NewExtension creates a new Toolhouse extension.
func NewExtension(apiKey string, opts ...Options) (e *Extension, err error) {
	e = &Extension{
		apiKey:   apiKey,
		baseURL:  defaultBaseURL,
		client:   http.DefaultClient,
		bundle:   "default",
		provider: "openai",
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

// Run runs the extension on the given history.
func (e *Extension) Run(
	ctx context.Context,
	response groq.ChatCompletionResponse,
) ([]groq.ChatCompletionMessage, error) {
	if response.Choices[0].FinishReason != groq.FinishReasonFunctionCall && response.Choices[0].FinishReason != "tool_calls" {
		return nil, fmt.Errorf("Not a function call")
	}
	respH := []groq.ChatCompletionMessage{}
	for _, tool := range response.Choices[0].Message.ToolCalls {
		buf := new(bytes.Buffer)
		bodyParams := request{
			Content:  tool,
			Provider: e.provider,
			Metadata: e.metadata,
			Bundle:   e.bundle,
		}
		err := json.NewEncoder(buf).Encode(bodyParams)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("%s%s", e.baseURL, runToolEndpoint),
			bytes.NewBuffer(buf.Bytes()),
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
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%v", resp)
		}
		bdy, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var runResp runResponse
		err = json.Unmarshal(bdy, &runResp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w: %s", err, string(bdy))
		}
		cCM := groq.ChatCompletionMessage{
			Content: runResp.Content.Content,
			Name:    runResp.Content.Name,
			Role:    groq.ChatMessageRoleFunction,
		}
		respH = append(respH, cCM)
	}
	return respH, nil
}

type request struct {
	Content  groq.ToolCall  `json:"content,omitempty"`
	Provider string         `json:"provider"`
	Metadata map[string]any `json:"metadata"`
	Bundle   string         `json:"bundle"`
}

// MustGetTools returns a list of tools that the extension can use.
//
// It panics if an error occurs.
func (e *Extension) MustGetTools(
	ctx context.Context,
) []groq.Tool {
	tools, err := e.GetTools(ctx)
	if err != nil {
		panic(err)
	}
	return tools
}

// GetTools returns a list of tools that the extension can use.
func (e *Extension) GetTools(
	ctx context.Context,
) ([]groq.Tool, error) {
	params := request{
		Bundle:   "default",
		Provider: "openai",
		Metadata: e.metadata,
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
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}
	bdy, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w: %s", err, string(bdy))
	}
	err = json.Unmarshal(bdy, &e.tools)
	if err != nil {
		return nil, err
	}
	return e.tools, nil
}
