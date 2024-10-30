package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	textToSqlEndpoint = "v1/ai/sql"
)

type (
	TextToSqlRequest struct {
		Prompt    string `json:"prompt"`
		SQLSchema string `json:"sql_schema"`
	}
	TextToSqlResponse struct {
		Success bool   `json:"success"`
		SQL     string `json:"sql"`
	}
)

// TextToSQL converts the text to SQL.
//
// Max text character is 5000.
func (j *JigsawStack) TextToSQL(
	ctx context.Context,
	request TextToSqlRequest,
) (response TextToSqlResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+textToSqlEndpoint,
		builders.WithBody(request),
	)
	if err != nil {
		return
	}
	var resp TextToSqlResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}
