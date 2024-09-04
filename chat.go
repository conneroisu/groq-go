package groq

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	url                                                                         = "https://api.groq.com/openai/v1/chat/completions"
	ChatMessageRoleSystem                                                       = "system"         // ChatMessageRoleSystem is the system chat message role.
	ChatMessageRoleUser                                                         = "user"           // ChatMessageRoleUser is the user chat message role.
	ChatMessageRoleAssistant                                                    = "assistant"      // ChatMessageRoleAssistant is the assistant chat message role.
	ChatMessageRoleFunction                                                     = "function"       // ChatMessageRoleFunction is the function chat message role.
	ChatMessageRoleTool                                                         = "tool"           // ChatMessageRoleTool is the tool chat message role.
	ImageURLDetailHigh                         ImageURLDetail                   = "high"           // ImageURLDetailHigh is the high image url detail.
	ImageURLDetailLow                          ImageURLDetail                   = "low"            // ImageURLDetailLow is the low image url detail.
	ImageURLDetailAuto                         ImageURLDetail                   = "auto"           // ImageURLDetailAuto is the auto image url detail.
	ChatMessagePartTypeText                    ChatMessagePartType              = "text"           // ChatMessagePartTypeText is the text chat message part type.
	ChatMessagePartTypeImageURL                ChatMessagePartType              = "image_url"      // ChatMessagePartTypeImageURL is the image url chat message part type.
	ChatCompletionResponseFormatTypeJSONObject ChatCompletionResponseFormatType = "json_object"    // ChatCompletionResponseFormatTypeJSONObject is the json object chat completion response format type.
	ChatCompletionResponseFormatTypeJSONSchema ChatCompletionResponseFormatType = "json_schema"    // ChatCompletionResponseFormatTypeJSONSchema is the json schema chat completion response format type.
	ChatCompletionResponseFormatTypeText       ChatCompletionResponseFormatType = "text"           // ChatCompletionResponseFormatTypeText is the text chat completion response format type.
	ToolTypeFunction                           ToolType                         = "function"       // ToolTypeFunction is the function tool type.
	FinishReasonStop                           FinishReason                     = "stop"           // FinishReasonStop is the stop finish reason.
	FinishReasonLength                         FinishReason                     = "length"         // FinishReasonLength is the length finish reason.
	FinishReasonFunctionCall                   FinishReason                     = "function_call"  // FinishReasonFunctionCall is the function call finish reason.
	FinishReasonToolCalls                      FinishReason                     = "tool_calls"     // FinishReasonToolCalls is the tool calls finish reason.
	FinishReasonContentFilter                  FinishReason                     = "content_filter" // FinishReasonContentFilter is the content filter finish reason.
	FinishReasonNull                           FinishReason                     = "null"           // FinishReasonNull is the null finish reason.
)

// Message is a message in a chat request.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// PromptAnnotation represents the prompt annotation.
type PromptAnnotation struct {
	PromptIndex          int                  `json:"prompt_index,omitempty"`
	ContentFilterResults ContentFilterResults `json:"content_filter_results,omitempty"`
}

// ImageURLDetail is the image url detail.
// string
type ImageURLDetail string

// ChatMessageImageURL represents the chat message image url.
type ChatMessageImageURL struct {
	URL    string         `json:"url,omitempty"`    // URL is the url of the image.
	Detail ImageURLDetail `json:"detail,omitempty"` // Detail is the detail of the image url.
}

// ChatMessagePartType is the chat message part type.
// string
type ChatMessagePartType string

// ChatMessagePart represents the chat message part of a chat completion message.
type ChatMessagePart struct {
	Type     ChatMessagePartType  `json:"type,omitempty"`
	Text     string               `json:"text,omitempty"`
	ImageURL *ChatMessageImageURL `json:"image_url,omitempty"`
}

// ChatCompletionMessage represents the chat completion message.
type ChatCompletionMessage struct {
	Role         string `json:"role"`
	Content      string `json:"content"`
	MultiContent []ChatMessagePart

	// This property isn't in the official documentation, but it's in
	// the documentation for the official library for python:
	// - https://github.com/openai/openai-python/blob/main/chatml.md
	// - https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	Name string `json:"name,omitempty"`

	FunctionCall *FunctionCall `json:"function_call,omitempty"`

	// For Role=assistant prompts this may be set to the tool calls generated by the model, such as function calls.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`

	// For Role=tool prompts this should be set to the ID given in the assistant's prior request to call a tool.
	ToolCallID string `json:"tool_call_id,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface.
func (m ChatCompletionMessage) MarshalJSON() ([]byte, error) {
	if m.Content != "" && m.MultiContent != nil {
		return nil, &ErrContentFieldsMisused{field: "Content"}
	}
	if len(m.MultiContent) > 0 {
		msg := struct {
			Role         string            `json:"role"`
			Content      string            `json:"-"`
			MultiContent []ChatMessagePart `json:"content,omitempty"`
			Name         string            `json:"name,omitempty"`
			FunctionCall *FunctionCall     `json:"function_call,omitempty"`
			ToolCalls    []ToolCall        `json:"tool_calls,omitempty"`
			ToolCallID   string            `json:"tool_call_id,omitempty"`
		}(m)
		return json.Marshal(msg)
	}
	msg := struct {
		Role         string            `json:"role"`
		Content      string            `json:"content"`
		MultiContent []ChatMessagePart `json:"-"`
		Name         string            `json:"name,omitempty"`
		FunctionCall *FunctionCall     `json:"function_call,omitempty"`
		ToolCalls    []ToolCall        `json:"tool_calls,omitempty"`
		ToolCallID   string            `json:"tool_call_id,omitempty"`
	}(m)
	return json.Marshal(msg)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *ChatCompletionMessage) UnmarshalJSON(bs []byte) (err error) {
	msg := struct {
		Role         string `json:"role"`
		Content      string `json:"content"`
		MultiContent []ChatMessagePart
		Name         string        `json:"name,omitempty"`
		FunctionCall *FunctionCall `json:"function_call,omitempty"`
		ToolCalls    []ToolCall    `json:"tool_calls,omitempty"`
		ToolCallID   string        `json:"tool_call_id,omitempty"`
	}{}
	err = json.Unmarshal(bs, &msg)
	if err == nil {
		*m = ChatCompletionMessage(msg)
		return nil
	}
	multiMsg := struct {
		Role         string `json:"role"`
		Content      string
		MultiContent []ChatMessagePart `json:"content"`
		Name         string            `json:"name,omitempty"`
		FunctionCall *FunctionCall     `json:"function_call,omitempty"`
		ToolCalls    []ToolCall        `json:"tool_calls,omitempty"`
		ToolCallID   string            `json:"tool_call_id,omitempty"`
	}{}
	err = json.Unmarshal(bs, &multiMsg)
	if err != nil {
		return err
	}
	*m = ChatCompletionMessage(multiMsg)
	return nil
}

// ToolCall represents a tool call.
type ToolCall struct {
	// Index is not nil only in chat completion chunk object
	Index    *int         `json:"index,omitempty"`
	ID       string       `json:"id"`
	Type     ToolType     `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call.
type FunctionCall struct {
	Name string `json:"name,omitempty"`
	// call function with arguments in JSON format
	Arguments string `json:"arguments,omitempty"`
}

// ChatCompletionResponseFormatType is the chat completion response format type.
//
// string
type ChatCompletionResponseFormatType string

// ChatCompletionResponseFormat is the chat completion response format.
type ChatCompletionResponseFormat struct {
	Type       ChatCompletionResponseFormatType        `json:"type,omitempty"`
	JSONSchema *ChatCompletionResponseFormatJSONSchema `json:"json_schema,omitempty"`
}

// ChatCompletionResponseFormatJSONSchema is the chat completion response format json schema.
type ChatCompletionResponseFormatJSONSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Schema      json.Marshaler `json:"schema"`
	Strict      bool           `json:"strict"`
}

// ChatCompletionRequest represents a request structure for the chat completion API.
type ChatCompletionRequest struct {
	Model            string                        `json:"model"`
	Messages         []ChatCompletionMessage       `json:"messages"`
	MaxTokens        int                           `json:"max_tokens,omitempty"`
	Temperature      float32                       `json:"temperature,omitempty"`
	TopP             float32                       `json:"top_p,omitempty"`
	N                int                           `json:"n,omitempty"`
	Stream           bool                          `json:"stream,omitempty"`
	Stop             []string                      `json:"stop,omitempty"`
	PresencePenalty  float32                       `json:"presence_penalty,omitempty"`
	ResponseFormat   *ChatCompletionResponseFormat `json:"response_format,omitempty"`
	Seed             *int                          `json:"seed,omitempty"`
	FrequencyPenalty float32                       `json:"frequency_penalty,omitempty"`
	// LogitBias is must be a token id string (specified by their token ID in the tokenizer), not a word string.
	// incorrect: `"logit_bias":{"You": 6}`, correct: `"logit_bias":{"1639": 6}`
	// refs: https://platform.openai.com/docs/api-reference/chat/create#chat/create-logit_bias
	LogitBias map[string]int `json:"logit_bias,omitempty"`
	// LogProbs indicates whether to return log probabilities of the output tokens or not.
	// If true, returns the log probabilities of each output token returned in the content of message.
	// This option is currently not available on the gpt-4-vision-preview model.
	LogProbs bool `json:"logprobs,omitempty"`
	// TopLogProbs is an integer between 0 and 5 specifying the number of most likely tokens to return at each
	// token position, each with an associated log probability.
	// logprobs must be set to true if this parameter is used.
	TopLogProbs int    `json:"top_logprobs,omitempty"`
	User        string `json:"user,omitempty"`
	// Deprecated: use Tools instead.
	Functions []FunctionDefinition `json:"functions,omitempty"`
	// Deprecated: use ToolChoice instead.
	FunctionCall any    `json:"function_call,omitempty"`
	Tools        []Tool `json:"tools,omitempty"`
	// This can be either a string or an ToolChoice object.
	ToolChoice any `json:"tool_choice,omitempty"`
	// Options for streaming response. Only set this when you set stream: true.
	StreamOptions *StreamOptions `json:"stream_options,omitempty"`
	// Disable the default behavior of parallel tool calls by setting it: false.
	ParallelToolCalls any `json:"parallel_tool_calls,omitempty"`
}

// StreamOptions represents the stream options.
type StreamOptions struct {
	// If set, an additional chunk will be streamed before the data: [DONE] message.
	// The usage field on this chunk shows the token usage statistics for the entire request,
	// and the choices field will always be an empty array.
	// All other chunks will also include a usage field, but with a null value.
	IncludeUsage bool `json:"include_usage,omitempty"`
}

// ToolType is the tool type.
// string
type ToolType string

// Tool represents the tool.
type Tool struct {
	Type     ToolType            `json:"type"`
	Function *FunctionDefinition `json:"function,omitempty"`
}

// ToolChoice represents the tool choice.
type ToolChoice struct {
	Type     ToolType     `json:"type"`
	Function ToolFunction `json:"function,omitempty"`
}

// ToolFunction represents the tool function.
type ToolFunction struct {
	Name string `json:"name"`
}

// FunctionDefinition represents the function definition.
type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Strict      bool   `json:"strict,omitempty"`
	// Parameters is an object describing the function.
	// You can pass json.RawMessage to describe the schema,
	// or you can pass in a struct which serializes to the proper JSON schema.
	// The jsonschema package is provided for convenience, but you should
	// consider another specialized library if you require more complex schemas.
	Parameters any `json:"parameters"`
}

// TopLogProbs represents the top log probs.
type TopLogProbs struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []byte  `json:"bytes,omitempty"`
}

// LogProb represents the probability information for a token.
type LogProb struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []byte  `json:"bytes,omitempty"` // Omitting the field if it is null
	// TopLogProbs is a list of the most likely tokens and their log probability, at this token position.
	// In rare cases, there may be fewer than the number of requested top_logprobs returned.
	TopLogProbs []TopLogProbs `json:"top_logprobs"`
}

// LogProbs is the top-level structure containing the log probability information.
type LogProbs struct {
	// Content is a list of message content tokens with log probability information.
	Content []LogProb `json:"content"`
}

// FinishReason is the finish reason.
// string
type FinishReason string

// MarshalJSON implements the json.Marshaler interface.
func (r FinishReason) MarshalJSON() ([]byte, error) {
	if r == FinishReasonNull || r == "" {
		return []byte("null"), nil
	}
	return []byte(
		`"` + string(r) + `"`,
	), nil // best effort to not break future API changes
}

// ChatCompletionChoice represents the chat completion choice.
type ChatCompletionChoice struct {
	Index   int                   `json:"index"`
	Message ChatCompletionMessage `json:"message"`
	// FinishReason
	// stop: API returned complete message,
	// or a message terminated by one of the stop sequences provided via the stop parameter
	// length: Incomplete model output due to max_tokens parameter or token limit
	// function_call: The model decided to call a function
	// content_filter: Omitted content due to a flag from our content filters
	// null: API response still in progress or incomplete
	FinishReason FinishReason `json:"finish_reason"`      // FinishReason is the finish reason of the choice.
	LogProbs     *LogProbs    `json:"logprobs,omitempty"` // LogProbs is the log probs of the choice.
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	ID                string                 `json:"id"`                 // ID is the id of the response.
	Object            string                 `json:"object"`             // Object is the object of the response.
	Created           int64                  `json:"created"`            // Created is the created time of the response.
	Model             string                 `json:"model"`              // Model is the model of the response.
	Choices           []ChatCompletionChoice `json:"choices"`            // Choices is the choices of the response.
	Usage             Usage                  `json:"usage"`              // Usage is the usage of the response.
	SystemFingerprint string                 `json:"system_fingerprint"` // SystemFingerprint is the system fingerprint of the response.

	http.Header // Header is the header of the response.
}

// SetHeader sets the header of the response.
func (r *ChatCompletionResponse) SetHeader(h http.Header) {
	r.Header = h
}

// CreateChatCompletion is an API call to create a chat completion.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
) (response ChatCompletionResponse, err error) {
	if request.Stream {
		return response, ErrChatCompletionStreamNotSupported{model: request.Model}
	}
	if !checkEndpointSupportsModel(chatCompletionsSuffix, request.Model) {
		return response, ErrChatCompletionInvalidModel{model: request.Model}
	}
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(chatCompletionsSuffix, withModel(request.Model)),
		withBody(request))
	if err != nil {
		return response, err
	}

	return response, c.sendRequest(req, &response)
}
