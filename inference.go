package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/groqerr"
	"github.com/conneroisu/groq-go/pkg/schema"
)

const (
	// groqAPIURLv1 is the base URL for the Groq API.
	groqAPIURLv1 = "https://api.groq.com/openai/v1"

	chatCompletionsSuffix endpoint = "/chat/completions"
	transcriptionsSuffix  endpoint = "/audio/transcriptions"
	translationsSuffix    endpoint = "/audio/translations"
	embeddingsSuffix      endpoint = "/embeddings"
	moderationsSuffix     endpoint = "/moderations"
)

// ChatCompletion method is an API call to create a chat completion.
func (c *Client) ChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
) (response ChatCompletionResponse, err error) {
	request.Stream = false
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodPost,
		c.fullURL(chatCompletionsSuffix, withModel(request.Model)),
		builders.WithBody(request))
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	reqErr, ok := err.(*groqerr.APIError)
	if ok && (reqErr.HTTPStatusCode == http.StatusServiceUnavailable ||
		reqErr.HTTPStatusCode == http.StatusInternalServerError) {
		time.Sleep(request.RetryDelay)
		return c.ChatCompletion(ctx, request)
	}
	return
}

// ChatCompletionStream method is an API call to create a chat completion
// w/ streaming support.
func (c *Client) ChatCompletionStream(
	ctx context.Context,
	request ChatCompletionRequest,
) (stream *ChatCompletionStream, err error) {
	request.Stream = true
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodPost,
		c.fullURL(
			chatCompletionsSuffix,
			withModel(request.Model)),
		builders.WithBody(request),
	)
	if err != nil {
		return nil, err
	}
	resp, err := sendRequestStream(c, req)
	if err != nil {
		return
	}
	return &ChatCompletionStream{
		StreamReader: resp,
	}, nil
}

// ChatCompletionJSON method is an API call to create a chat completion
// w/ object output.
func (c *Client) ChatCompletionJSON(
	ctx context.Context,
	request ChatCompletionRequest,
	output any,
) (err error) {
	schema, err := schema.ReflectSchema(reflect.TypeOf(output))
	if err != nil {
		return err
	}
	request.ResponseFormat = &ChatResponseFormat{
		JSONSchema: &JSONSchema{
			Name:        schema.Title,
			Description: schema.Description,
			Schema:      *schema,
			Strict:      true,
		},
		Type: FormatJSON,
	}
	response, err := c.ChatCompletion(ctx, request)
	if err != nil {
		reqErr, ok := err.(*groqerr.APIError)
		if ok && (reqErr.HTTPStatusCode == http.StatusServiceUnavailable ||
			reqErr.HTTPStatusCode == http.StatusInternalServerError) {
			time.Sleep(request.RetryDelay)
			return c.ChatCompletionJSON(ctx, request, output)
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
			response.ID,
			err,
		)
	}
	return
}

// Moderate performs a moderation api call over a string.
// Input can be an array or slice but a string will reduce the complexity.
func (c *Client) Moderate(
	ctx context.Context,
	messages []ChatCompletionMessage,
	model ModerationModel,
) (response []Moderation, err error) {
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodPost,
		c.fullURL(chatCompletionsSuffix, withModel(model)),
		builders.WithBody(&struct {
			Messages []ChatCompletionMessage `json:"messages"`
			Model    ModerationModel         `json:"model,omitempty"`
		}{
			Messages: messages,
			Model:    model,
		}),
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
		split := strings.Split(
			strings.Split(resp.Choices[0].Message.Content, "\n")[1],
			",",
		)
		for _, s := range split {
			response = append(
				response,
				sectionMap[strings.TrimSpace(s)],
			)
		}
	}
	return
}

// Transcribe calls the transcriptions endpoint with the given request.
//
// Returns transcribed text in the response_format specified in the request.
func (c *Client) Transcribe(
	ctx context.Context,
	request AudioRequest,
) (AudioResponse, error) {
	return c.callAudioAPI(ctx, request, transcriptionsSuffix)
}

// Translate calls the translations endpoint with the given request.
//
// Returns the translated text in the response_format specified in the request.
func (c *Client) Translate(
	ctx context.Context,
	request AudioRequest,
) (AudioResponse, error) {
	return c.callAudioAPI(ctx, request, translationsSuffix)
}

// callAudioAPI calls the audio API with the given request.
//
// Currently supports both the transcription and translation APIs.
func (c *Client) callAudioAPI(
	ctx context.Context,
	request AudioRequest,
	endpointSuffix endpoint,
) (response AudioResponse, err error) {
	var formBody bytes.Buffer
	c.requestFormBuilder = builders.NewFormBuilder(&formBody)
	err = audioMultipartForm(request, c.requestFormBuilder)
	if err != nil {
		return AudioResponse{}, err
	}
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodPost,
		c.fullURL(endpointSuffix, withModel(request.Model)),
		builders.WithBody(&formBody),
		builders.WithContentType(c.requestFormBuilder.FormDataContentType()),
	)
	if err != nil {
		return AudioResponse{}, err
	}

	if request.hasJSONResponse() {
		err = c.sendRequest(req, &response)
	} else {
		var textResponse audioTextResponse
		err = c.sendRequest(req, &textResponse)
		response = textResponse.toAudioResponse()
	}
	if err != nil {
		return AudioResponse{}, err
	}
	return
}
