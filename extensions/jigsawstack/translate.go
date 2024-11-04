package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	translateEndpoint = "/v1/ai/translate"
)

type (
	// Language is a language.
	Language string
	// TranslateRequest represents a request structure for translate API.
	TranslateRequest struct {
		CurrentLanguage Language `json:"current_language"`
		TargetLanguage  Language `json:"target_language"`
		Text            string   `json:"text"`
	}
	// TranslateResponse represents a response structure for translate API.
	TranslateResponse struct {
		Success        bool   `json:"success"`
		TranslatedText string `json:"translated_text"`
	}
)

// Translate translates the text from the current language to the target language.
//
// Max text character is 5000.
func (j *JigsawStack) Translate(
	ctx context.Context,
	request TranslateRequest,
) (response TranslateResponse, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+translateEndpoint,
		builders.WithBody(request),
	)
	if err != nil {
		return
	}
	var resp TranslateResponse
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}
