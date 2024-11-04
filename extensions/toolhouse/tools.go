package toolhouse

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/tools"
)

// GetTools returns a list of tools that the extension can use.
func (e *Toolhouse) GetTools(
	ctx context.Context,
) ([]tools.Tool, error) {
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
	var tooling []tools.Tool
	err = e.sendRequest(req, &tooling)
	if err != nil {
		return nil, err
	}
	return tooling, nil
}
