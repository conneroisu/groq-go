package groq

import (
	"context"
	"net/http"
)

// GPT3 Defines the models provided by OpenAI to use when generating
// completions from OpenAI.
//
// GPT3 Models are designed for text-based tasks. For code-specific
// tasks, please refer to the Codex series of models.
const (
	completionsSuffix     = "/completions"
	chatCompletionsSuffix = "/chat/completions"

	GPT432K0613           = "gpt-4-32k-0613"
	GPT432K0314           = "gpt-4-32k-0314"
	GPT432K               = "gpt-4-32k"
	GPT40613              = "gpt-4-0613"
	GPT40314              = "gpt-4-0314"
	GPT4o                 = "gpt-4o"
	GPT4o20240513         = "gpt-4o-2024-05-13"
	GPT4o20240806         = "gpt-4o-2024-08-06"
	GPT4oLatest           = "chatgpt-4o-latest"
	GPT4oMini             = "gpt-4o-mini"
	GPT4oMini20240718     = "gpt-4o-mini-2024-07-18"
	GPT4Turbo             = "gpt-4-turbo"
	GPT4Turbo20240409     = "gpt-4-turbo-2024-04-09"
	GPT4Turbo0125         = "gpt-4-0125-preview"
	GPT4Turbo1106         = "gpt-4-1106-preview"
	GPT4TurboPreview      = "gpt-4-turbo-preview"
	GPT4VisionPreview     = "gpt-4-vision-preview"
	GPT4                  = "gpt-4"
	GPT3Dot5Turbo0125     = "gpt-3.5-turbo-0125"
	GPT3Dot5Turbo1106     = "gpt-3.5-turbo-1106"
	GPT3Dot5Turbo0613     = "gpt-3.5-turbo-0613"
	GPT3Dot5Turbo0301     = "gpt-3.5-turbo-0301"
	GPT3Dot5Turbo16K      = "gpt-3.5-turbo-16k"
	GPT3Dot5Turbo16K0613  = "gpt-3.5-turbo-16k-0613"
	GPT3Dot5Turbo         = "gpt-3.5-turbo"
	GPT3Dot5TurboInstruct = "gpt-3.5-turbo-instruct"
)

// Codex Defines the models provided by OpenAI.
// These models are designed for code-specific tasks, and use
// a different tokenizer which optimizes for whitespace.
const (
	CodexCodeDavinci002 = "code-davinci-002"
	CodexCodeCushman001 = "code-cushman-001"
	CodexCodeDavinci001 = "code-davinci-001"
)

var disabledModelsForEndpoints = map[string]map[string]bool{
	completionsSuffix: {
		GPT3Dot5Turbo:        true,
		GPT3Dot5Turbo0301:    true,
		GPT3Dot5Turbo0613:    true,
		GPT3Dot5Turbo1106:    true,
		GPT3Dot5Turbo0125:    true,
		GPT3Dot5Turbo16K:     true,
		GPT3Dot5Turbo16K0613: true,
		GPT4:                 true,
		GPT4o:                true,
		GPT4o20240513:        true,
		GPT4o20240806:        true,
		GPT4oLatest:          true,
		GPT4oMini:            true,
		GPT4oMini20240718:    true,
		GPT4TurboPreview:     true,
		GPT4VisionPreview:    true,
		GPT4Turbo1106:        true,
		GPT4Turbo0125:        true,
		GPT4Turbo:            true,
		GPT4Turbo20240409:    true,
		GPT40314:             true,
		GPT40613:             true,
		GPT432K:              true,
		GPT432K0314:          true,
		GPT432K0613:          true,
	},
	chatCompletionsSuffix: {
		CodexCodeDavinci002: true,
		CodexCodeCushman001: true,
		CodexCodeDavinci001: true,
	},
}

func endpointSupportsModel(endpoint, model string) bool {
	return !disabledModelsForEndpoints[endpoint][model]
}

func checkPromptType(prompt any) bool {
	_, isString := prompt.(string)
	_, isStringSlice := prompt.([]string)
	return isString || isStringSlice
}

// CompletionRequest represents a request structure for completion API.
type CompletionRequest struct {
	Model            string         `json:"model"`                       // Model is the model to use for the completion.
	Prompt           any            `json:"prompt,omitempty"`            // Prompt is the prompt for the completion.
	BestOf           int            `json:"best_of,omitempty"`           // BestOf is the number of completions to generate.
	Echo             bool           `json:"echo,omitempty"`              // Echo is whether to echo back the prompt in the completion.
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"` // FrequencyPenalty is the frequency penalty for the completion.
	LogitBias        map[string]int `json:"logit_bias,omitempty"`        // LogitBias is must be a token id string (specified by their token ID in the tokenizer), not a word string. incorrect: `"logit_bias":{"You": 6}`, correct: `"logit_bias":{"1639": 6}` refs: https://platform.openai.com/docs/api-reference/completions/create#completions/create-logit_bias
	LogProbs         int            `json:"logprobs,omitempty"`          // LogProbs is whether to include the log probabilities in the response.
	MaxTokens        int            `json:"max_tokens,omitempty"`        // MaxTokens is the maximum number of tokens to generate.
	N                int            `json:"n,omitempty"`                 // N is the number of completions to generate.
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`  // PresencePenalty is the presence penalty for the completion.
	Seed             *int           `json:"seed,omitempty"`              // Seed is the seed for the completion.
	Stop             []string       `json:"stop,omitempty"`              // Stop is the stop sequence for the completion.
	Stream           bool           `json:"stream,omitempty"`            // Stream is whether to stream the response.
	Suffix           string         `json:"suffix,omitempty"`            // Suffix is the suffix for the completion.
	Temperature      float32        `json:"temperature,omitempty"`       // Temperature is the temperature for the completion.
	TopP             float32        `json:"top_p,omitempty"`             // TopP is the top p for the completion.
	User             string         `json:"user,omitempty"`              // User is the user for the completion.
}

// CompletionChoice represents one of possible completions.
type CompletionChoice struct {
	Text         string        `json:"text"`          // Text is the text of the completion.
	Index        int           `json:"index"`         // Index is the index of the completion.
	FinishReason string        `json:"finish_reason"` // FinishReason is the finish reason of the completion.
	LogProbs     LogprobResult `json:"logprobs"`      // LogProbs is the log probabilities of the completion.
}

// LogprobResult represents logprob result of Choice.
type LogprobResult struct {
	Tokens        []string             `json:"tokens"`         // Tokens is the tokens of the completion.
	TokenLogprobs []float32            `json:"token_logprobs"` // TokenLogprobs is the token log probabilities of the completion.
	TopLogprobs   []map[string]float32 `json:"top_logprobs"`   // TopLogprobs is the top log probabilities of the completion.
	TextOffset    []int                `json:"text_offset"`    // TextOffset is the text offset of the completion.
}

// CompletionResponse represents a response structure for completion API.
type CompletionResponse struct {
	ID      string             `json:"id"`      // ID is the ID of the completion.
	Object  string             `json:"object"`  // Object is the object of the completion.
	Created int64              `json:"created"` // Created is the created time of the completion.
	Model   string             `json:"model"`   // Model is the model of the completion.
	Choices []CompletionChoice `json:"choices"` // Choices is the choices of the completion.
	Usage   Usage              `json:"usage"`   // Usage is the usage of the completion.

	http.Header
}

// SetHeader sets the header of the response.
func (r *CompletionResponse) SetHeader(header http.Header) {
	r.Header = header
}

// CreateCompletion — API call to create a completion. This is the main endpoint of the API. Returns new text as well
// as, if requested, the probabilities over each alternative token at each position.
//
// If using a fine-tuned model, simply provide the model's ID in the CompletionRequest object,
// and the server will use the model's parameters to generate the completion.
func (c *Client) CreateCompletion(
	ctx context.Context,
	request CompletionRequest,
) (response CompletionResponse, err error) {
	if request.Stream {
		err = ErrCompletionStreamNotSupported{}
		return
	}

	urlSuffix := "/completions"
	if !endpointSupportsModel(urlSuffix, request.Model) {
		err = ErrCompletionUnsupportedModel{}
		return
	}

	if !checkPromptType(request.Prompt) {
		err = ErrCompletionRequestPromptTypeNotSupported{}
		return
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix, withModel(request.Model)),
		withBody(request),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
