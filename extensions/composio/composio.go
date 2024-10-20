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
	defaultBaseURL = "https://backend.composio.dev/api/v2"
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
)

func NewComposer(apiKey string, opts ...ComposerOption) (*Composio, error) {
	c := &Composio{
		apiKey: apiKey,
		header: builders.Header{SetCommonHeaders: func(req *http.Request) {
			req.Header.Set("X-API-Key", apiKey)
		}},
		baseURL: defaultBaseURL,
		client:  http.DefaultClient,
		logger:  slog.Default(),
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.client == nil {
		c.client = &http.Client{}
	}
	if c.logger == nil {
		c.logger = slog.Default()
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
func (c *Composio) GetTools() ([]groq.Tool, error) {
	req, err := builders.NewRequest(
		context.Background(),
		c.header,
		http.MethodGet,
		fmt.Sprintf("%s/actions", c.baseURL),
		builders.WithBody(nil),
	)
	if err != nil {
		return nil, err
	}
	var tools []groq.Tool
	err = c.doRequest(req, &tools)
	if err != nil {
		return nil, err
	}
	return tools, nil
}
