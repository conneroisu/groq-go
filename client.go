package groq

import (
	"encoding/json"
	"io"
	"net/http"

	"log/slog"
)

// Client is a Groq api client.
type Client struct {
	groqApiKey string         // Groq API key
	client     *http.Client   // Client is the HTTP client to use
	models     *ModelResponse // Models is the list of models available to the client.
	verbosity  slog.Level     // Verbosity is the verbosity level for the client.
}

// Contains returns true if the model is in the list of models.
func (m *ModelResponse) contains(model string) bool {
	for _, m := range m.Data {
		if m.ID == model {
			return true
		}
	}
	return false
}

// GroqOpts is a function that sets options for a Groq client.
type GroqOpts func(*Client)

// WithClient sets the client for the Groq client.
func WithClient(client *http.Client) GroqOpts {
	return func(c *Client) {
		c.client = client
	}
}

// NewClient creates a new Groq client.
func NewClient(groqApiKey string, opts ...GroqOpts) (*Client, error) {
	c := &Client{
		groqApiKey: groqApiKey,
		client:     http.DefaultClient,
		verbosity:  slog.LevelError,
	}
	err := c.GetModels()
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// GetModels gets the list of models from the Groq API.
func (c *Client) GetModels() error {
	req, err := http.NewRequest("GET", "https://api.groq.com/openai/v1/models", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.groqApiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var modelsResponse ModelResponse
	err = json.Unmarshal(bodyText, &modelsResponse)
	if err != nil {
		return err
	}
	c.models = &modelsResponse
	return nil
}

type ModelResponse struct {
	Object string `json:"object"`
	Data   []struct {
		ID             string `json:"id"`
		Object         string `json:"object"`
		Created        int    `json:"created"`
		OwnedBy        string `json:"owned_by"`
		Active         bool   `json:"active"`
		ContextWindow  int    `json:"context_window,omitempty"`
		PublicApps     any    `json:"public_apps"`
		ContextWindow0 int    `json:"context_ window,omitempty"`
	} `json:"data"`
}
