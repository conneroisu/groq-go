package composio

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

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
	return items.Tools, nil
}
