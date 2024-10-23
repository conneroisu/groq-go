package composio

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/builders"
)

var _ groq.Tool = &Tool{}

type (
	// Tool represents a composio tool as returned by the api.
	Tool struct {
		Name        string                  `json:"name"`
		Enum        string                  `json:"enum"`
		Tags        []string                `json:"tags"`
		Logo        string                  `json:"logo"`
		AppID       string                  `json:"appId"`
		AppName     string                  `json:"appName"`
		DisplayName string                  `json:"displayName"`
		Description string                  `json:"description"`
		Parameters  groq.FunctionParameters `json:"parameters"`
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
) ([]groq.Tool, error) {
	uri := fmt.Sprintf("%s/actions", c.baseURL)
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	ps := url.Values{}
	for _, opt := range opts {
		opt(u)
	}
	u.RawQuery = ps.Encode()
	uri = u.String()
	c.logger.Debug("tools", "url", uri)

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
func groqTools(tools []Tool) []groq.Tool {
	groqTools := make([]groq.Tool, 0, len(tools))
	for _, tool := range tools {
		groqTools = append(groqTools, &tool)
	}
	return groqTools
}

// Function returns the function definition of the tool.
func (t *Tool) Function() groq.FunctionDefinition {
	return groq.FunctionDefinition{
		Name:        t.Name,
		Description: t.Description,
		Parameters:  t.Parameters,
	}
}

// WithTags sets the tags for the tools request.
func WithTags(tags ...string) ToolsOption {
	return func(u *url.URL) {
		ps := u.Query()
		ps.Add("tags", strings.Join(tags, ","))
		u.RawQuery = ps.Encode()
	}
}

// WithApp sets the app for the tools request.
func WithApp(app string) ToolsOption {
	return func(u *url.URL) {
		ps := u.Query()
		ps.Add("appNames", app)
		u.RawQuery = ps.Encode()
	}
}

// WithEntityID sets the entity id for the tools request.
func WithEntityID(entityID string) ToolsOption {
	return func(u *url.URL) {
		ps := u.Query()
		ps.Add("user_uuid", entityID)
		u.RawQuery = ps.Encode()
	}
}

// WithUseCase sets the use case for the tools request.
func WithUseCase(useCase string) ToolsOption {
	return func(u *url.URL) {
		ps := u.Query()
		ps.Add("useCase", useCase)
		u.RawQuery = ps.Encode()
	}
}
