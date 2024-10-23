package composio

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/builders"
)

type (
	// ToolsParams represents the parameters for the tools request.
	ToolsParams struct {
		App      string `url:"appNames"`
		Tags     string `url:"tags"`
		EntityID string `url:"user_uuid"`
		UseCase  string `url:"useCase"`
	}

	// Tools is a map of tools.
	Tools map[string]Tool
	// Tool represents a composio tool.
	Tool struct {
		groqTool    groq.Tool
		Enum        string       `json:"enum"`
		Tags        []string     `json:"tags"`
		Logo        string       `json:"logo"`
		AppID       string       `json:"appId"`
		AppName     string       `json:"appName"`
		DisplayName string       `json:"displayName"`
		Response    ToolResponse `json:"response"`
		Deprecated  bool         `json:"deprecated"`
	}
	// ToolResponse represents the response for a tool.
	ToolResponse struct {
		Response struct {
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
	}
)

// GetTools returns the tools for the composio client.
func (c *Composio) GetTools(params ToolsParams) ([]Tool, error) {
	ul := fmt.Sprintf("%s/actions", c.baseURL)
	u, err := url.Parse(ul)
	if err != nil {
		return nil, err
	}
	ps := url.Values{}
	if params.App != "" {
		ps.Add("appNames", params.App)
	}
	if params.Tags != "" {
		ps.Add("tags", params.Tags)
	}
	if params.EntityID != "" {
		ps.Add("user_uuid", params.EntityID)
	}
	if params.UseCase != "" {
		ps.Add("useCase", params.UseCase)
	}
	u.RawQuery = ps.Encode()
	uuuu := u.String()
	c.logger.Debug("tools", "url", uuuu)
	req, err := builders.NewRequest(
		context.Background(),
		c.header,
		http.MethodGet,
		uuuu,
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
	for _, tool := range items.Tools {
		c.tools[tool.groqTool.Function.Name] = tool
	}
	return items.Tools, nil
}
