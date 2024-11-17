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
	// Tool represents a composio tool as returned by the api.
	Tool struct {
		// Name is the name of the tool returned by the composio api.
		Name string `json:"name"`
		// Enum is the enum of the tool returned by the composio api.
		Enum string `json:"enum"`
		// Tags are the tags of the tool returned by the composio api.
		Tags []string `json:"tags"`
		// Logo is the logo of the tool returned by the composio api.
		Logo string `json:"logo"`
		// AppID is the app id of the tool returned by the composio api.
		AppID string `json:"appId"`
		// AppName is the app name of the tool returned by the composio
		// api.
		AppName string `json:"appName"`
		// DisplayName is the display name of the tool returned by the
		// composio api.
		DisplayName string `json:"displayName"`
		// Description is the description of the tool returned by the
		// composio api.
		Description string `json:"description"`
		// Parameters are the parameters of the tool returned by the
		// composio api.
		Parameters tools.Parameters `json:"parameters"`
		// Response is the response of the tool returned by the
		// composio api.
		Response struct {
			// Properties are the properties of the response
			// returned by the composio api.
			Properties struct {
				// Data is the data of the response returned by
				// the composio api.
				Data struct {
					// Title is the title of the data in the
					// response returned by the composio
					// api.
					Title string `json:"title"`
					// Type is the type of the data in the
					// response returned by the composio
					// api.
					Type string `json:"type"`
				} `json:"data"`
				// Successful is the successful response of the
				// composio api.
				Successful struct {
					// Description is the description of the
					// successful response of the composio
					// api.
					Description string `json:"description"`
					// Title is the title of the successful
					// response of the composio api.
					Title string `json:"title"`
					// Type is the type of the successful
					// response of the composio api.
					Type string `json:"type"`
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
	return groqTools(items.Tools), nil
}
func groqTools(tooling []Tool) []tools.Tool {
	groqTools := make([]tools.Tool, 0, len(tooling))
	for _, tool := range tooling {
		groqTools = append(groqTools, tools.Tool{
			Function: tools.Defintion{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
			},
			Type: tools.ToolTypeFunction,
		})
	}
	return groqTools
}
