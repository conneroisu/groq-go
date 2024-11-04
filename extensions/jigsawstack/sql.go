package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	textToSQLEndpoint Endpoint = "/v1/ai/sql"
)

type (
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
	prompt string,
	sqlSchema string,
) (response TextToSQLResponse, err error) {
	body := struct {
		Prompt    string `json:"prompt"`
		SQLSchema string `json:"sql_schema"`
	}{
		Prompt:    prompt,
		SQLSchema: sqlSchema,
	}
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(textToSQLEndpoint),
		builders.WithBody(body),
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
