package jigsawstack

import (
	"context"
	"net/http"
	"time"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	predictEndpoint Endpoint = "v1/ai/prediction"
)

type (
	// DatasetEntry represents a dataset entry.
	DatasetEntry struct {
		Date  time.Time `json:"date"`
		Value float64   `json:"value"`
	}
	// PredictRequest represents a request structure for prediction API.
	PredictRequest struct {
		Dataset []DatasetEntry `json:"dataset"`
	}
	// PredictResponse represents a response structure for prediction API.
	PredictResponse struct {
		Success bool           `json:"success"`
		Answer  []DatasetEntry `json:"answer"`
	}
)

// Predict predicts the text.
//
// Max text character is 5000.
func (j *JigsawStack) Predict(
	ctx context.Context,
	request PredictRequest,
) (response PredictResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(predictEndpoint),
		builders.WithBody(request),
	)
	if err != nil {
		return
	}
	var resp PredictResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}
