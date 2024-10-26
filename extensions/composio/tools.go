package composio

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/tools"
)

type (
	// Tooler is an interface for retreiving composio tools.
	Tooler interface {
		GetTools(ctx context.Context, opts ...ToolsOption) (
			[]tools.Tool, error)
	}
	// Tool represents a composio tool as returned by the api.
	Tool struct {
		Name        string                   `json:"name"`
		Enum        string                   `json:"enum"`
		Tags        []string                 `json:"tags"`
		Logo        string                   `json:"logo"`
		AppID       string                   `json:"appId"`
		AppName     string                   `json:"appName"`
		DisplayName string                   `json:"displayName"`
		Description string                   `json:"description"`
		Parameters  tools.FunctionParameters `json:"parameters"`
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
		Deprecated   bool   `json:"deprecated"`
		DisplayName0 string `json:"display_name"`
	}
)

// GetTools returns the tools for the composio client.
func (c *Composio) GetTools(
	ctx context.Context,
	opts ...ToolsOption,
) ([]tools.Tool, error) {
	uri := fmt.Sprintf("%s/v1/actions", c.baseURL)
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for _, opt := range opts {
		opt(&q)
	}
	u.RawQuery = q.Encode()
	uri = u.String()
	c.logger.Debug("tools", "uri", uri)
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodGet,
		uri,
		builders.WithBody(nil),
	)
	if err != nil {
		return nil, err
	}
	var items struct {
		Tools []Tool `json:"items"`
	}
	err = c.doRequest(req, &items)
	if err != nil {
		return nil, err
	}
	c.logger.Debug("tools", "toolslen", len(items.Tools))
	return groqTools(items.Tools), nil
}
func groqTools(localTools []Tool) []tools.Tool {
	groqTools := make([]tools.Tool, 0, len(localTools))
	for _, tool := range localTools {
		groqTools = append(groqTools, tools.Tool{
			Function: tools.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
			},
			Type: tools.ToolTypeFunction,
		})
	}
	return groqTools
}
