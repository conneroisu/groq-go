package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	textToSQLEndpoint = "/v1/ai/sql"
)

type (
	// TextToSQLRequest represents a request structure for text to SQL API.
	TextToSQLRequest struct {
		Prompt    string `json:"prompt"`
		SQLSchema string `json:"sql_schema"`
	}
	// TextToSQLResponse represents a response structure for text to SQL API.
	TextToSQLResponse struct {
		Success bool   `json:"success"`
		SQL     string `json:"sql"`
	}
)

// TextToSQL converts the text to SQL.
//
// Max text character is 5000.
func (j *JigsawStack) TextToSQL(
	ctx context.Context,
	request TextToSQLRequest,
) (response TextToSQLResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+textToSQLEndpoint,
		builders.WithBody(request),
	)
	if err != nil {
		return
	}
	var resp TextToSQLResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}
