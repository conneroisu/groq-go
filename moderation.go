package groq

import (
	"context"
	"fmt"
	"net/http"
)

// The moderation endpoint is a tool you can use to check whether content complies with OpenAI's usage policies.
// Developers can thus identify content that our usage policies prohibits and take action, for instance by filtering it.

// The default is text-moderation-latest which will be automatically upgraded over time.
// This ensures you are always using our most accurate model.
// If you use text-moderation-stable, we will provide advanced notice before updating the model.
// Accuracy of text-moderation-stable may be slightly lower than for text-moderation-latest.
const (
	ModerationTextStable = "text-moderation-stable"
	ModerationTextLatest = "text-moderation-latest"
)

// ErrModerationInvalidModel is returned when the model is not supported with the moderation endpoint.
type ErrModerationInvalidModel struct {
	Model string
}

// Error implements the error interface.
func (e ErrModerationInvalidModel) Error() string {
	return fmt.Sprintf(
		"this model (%s) is not supported with moderation, please use text-moderation-stable or text-moderation-latest instead",
		e.Model,
	)
}

var validModerationModel = map[string]struct{}{
	ModerationTextStable: {},
	ModerationTextLatest: {},
}

// ModerationRequest represents a request structure for moderation API.
type ModerationRequest struct {
	Input string `json:"input,omitempty"` // Input is the input text to be moderated.
	Model string `json:"model,omitempty"` // Model is the model to use for the moderation.
}

// Result represents one of possible moderation results.
type Result struct {
	Categories     ResultCategories     `json:"categories"`      // Categories is the categories of the result.
	CategoryScores ResultCategoryScores `json:"category_scores"` // CategoryScores is the category scores of the result.
	Flagged        bool                 `json:"flagged"`         // Flagged is the flagged of the result.
}

// Hate represents a hate message.
type Hate struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity,omitempty"`
}

// SelfHarm represents a self-harm message.
type SelfHarm struct {
	Filtered bool   `json:"filtered"`           // Filtered is the filtered of the self-harm message.
	Severity string `json:"severity,omitempty"` // Severity is the severity of the self-harm message.
}

// Sexual represents a sexual message.
type Sexual struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity,omitempty"`
}

// Violence represents a violence message.
type Violence struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity,omitempty"`
}

// ContentFilterResults represents the content filter results.
type ContentFilterResults struct {
	Hate     Hate     `json:"hate,omitempty"`
	SelfHarm SelfHarm `json:"self_harm,omitempty"`
	Sexual   Sexual   `json:"sexual,omitempty"`
	Violence Violence `json:"violence,omitempty"`
}

// ResultCategories represents Categories of Result.
type ResultCategories struct {
	Hate                  bool `json:"hate"`
	HateThreatening       bool `json:"hate/threatening"`
	Harassment            bool `json:"harassment"`
	HarassmentThreatening bool `json:"harassment/threatening"`
	SelfHarm              bool `json:"self-harm"`
	SelfHarmIntent        bool `json:"self-harm/intent"`
	SelfHarmInstructions  bool `json:"self-harm/instructions"`
	Sexual                bool `json:"sexual"`
	SexualMinors          bool `json:"sexual/minors"`
	Violence              bool `json:"violence"`
	ViolenceGraphic       bool `json:"violence/graphic"`
}

// ResultCategoryScores represents CategoryScores of Result.
type ResultCategoryScores struct {
	Hate                  float32 `json:"hate"`
	HateThreatening       float32 `json:"hate/threatening"`
	Harassment            float32 `json:"harassment"`
	HarassmentThreatening float32 `json:"harassment/threatening"`
	SelfHarm              float32 `json:"self-harm"`
	SelfHarmIntent        float32 `json:"self-harm/intent"`
	SelfHarmInstructions  float32 `json:"self-harm/instructions"`
	Sexual                float32 `json:"sexual"`
	SexualMinors          float32 `json:"sexual/minors"`
	Violence              float32 `json:"violence"`
	ViolenceGraphic       float32 `json:"violence/graphic"`
}

// ModerationResponse represents a response structure for moderation API.
type ModerationResponse struct {
	ID      string   `json:"id"`      // ID is the ID of the response.
	Model   string   `json:"model"`   // Model is the model of the response.
	Results []Result `json:"results"` // Results is the results of the response.

	http.Header // Header is the header of the response.
}

// SetHeader sets the header of the response.
func (r *ModerationResponse) SetHeader(header http.Header) {
	r.Header = header
}

// Moderations — perform a moderation api call over a string.
// Input can be an array or slice but a string will reduce the complexity.
func (c *Client) Moderations(
	ctx context.Context,
	request ModerationRequest,
) (response ModerationResponse, err error) {
	if _, ok := validModerationModel[request.Model]; len(request.Model) > 0 &&
		!ok {
		err = ErrModerationInvalidModel{Model: request.Model}
		return
	}
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL("/moderations", withModel(request.Model)),
		withBody(&request),
	)
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	return
}
