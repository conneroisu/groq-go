package toolhouse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/conneroisu/groq-go"
	"github.com/conneroisu/groq-go/pkg/builders"
)

// MustGetTools returns a list of tools that the extension can use.
//
// It panics if an error occurs.
func (e *Toolhouse) MustGetTools(
	ctx context.Context,
) []groq.Tool {
	tools, err := e.GetTools(ctx)
	if err != nil {
		panic(err)
	}
	return tools
}

// GetTools returns a list of tools that the extension can use.
func (e *Toolhouse) GetTools(
	ctx context.Context,
) ([]groq.Tool, error) {
	if len(e.tools) > 0 {
		return e.tools, nil
	}
	e.logger.Debug("Getting tools from Toolhouse extension")
	url := e.baseURL + getToolsEndpoint
	req, err := builders.NewRequest(
		ctx,
		e.header,
		http.MethodPost,
		url,
		builders.WithBody(
			request{
				Bundle:   "default",
				Provider: "openai",
				Metadata: e.metadata,
			}),
	)
	if err != nil {
		return nil, err
	}
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}
	bdy, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w: %s", err, string(bdy))
	}
	err = json.Unmarshal(bdy, &e.tools)
	if err != nil {
		return nil, err
	}
	return e.tools, nil
}
