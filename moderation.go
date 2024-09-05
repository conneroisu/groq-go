package groq

import (
	"context"
	"net/http"
	"strings"
)

// HarmfulCategory is the category of the harmful content.
//
// Llama Guard 3 provides content safety support for the following languages: English, French, German, Hindi, Italian, Portuguese, Spanish, and Thai.
//
// string
type HarmfulCategory string

const (
	// ModerationTextStable is the text-moderation-stable model.
	ModerationTextStable = "text-moderation-stable"
	// ModerationTextLatest is the text-moderation-latest model.
	ModerationTextLatest = "text-moderation-latest"

	// ViolentCrimes (S1) is the violent crimes category. Responses that enable, encourage, or endorse the commission of violent crimes, including: (1) unlawful violence toward people (ex: terrorism, genocide, murder, hate-motivated violent crimes, child abuse, assault, battery, kidnapping) and (2) unlawful violence toward animals (ex: animal abuse)
	ViolentCrimes HarmfulCategory = "violent_crimes"
	// NonViolentCrimes (S2) is the non-violent crimes category.
	NonViolentCrimes HarmfulCategory = "non_violent_crimes"
	// SexRelatedCrimes (S3) is the sex-related crimes category.
	SexRelatedCrimes HarmfulCategory = "sex_related_crimes"
	// ChildSexualExploitation (S4) is the child sexual exploitation category. Responses that contain, describe, enable, encourage, or endorse the sexual abuse of children.
	ChildSexualExploitation HarmfulCategory = "child_sexual_exploitation"
	// Defamation (S5) is the defamation category. Responses that contain, describe, enable, encourage, or endorse defamation.
	Defamation HarmfulCategory = "defamation"
	// SpecializedAdvice (S6) is the specialized advice category. Responses that contain, describe, enable, encourage, or endorse specialized advice.
	SpecializedAdvice HarmfulCategory = "specialized_advice"
	// Privacy (S7) is the privacy category. Responses that contain, describe, enable, encourage, or endorse privacy.
	Privacy HarmfulCategory = "privacy"
	// IntellectualProperty (S8) is the intellectual property category. Responses that contain, describe, enable, encourage, or endorse intellectual property.
	IntellectualProperty HarmfulCategory = "intellectual_property"
	// IndiscriminateWeapons (S9) is the indiscriminate weapons category. Responses that contain, describe, enable, encourage, or endorse indiscriminate weapons.
	IndiscriminateWeapons HarmfulCategory = "indiscriminate_weapons"
	// Hate (S10) is the hate category. Responses that contain, describe, enable, encourage, or endorse hate.
	Hate HarmfulCategory = "hate"
	// SuicideNSelfHarm (S11) is the suicide/self-harm category. Responses that contain, describe, enable, encourage, or endorse suicide or self-harm.
	SuicideNSelfHarm HarmfulCategory = "suicide_and_self_harm"
	// SexualContent (S12) is the sexual content category. Responses that contain, describe, enable, encourage, or endorse sexual content.
	SexualContent HarmfulCategory = "sexual_content"
	// Elections (S13) is the elections category. Responses that contain factually incorrect information about electoral systems and processes, including in the time, place, or manner of voting in civic elections.
	Elections HarmfulCategory = "elections"
	// CodeInterpreterAbuse (S14) is the code interpreter abuse category. Responses that contain, describe, enable, encourage, or endorse code interpreter abuse.
	CodeInterpreterAbuse HarmfulCategory = "code_interpreter_abuse"
)

var (
	cateMap = map[string]HarmfulCategory{
		"S1":  ViolentCrimes,
		"S2":  NonViolentCrimes,
		"S3":  SexRelatedCrimes,
		"S4":  ChildSexualExploitation,
		"S5":  Defamation,
		"S6":  SpecializedAdvice,
		"S7":  Privacy,
		"S8":  IntellectualProperty,
		"S9":  IndiscriminateWeapons,
		"S10": Hate,
		"S11": SuicideNSelfHarm,
		"S12": SexualContent,
		"S13": Elections,
		"S14": CodeInterpreterAbuse,
	}
)

var validModerationModel = map[Model]struct{}{
	ModerationTextStable: {},
	ModerationTextLatest: {},
}

// ModerationRequest represents a request structure for moderation API.
type ModerationRequest struct {
	Input string `json:"input,omitempty"` // Input is the input text to be moderated.
	Model Model  `json:"model,omitempty"` // Model is the model to use for the moderation.
}

// Moderation represents one of possible moderation results.
type Moderation struct {
	Categories HarmfulCategory `json:"categories"` // Categories is the categories of the result.
	Flagged    bool            `json:"flagged"`    // Flagged is the flagged of the result.
}

// Moderate — perform a moderation api call over a string.
// Input can be an array or slice but a string will reduce the complexity.
func (c *Client) Moderate(
	ctx context.Context,
	request ModerationRequest,
) (response Moderation, err error) {
	if _, ok := validModerationModel[request.Model]; len(request.Model) > 0 &&
		!ok {
		err = ErrChatCompletionInvalidModel{Model: request.Model}
		return
	}
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(chatCompletionsSuffix, withModel(request.Model)),
		withBody(&request),
	)
	if err != nil {
		return
	}
	var resp ChatCompletionResponse
	err = c.sendRequest(req, &resp)
	if err != nil {
		return
	}
	if strings.Contains(resp.Choices[0].Message.Content, "unsafe") {
		split := strings.Split(resp.Choices[0].Message.Content, "\n")[1]
		response.Categories = cateMap[strings.TrimSpace(split)]
		response.Flagged = true
	}
	return
}
