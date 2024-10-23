package composio

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	composioBaseURL = "https://backend.composio.dev/api/v1"
)

type (
	// Composio is a composio client.
	Composio struct {
		apiKey  string
		client  *http.Client
		logger  *slog.Logger
		header  builders.Header
		baseURL string
	}
	// Composer is an interface for composio.
	Composer interface {
		GetTools(opts ...ToolsOption) ([]groq.Tool, error)
		ListIntegrations() []Integration
	}
	// Integration represents a composio integration.
	Integration struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}
	// ComposerOption is an option for the composio client.
	ComposerOption func(*Composio)
	// ToolsOption is an option for the tools request.
	ToolsOption func(*url.URL)
)

// NewComposer creates a new composio client.
func NewComposer(apiKey string, opts ...ComposerOption) (*Composio, error) {
	c := &Composio{
		apiKey: apiKey,
		header: builders.Header{SetCommonHeaders: func(req *http.Request) {
			req.Header.Set("X-API-Key", apiKey)
		}},
		baseURL: composioBaseURL,
		client:  http.DefaultClient,
		logger:  slog.Default(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

func (c *Composio) doRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json")
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK ||
		res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("failed to create sandbox: %s", res.Status)
	}
	if v == nil {
		return nil
	}
	switch o := v.(type) {
	case *string:
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		*o = string(b)
		return nil
	default:
		return json.NewDecoder(res.Body).Decode(v)
	}
}

// WithLogger sets the logger for the composio client.
func WithLogger(logger *slog.Logger) ComposerOption {
	return func(c *Composio) { c.logger = logger }
}
