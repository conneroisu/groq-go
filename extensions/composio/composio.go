package composio

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/tools"
)

const (
	composioBaseURL = "https://backend.composio.dev/api"
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
		GetTools(opts ...ToolsOption) ([]tools.Tool, error)
		ListIntegrations() []Integration
	}
	// Integration represents a composio integration.
	Integration struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}
)

// NewComposer creates a new composio client.
func NewComposer(apiKey string, opts ...Option) (*Composio, error) {
	c := &Composio{
		apiKey: apiKey,
		header: builders.Header{SetCommonHeaders: func(r *http.Request) {
			r.Header.Set("X-API-Key", apiKey)
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
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK ||
		res.StatusCode >= http.StatusBadRequest {
		bodyText, _ := io.ReadAll(res.Body)
		return fmt.Errorf("request failed: %s\nbody: %s", res.Status, bodyText)
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
		err = json.NewDecoder(res.Body).Decode(v)
		if err != nil {
			bodyText, _ := io.ReadAll(res.Body)
			return fmt.Errorf("failed to decode response: %w\nbody: %s", err, bodyText)
		}
		return nil
	}
}
