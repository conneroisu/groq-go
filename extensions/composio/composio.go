package composio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	composioBaseURL = "https://backend.composio.dev/api/v2"
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
		GetTools() []groq.Tool
		ListIntegrations() []Integration
	}
	// Integration represents a composio integration.
	Integration struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}
	// ComposerOption is an option for the composio client.
	ComposerOption func(*Composio)
	// Tool represents a composio tool.
	Tool struct {
		groq.Tool
		Enum        string   `json:"enum"`
		Tags        []string `json:"tags"`
		Logo        string   `json:"logo"`
		AppID       string   `json:"appId"`
		AppName     string   `json:"appName"`
		DisplayName string   `json:"displayName"`
		Response    struct {
			Properties struct {
				Data struct {
					Title string `json:"title"`
					Type  string `json:"type"`
				} `json:"data"`
				Successful struct {
					Description string `json:"description"`
					Title       string `json:"title"`
					Type        string `json:"type"`
				} `json:"successful"`
				Error struct {
					AnyOf []struct {
						Type string `json:"type"`
					} `json:"anyOf"`
					Default     any    `json:"default"`
					Description string `json:"description"`
					Title       string `json:"title"`
				} `json:"error"`
			} `json:"properties"`
			Required []string `json:"required"`
			Title    string   `json:"title"`
			Type     string   `json:"type"`
		} `json:"response"`
		Deprecated bool `json:"deprecated"`
	}
	// ToolsParams represents the parameters for the tools request.
	ToolsParams struct {
		App      string `url:"appNames"`
		Tags     string `url:"tags"`
		EntityID string `url:"user_uuid"`
		UseCase  string `url:"useCase"`
	}
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

// GetTools returns the tools for the composio client.
func (c *Composio) GetTools(params ToolsParams) ([]Tool, error) {
	url := fmt.Sprintf("%s/actions", c.baseURL)
	if params.App != "" {
		url = fmt.Sprintf("%s?appNames=%s", url, params.App)
	}
	if params.Tags != "" {
		url = fmt.Sprintf("%s?tags=%s", url, params.Tags)
	}
	if params.EntityID != "" {
		url = fmt.Sprintf("%s?user_uuid=%s", url, params.EntityID)
	}
	if params.UseCase != "" {
		url = fmt.Sprintf("%s?useCase=%s", url, params.UseCase)
	}
	req, err := builders.NewRequest(
		context.Background(),
		c.header,
		http.MethodGet,
		url,
		builders.WithBody(nil),
	)
	if err != nil {
		return nil, err
	}
	var tools struct {
		Tools []Tool `json:"items"`
	}
	err = c.doRequest(req, &tools)
	if err != nil {
		return nil, err
	}
	c.logger.Debug("tools", "toolslen", len(tools.Tools))
	return tools.Tools, nil
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
	return func(c *Composio) {
		c.logger = logger
	}
}
