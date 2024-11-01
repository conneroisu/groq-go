package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	summaryEndpoint = "/v1/ai/summarize"
)

type (
	// SummaryRequest represents a request structure for summary API.
	SummaryRequest struct {
		Text string `json:"text"`
	}
	// SummaryResponse represents a response structure for summary API.
	SummaryResponse struct {
		Success bool   `json:"success"`
		Summary string `json:"summary"`
	}
)

// Summarize summarizes the give text.
//
// Max text character is 5000.
func (j *JigsawStack) Summarize(
	ctx context.Context,
	request SummaryRequest,
) (response SummaryResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+summaryEndpoint,
		builders.WithBody(request),
	)
	if err != nil {
		return
	}
	var resp SummaryResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}
