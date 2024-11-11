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

type (
	request struct {
		ConnectedAccountID string         `json:"connectedAccountId"`
		EntityID           string         `json:"entityId"`
		AppName            string         `json:"appName"`
		Input              map[string]any `json:"input"`
		Text               string         `json:"text,omitempty"`
		AuthConfig         map[string]any `json:"authConfig,omitempty"`
	}
)

// Run runs the composio client on a chat completion response.
func (c *Composio) Run(
	ctx context.Context,
	user ConnectedAccount,
	response groq.ChatCompletionResponse,
) ([]groq.ChatCompletionMessage, error) {
	var respH []groq.ChatCompletionMessage
	if response.Choices[0].FinishReason != groq.ReasonFunctionCall &&
		response.Choices[0].FinishReason != "tool_calls" {
		return nil, fmt.Errorf("not a function call")
	}
	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		var args map[string]any
		if json.Valid([]byte(toolCall.Function.Arguments)) {
			err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
			if err != nil {
				return nil, err
			}
			c.logger.Debug("arguments", "args", args)
		}
		req, err := builders.NewRequest(
			ctx,
			c.header,
			http.MethodPost,
			fmt.Sprintf("%s/v2/actions/%s/execute", c.baseURL, toolCall.Function.Name),
			builders.WithBody(&request{
				ConnectedAccountID: user.ID,
				EntityID:           "default",
				AppName:            toolCall.Function.Name,
				Input:              args,
				AuthConfig:         map[string]any{},
			}),
		)
		if err != nil {
			return nil, err
		}
		var body string
		err = c.doRequest(req, &body)
		if err != nil {
			return nil, err
		}
		respH = append(respH, groq.ChatCompletionMessage{
			Content: string(body),
			Name:    toolCall.ID,
			Role:    groq.ChatMessageRoleFunction,
		})
	}
	return respH, nil
}
