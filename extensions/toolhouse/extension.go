// Package toolhouse provides a Toolhouse extension for groq-go.
package toolhouse

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/builders"
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
		logger   *slog.Logger
		header   builders.Header
	}

	// Options is a function that sets options for a Toolhouse extension.
	Options func(*Extension)
)

// NewExtension creates a new Toolhouse extension.
func NewExtension(apiKey string, opts ...Options) (e *Extension, err error) {
	e = &Extension{
		apiKey:   apiKey,
		baseURL:  defaultBaseURL,
		client:   http.DefaultClient,
		bundle:   "default",
		provider: "openai",
		logger:   slog.Default(),
	}
	e.header.SetCommonHeaders = func(req *http.Request) {
		req.Header.Set("User-Agent", "Toolhouse/1.2.1 Python/3.11.0")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.apiKey))
		req.Header.Set("Content-Type", applicationJSON)
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
