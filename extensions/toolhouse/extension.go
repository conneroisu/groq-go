// Package toolhouse provides a Toolhouse extension for groq-go.
package toolhouse

import (
	"context"

	"github.com/conneroisu/groq-go"
)

// Extension is a Toolhouse extension.
type Extension struct {
	apiKey  string
	baseURL string
}

// Options is a function that sets options for a Toolhouse extension.
type Options func(*Extension)

// WithBaseURL sets the base URL for the Toolhouse extension.
func WithBaseURL(baseURL string) Options {
	return func(e *Extension) {
		e.baseURL = baseURL
	}
}

// NewExtension creates a new Toolhouse extension.
func NewExtension(apiKey string, opts ...Options) (e *Extension) {
	e.apiKey = apiKey
	for _, opt := range opts {
		opt(e)
	}
	if e.baseURL == "" {
		e.baseURL = "https://api.toolhouse.ai/v1"
	}
	return e
}

// Run runs the extension on the given history.
func (e *Extension) Run(
	ctx context.Context,
	history []groq.ChatCompletionMessage,
) ([]groq.ChatCompletionMessage, error) {
	return history, nil
}

// GetTools returns a list of tools that the extension can use.
func (e *Extension) GetTools(
	ctx context.Context,
	params GetToolsParams,
) ([]groq.Tool, error) {
	return []groq.Tool{}, nil
}

type GetToolsParams struct {
	Provider string `json:"provider,omitempty"`
	Metadata string `json:"metadata,omitempty"`
	Bundle   string `json:"bundle,omitempty"`
}

type getToolsResponse struct {
}
