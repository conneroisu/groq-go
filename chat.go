package groq

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/tools"
)

const (
	ChatMessageRoleSystem                      Role                             = "system"         // ChatMessageRoleSystem is the system chat message role.
	ChatMessageRoleUser                        Role                             = "user"           // ChatMessageRoleUser is the user chat message role.
	ChatMessageRoleAssistant                   Role                             = "assistant"      // ChatMessageRoleAssistant is the assistant chat message role.
	ChatMessageRoleFunction                    Role                             = "function"       // ChatMessageRoleFunction is the function chat message role.
	ChatMessageRoleTool                        Role                             = "tool"           // ChatMessageRoleTool is the tool chat message role.
	ImageURLDetailHigh                         ImageURLDetail                   = "high"           // ImageURLDetailHigh is the high image url detail.
	ImageURLDetailLow                          ImageURLDetail                   = "low"            // ImageURLDetailLow is the low image url detail.
	ImageURLDetailAuto                         ImageURLDetail                   = "auto"           // ImageURLDetailAuto is the auto image url detail.
	ChatMessagePartTypeText                    ChatMessagePartType              = "text"           // ChatMessagePartTypeText is the text chat message part type.
	ChatMessagePartTypeImageURL                ChatMessagePartType              = "image_url"      // ChatMessagePartTypeImageURL is the image url chat message part type.
	ChatCompletionResponseFormatTypeJSONObject ChatCompletionResponseFormatType = "json_object"    // ChatCompletionResponseFormatTypeJSONObject is the json object chat completion response format type.
	ChatCompletionResponseFormatTypeJSONSchema ChatCompletionResponseFormatType = "json_schema"    // ChatCompletionResponseFormatTypeJSONSchema is the json schema chat completion response format type.
	ChatCompletionResponseFormatTypeText       ChatCompletionResponseFormatType = "text"           // ChatCompletionResponseFormatTypeText is the text chat completion response format type.
	FinishReasonStop                           FinishReason                     = "stop"           // FinishReasonStop is the stop finish reason.
	FinishReasonLength                         FinishReason                     = "length"         // FinishReasonLength is the length finish reason.
	FinishReasonFunctionCall                   FinishReason                     = "function_call"  // FinishReasonFunctionCall is the function call finish reason.
	FinishReasonToolCalls                      FinishReason                     = "tool_calls"     // FinishReasonToolCalls is the tool calls finish reason.
	FinishReasonContentFilter                  FinishReason                     = "content_filter" // FinishReasonContentFilter is the content filter finish reason.
	FinishReasonNull                           FinishReason                     = "null"           // FinishReasonNull is the null finish reason.
)

type (
	// ImageURLDetail is the image url detail.
	//
	// string
	ImageURLDetail string
	// ChatMessagePartType is the chat message part type.
	//
	// string
	ChatMessagePartType string
	// Role is the role of the chat completion message.
	//
	// string
	Role string
	// PromptAnnotation represents the prompt annotation.
	PromptAnnotation struct {
		PromptIndex int `json:"prompt_index,omitempty"`
	}
	// ChatMessageImageURL represents the chat message image url.
	ChatMessageImageURL struct {
		URL    string         `json:"url,omitempty"`    // URL is the url of the image.
		Detail ImageURLDetail `json:"detail,omitempty"` // Detail is the detail of the image url.
	}
	// ChatMessagePart represents the chat message part of a chat completion
	// message.
	ChatMessagePart struct {
		Text     string               `json:"text,omitempty"`      // Text is the text of the chat message part.
		Type     ChatMessagePartType  `json:"type,omitempty"`      // Type is the type of the chat message part.
		ImageURL *ChatMessageImageURL `json:"image_url,omitempty"` // ImageURL is the image url of the chat message part.
	}
	// ChatCompletionMessage represents the chat completion message.
	ChatCompletionMessage struct {
		Name         string              `json:"name"`                    // Name is the name of the chat completion message.
		Role         Role                `json:"role"`                    // Role is the role of the chat completion message.
		Content      string              `json:"content"`                 // Content is the content of the chat completion message.
		MultiContent []ChatMessagePart   `json:"-"`                       // MultiContent is the multi content of the chat completion message.
		FunctionCall *tools.FunctionCall `json:"function_call,omitempty"` // FunctionCall setting for Role=assistant prompts this may be set to the function call generated by the model.
		ToolCalls    []tools.ToolCall    `json:"tool_calls,omitempty"`    // ToolCalls setting for Role=assistant prompts this may be set to the tool calls generated by the model, such as function calls.
		ToolCallID   string              `json:"tool_call_id,omitempty"`  // ToolCallID is setting for Role=tool prompts this should be set to the ID given in the assistant's prior request to call a tool.
	}
	// ChatCompletionResponseFormatType is the chat completion response format type.
	//
	// string
	ChatCompletionResponseFormatType string
	// ChatCompletionResponseFormat is the chat completion response format.
	ChatCompletionResponseFormat struct {
		// Type is the type of the chat completion response format.
		Type ChatCompletionResponseFormatType `json:"type,omitempty"`
		// JSONSchema is the json schema of the chat completion response format.
		JSONSchema *ChatCompletionResponseFormatJSONSchema `json:"json_schema,omitempty"`
	}
	// ChatCompletionResponseFormatJSONSchema is the chat completion
	// response format json schema.
	ChatCompletionResponseFormatJSONSchema struct {
		// Name is the name of the chat completion response format json
		// schema.
		//
		// it is used to further identify the schema in the response.
		Name string `json:"name"`
		// Description is the description of the chat completion response
		// format
		// json schema.
		Description string `json:"description,omitempty"`
		// description of the chat completion response format json schema.
		// Schema is the schema of the chat completion response format json schema.
		Schema Schema `json:"schema"`
		// Strict determines whether to enforce the schema upon the generated
		// content.
		Strict bool `json:"strict"`
	}
	// ChatCompletionRequest represents a request structure for the chat completion API.
	ChatCompletionRequest struct {
		Model             ChatModel                     `json:"model"`                         // Model is the model of the chat completion request.
		Messages          []ChatCompletionMessage       `json:"messages"`                      // Messages is the messages of the chat completion request. These act as the prompt for the model.
		MaxTokens         int                           `json:"max_tokens,omitempty"`          // MaxTokens is the max tokens of the chat completion request.
		Temperature       float32                       `json:"temperature,omitempty"`         // Temperature is the temperature of the chat completion request.
		TopP              float32                       `json:"top_p,omitempty"`               // TopP is the top p of the chat completion request.
		N                 int                           `json:"n,omitempty"`                   // N is the n of the chat completion request.
		Stream            bool                          `json:"stream,omitempty"`              // Stream is the stream of the chat completion request.
		Stop              []string                      `json:"stop,omitempty"`                // Stop is the stop of the chat completion request.
		PresencePenalty   float32                       `json:"presence_penalty,omitempty"`    // PresencePenalty is the presence penalty of the chat completion request.
		ResponseFormat    *ChatCompletionResponseFormat `json:"response_format,omitempty"`     // ResponseFormat is the response format of the chat completion request.
		Seed              *int                          `json:"seed,omitempty"`                // Seed is the seed of the chat completion request.
		FrequencyPenalty  float32                       `json:"frequency_penalty,omitempty"`   // FrequencyPenalty is the frequency penalty of the chat completion request.
		LogitBias         map[string]int                `json:"logit_bias,omitempty"`          // LogitBias is must be a token id string (specified by their token ID in the tokenizer), not a word string. incorrect: `"logit_bias":{ "You": 6}`, correct: `"logit_bias":{"1639": 6}` refs: https://platform.openai.com/docs/api-reference/chat/create#chat/create-logit_bias
		LogProbs          bool                          `json:"logprobs,omitempty"`            // LogProbs indicates whether to return log probabilities of the output tokens or not. If true, returns the log probabilities of each output token returned in the content of message. This option is currently not available on the gpt-4-vision-preview model.
		TopLogProbs       int                           `json:"top_logprobs,omitempty"`        // TopLogProbs is an integer between 0 and 5 specifying the number of most likely tokens to return at each token position, each with an associated log probability. logprobs must be set to true if this parameter is used.
		User              string                        `json:"user,omitempty"`                // User is the user of the chat completion request.
		Tools             []tools.Tool                  `json:"tools,omitempty"`               // Tools is the tools of the chat completion request.
		ToolChoice        any                           `json:"tool_choice,omitempty"`         // This can be either a string or an ToolChoice object.
		StreamOptions     *StreamOptions                `json:"stream_options,omitempty"`      // Options for streaming response. Only set this when you set stream: true.
		ParallelToolCalls any                           `json:"parallel_tool_calls,omitempty"` // Disable the default behavior of parallel tool calls by setting it: false.
		RetryDelay        time.Duration                 `json:"-"`                             // RetryDelay is the delay between retries.
	}
	// LogProbs is the top-level structure containing the log probability information.
	LogProbs struct {
		Content []struct {
			Token       string        `json:"token"`           // Token is the token of the log prob.
			LogProb     float64       `json:"logprob"`         // LogProb is the log prob of the log prob.
			Bytes       []byte        `json:"bytes,omitempty"` // Omitting the field if it is null
			TopLogProbs []TopLogProbs `json:"top_logprobs"`    // TopLogProbs is a list of the most likely tokens and their log probability, at this token position. In rare cases, there may be fewer than the number of requested top_logprobs returned.
		} `json:"content"` // Content is a list of message content tokens with log probability information.
	}
	// TopLogProbs represents the top log probs.
	TopLogProbs struct {
		Token   string  `json:"token"`           // Token is the token of the top log probs.
		LogProb float64 `json:"logprob"`         // LogProb is the log prob of the top log probs.
		Bytes   []byte  `json:"bytes,omitempty"` // Bytes is the bytes of the top log probs.
	}
	// FinishReason is the finish reason.
	// string
	FinishReason string
	// ChatCompletionChoice represents the chat completion choice.
	ChatCompletionChoice struct {
		Index int `json:"index"` // Index is the index of the choice.
		// Message is the chat completion message of the choice.
		Message ChatCompletionMessage `json:"message"` // Message is the chat completion message of the choice.
		// FinishReason is the finish reason of the choice.
		//
		// stop: API returned complete message,
		// or a message terminated by one of the stop sequences provided via the stop parameter
		// length: Incomplete model output due to max_tokens parameter or token limit
		// function_call: The model decided to call a function
		// content_filter: Omitted content due to a flag from our content filters
		// null: API response still in progress or incomplete
		FinishReason FinishReason `json:"finish_reason"` // FinishReason is the finish reason of the choice.
		// LogProbs is the log probs of the choice.
		//
		// This is basically the probability of the model choosing the token.
		LogProbs *LogProbs `json:"logprobs,omitempty"` // LogProbs is the log probs of the choice.
	}
	// ChatCompletionResponse represents a response structure for chat completion API.
	ChatCompletionResponse struct {
		ID                string                 `json:"id"`                 // ID is the id of the response.
		Object            string                 `json:"object"`             // Object is the object of the response.
		Created           int64                  `json:"created"`            // Created is the created time of the response.
		Model             ChatModel              `json:"model"`              // Model is the model of the response.
		Choices           []ChatCompletionChoice `json:"choices"`            // Choices is the choices of the response.
		Usage             Usage                  `json:"usage"`              // Usage is the usage of the response.
		SystemFingerprint string                 `json:"system_fingerprint"` // SystemFingerprint is the system fingerprint of the response.
		http.Header                              // Header is the header of the response.
	}
	// ChatCompletionStreamChoiceDelta represents a response structure for chat completion API.
	ChatCompletionStreamChoiceDelta struct {
		Content      string              `json:"content,omitempty"`
		Role         string              `json:"role,omitempty"`
		FunctionCall *tools.FunctionCall `json:"function_call,omitempty"`
		ToolCalls    []tools.ToolCall    `json:"tool_calls,omitempty"`
	}
	// ChatCompletionStreamChoice represents a response structure for chat completion API.
	ChatCompletionStreamChoice struct {
		Index        int                             `json:"index"`
		Delta        ChatCompletionStreamChoiceDelta `json:"delta"`
		FinishReason FinishReason                    `json:"finish_reason"`
	}
	streamer interface {
		ChatCompletionStreamResponse
	}
	// StreamOptions represents the stream options.
	StreamOptions struct {
		// If set, an additional chunk will be streamed before the data: [DONE] message.
		// The usage field on this chunk shows the token usage statistics for the entire request,
		// and the choices field will always be an empty array.
		// All other chunks will also include a usage field, but with a null value.
		IncludeUsage bool `json:"include_usage,omitempty"`
	}
	// ChatCompletionStreamResponse represents a response structure for chat completion API.
	ChatCompletionStreamResponse struct {
		ID                  string                       `json:"id"`                           // ID is the identifier for the chat completion stream response.
		Object              string                       `json:"object"`                       // Object is the object type of the chat completion stream response.
		Created             int64                        `json:"created"`                      // Created is the creation time of the chat completion stream response.
		Model               ChatModel                    `json:"model"`                        // Model is the model used for the chat completion stream response.
		Choices             []ChatCompletionStreamChoice `json:"choices"`                      // Choices is the choices for the chat completion stream response.
		SystemFingerprint   string                       `json:"system_fingerprint"`           // SystemFingerprint is the system fingerprint for the chat completion stream response.
		PromptAnnotations   []PromptAnnotation           `json:"prompt_annotations,omitempty"` // PromptAnnotations is the prompt annotations for the chat completion stream response.
		PromptFilterResults []struct {
			Index int `json:"index"`
		} `json:"prompt_filter_results,omitempty"` // PromptFilterResults is the prompt filter results for the chat completion stream response.
		// Usage is an optional field that will only be present when you set stream_options: {"include_usage": true} in your request.
		//
		// When present, it contains a null value except for the last chunk which contains the token usage statistics
		// for the entire request.
		Usage *Usage `json:"usage,omitempty"`
	}
	// ChatCompletionStream is a stream of ChatCompletionStreamResponse.
	//
	// Note: Perhaps it is more elegant to abstract Stream using generics.
	ChatCompletionStream struct {
		*streamReader[ChatCompletionStreamResponse]
	}
)

// MarshalJSON method implements the json.Marshaler interface.
func (m ChatCompletionMessage) MarshalJSON() ([]byte, error) {
	if m.Content != "" && m.MultiContent != nil {
		return nil, &ErrContentFieldsMisused{field: "Content"}
	}
	if len(m.MultiContent) > 0 {
		msg := struct {
			Name         string              `json:"name,omitempty"`
			Role         Role                `json:"role"`
			Content      string              `json:"-"`
			MultiContent []ChatMessagePart   `json:"content,omitempty"`
			FunctionCall *tools.FunctionCall `json:"function_call,omitempty"`
			ToolCalls    []tools.ToolCall    `json:"tool_calls,omitempty"`
			ToolCallID   string              `json:"tool_call_id,omitempty"`
		}(m)
		return json.Marshal(msg)
	}
	msg := struct {
		Name         string              `json:"name,omitempty"`
		Role         Role                `json:"role"`
		Content      string              `json:"content"`
		MultiContent []ChatMessagePart   `json:"-"`
		FunctionCall *tools.FunctionCall `json:"function_call,omitempty"`
		ToolCalls    []tools.ToolCall    `json:"tool_calls,omitempty"`
		ToolCallID   string              `json:"tool_call_id,omitempty"`
	}(m)
	return json.Marshal(msg)
}

// UnmarshalJSON method implements the json.Unmarshaler interface.
func (m *ChatCompletionMessage) UnmarshalJSON(bs []byte) (err error) {
	msg := struct {
		Name         string `json:"name,omitempty"`
		Role         Role   `json:"role"`
		Content      string `json:"content"`
		MultiContent []ChatMessagePart
		FunctionCall *tools.FunctionCall `json:"function_call,omitempty"`
		ToolCalls    []tools.ToolCall    `json:"tool_calls,omitempty"`
		ToolCallID   string              `json:"tool_call_id,omitempty"`
	}{}
	err = json.Unmarshal(bs, &msg)
	if err == nil {
		*m = ChatCompletionMessage(msg)
		return nil
	}
	multiMsg := struct {
		Name         string `json:"name,omitempty"`
		Role         Role   `json:"role"`
		Content      string
		MultiContent []ChatMessagePart   `json:"content"`
		FunctionCall *tools.FunctionCall `json:"function_call,omitempty"`
		ToolCalls    []tools.ToolCall    `json:"tool_calls,omitempty"`
		ToolCallID   string              `json:"tool_call_id,omitempty"`
	}{}
	err = json.Unmarshal(bs, &multiMsg)
	if err != nil {
		return err
	}
	*m = ChatCompletionMessage(multiMsg)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (r FinishReason) MarshalJSON() ([]byte, error) {
	if r == FinishReasonNull || r == "" {
		return []byte("null"), nil
	}
	return []byte(
		`"` + string(r) + `"`,
	), nil // best effort to not break future API changes
}

// SetHeader sets the header of the response.
func (r *ChatCompletionResponse) SetHeader(h http.Header) {
	r.Header = h
}

// MustCreateChatCompletion method is an API call to create a chat completion.
//
// It panics if an error occurs.
func (c *Client) MustCreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
) (response ChatCompletionResponse) {
	response, err := c.CreateChatCompletion(ctx, request)
	if err != nil {
		panic(err)
	}
	return response
}

// CreateChatCompletion method is an API call to create a chat completion.
func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
) (response ChatCompletionResponse, err error) {
	request.Stream = false
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodPost,
		c.fullURL(chatCompletionsSuffix, withModel(model(request.Model))),
		builders.WithBody(request))
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	reqErr, ok := err.(*APIError)
	if ok && (reqErr.HTTPStatusCode == http.StatusServiceUnavailable || reqErr.HTTPStatusCode == http.StatusInternalServerError) {
		time.Sleep(request.RetryDelay)
		return c.CreateChatCompletion(ctx, request)
	}
	return
}

// CreateChatCompletionStream method is an API call to create a chat completion w/ streaming
// support.
//
// If set, tokens will be sent as data-only server-sent events as they become
// available, with the stream terminated by a data: [DONE] message.
func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request ChatCompletionRequest,
) (stream *ChatCompletionStream, err error) {
	request.Stream = true
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodPost,
		c.fullURL(chatCompletionsSuffix, withModel(model(request.Model))),
		builders.WithBody(request),
	)
	if err != nil {
		return nil, err
	}
	resp, err := sendRequestStream[ChatCompletionStreamResponse](c, req)
	if err != nil {
		return
	}
	stream = &ChatCompletionStream{
		streamReader: resp,
	}
	return
}

// CreateChatCompletionJSON method is an API call to create a chat completion w/ object output.
func (c *Client) CreateChatCompletionJSON(
	ctx context.Context,
	request ChatCompletionRequest,
	output any,
) (err error) {
	r := &reflector{}
	schema := r.ReflectFromType(reflect.TypeOf(output))
	request.ResponseFormat = &ChatCompletionResponseFormat{}
	request.ResponseFormat.JSONSchema = &ChatCompletionResponseFormatJSONSchema{
		Name:        schema.Title,
		Description: schema.Description,
		Schema:      *schema,
		Strict:      true,
	}
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodPost,
		c.fullURL(chatCompletionsSuffix, withModel(model(request.Model))),
		builders.WithBody(request),
	)
	if err != nil {
		return
	}
	var response ChatCompletionResponse
	err = c.sendRequest(req, &response)
	if err != nil {
		reqErr, ok := err.(*APIError)
		if ok && (reqErr.HTTPStatusCode == http.StatusServiceUnavailable ||
			reqErr.HTTPStatusCode == http.StatusInternalServerError) {
			time.Sleep(request.RetryDelay)
			return c.CreateChatCompletionJSON(ctx, request, output)
		}
	}
	content := response.Choices[0].Message.Content
	split := strings.Split(content, "```")
	if len(split) > 1 {
		content = split[1]
	}
	err = json.Unmarshal([]byte(content), &output)
	if err != nil {
		return fmt.Errorf(
			"error unmarshalling response (%s) to output: %v",
			response.Choices[0].Message.Content,
			err,
		)
	}
	return
}

type streamReader[T streamer] struct {
	emptyMessagesLimit uint
	isFinished         bool
	reader             *bufio.Reader
	response           *http.Response
	errAccumulator     errorAccumulator
	Header             http.Header // Header is the header of the response.
}

// Recv receives a response from the stream.
func (stream *streamReader[T]) Recv() (response T, err error) {
	if stream.isFinished {
		err = io.EOF
		return response, err
	}
	return stream.processLines()
}

// processLines processes the lines of the current response in the stream.
func (stream *streamReader[T]) processLines() (T, error) {
	var (
		headerData         = []byte("data: ")
		errorPrefix        = []byte(`data: {"error":`)
		emptyMessagesCount uint
		hasErrorPrefix     bool
	)
	for {
		rawLine, err := stream.reader.ReadBytes('\n')
		if err != nil || hasErrorPrefix {
			respErr := stream.unmarshalError()
			if respErr != nil {
				return *new(T),
					fmt.Errorf("error, %w", respErr.Error)
			}
			return *new(T), err
		}
		noSpaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasPrefix(noSpaceLine, errorPrefix) {
			hasErrorPrefix = true
		}
		if !bytes.HasPrefix(noSpaceLine, headerData) || hasErrorPrefix {
			if hasErrorPrefix {
				noSpaceLine = bytes.TrimPrefix(noSpaceLine, headerData)
			}
			err := stream.errAccumulator.Write(noSpaceLine)
			if err != nil {
				return *new(T), err
			}
			emptyMessagesCount++
			if emptyMessagesCount > stream.emptyMessagesLimit {
				return *new(T), ErrTooManyEmptyStreamMessages{}
			}
			continue
		}
		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)
		if string(noPrefixLine) == "[DONE]" {
			stream.isFinished = true
			return *new(T), io.EOF
		}
		var response T
		unmarshalErr := json.Unmarshal(noPrefixLine, &response)
		if unmarshalErr != nil {
			return *new(T), unmarshalErr
		}
		return response, nil
	}
}
func (stream *streamReader[T]) unmarshalError() (errResp *errorResponse) {
	errBytes := stream.errAccumulator.Bytes()
	if len(errBytes) == 0 {
		return
	}
	err := json.Unmarshal(errBytes, &errResp)
	if err != nil {
		errResp = nil
	}
	return
}
func (stream *streamReader[T]) Close() error {
	return stream.response.Body.Close()
}
