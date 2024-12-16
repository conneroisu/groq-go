package groq

import (
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
		return c.CreateChatCompletion(ctx, request)
	}
	return
}

// CreateChatCompletionStream method is an API call to create a chat completion
// w/ streaming support.
func (c *Client) CreateChatCompletionStream(
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

// CreateChatCompletionJSON method is an API call to create a chat completion
// w/ object output.
func (c *Client) CreateChatCompletionJSON(
	ctx context.Context,
	request ChatCompletionRequest,
	output any,
) (err error) {
	schema, err := schema.ReflectSchema(reflect.TypeOf(output))
	if err != nil {
		return err
	}
	request.ResponseFormat = &ChatCompletionResponseFormat{
		JSONSchema: &ChatCompletionResponseFormatJSONSchema{
			Name:        schema.Title,
			Description: schema.Description,
			Schema:      *schema,
			Strict:      true,
		},
	}
	response, err := c.CreateChatCompletion(ctx, request)
	if err != nil {
		reqErr, ok := err.(*groqerr.APIError)
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
			response.ID,
			err,
		)
	}
	return
}
