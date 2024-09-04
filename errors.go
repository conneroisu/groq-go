package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Helper is an interface for error helpers
type Helper interface {
	error
	Advice() string
}

// ErrorAccumulator is an interface for accumulating errors
type ErrorAccumulator interface {
	Write(p []byte) error
	Bytes() []byte
}

type errorBuffer interface {
	io.Writer
	Len() int
	Bytes() []byte
}

// DefaultErrorAccumulator is a default implementation of ErrorAccumulator
type DefaultErrorAccumulator struct {
	Buffer errorBuffer
}

// NewErrorAccumulator creates a new error accumulator
func NewErrorAccumulator() ErrorAccumulator {
	return &DefaultErrorAccumulator{
		Buffer: &bytes.Buffer{},
	}
}

// Write writes bytes to the error accumulator
func (e *DefaultErrorAccumulator) Write(p []byte) error {
	_, err := e.Buffer.Write(p)
	if err != nil {
		return fmt.Errorf("error accumulator write error, %w", err)
	}
	return nil
}

// Bytes returns the bytes of the error accumulator
func (e *DefaultErrorAccumulator) Bytes() (errBytes []byte) {
	if e.Buffer.Len() == 0 {
		return
	}
	errBytes = e.Buffer.Bytes()
	return
}

// APIError provides error information returned by the OpenAI API.
type APIError struct {
	Code           any     `json:"code,omitempty"`  // Code is the code of the error.
	Message        string  `json:"message"`         // Message is the message of the error.
	Param          *string `json:"param,omitempty"` // Param is the param of the error.
	Type           string  `json:"type"`            // Type is the type of the error.
	HTTPStatusCode int     `json:"-"`               // HTTPStatusCode is the status code of the error.
}

// RequestError provides information about generic request errors.
type RequestError struct {
	HTTPStatusCode int
	Err            error
}

// ErrorResponse is a response from the error endpoint.
type ErrorResponse struct {
	Error *APIError `json:"error,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.HTTPStatusCode > 0 {
		return fmt.Sprintf(
			"error, status code: %d, message: %s",
			e.HTTPStatusCode,
			e.Message,
		)
	}

	return e.Message
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *APIError) UnmarshalJSON(data []byte) (err error) {
	var rawMap map[string]json.RawMessage
	err = json.Unmarshal(data, &rawMap)
	if err != nil {
		return
	}

	err = json.Unmarshal(rawMap["message"], &e.Message)
	if err != nil {
		// If the parameter field of a function call is invalid as a JSON schema
		// refs: https://github.com/sashabaranov/go-openai/issues/381
		var messages []string
		err = json.Unmarshal(rawMap["message"], &messages)
		if err != nil {
			return
		}
		e.Message = strings.Join(messages, ", ")
	}

	// optional fields
	if _, ok := rawMap["param"]; ok {
		err = json.Unmarshal(rawMap["param"], &e.Param)
		if err != nil {
			return
		}
	}

	if _, ok := rawMap["code"]; !ok {
		return nil
	}

	// if the api returned a number, we need to force an integer
	// since the json package defaults to float64
	var intCode int
	err = json.Unmarshal(rawMap["code"], &intCode)
	if err == nil {
		e.Code = intCode
		return nil
	}

	return json.Unmarshal(rawMap["code"], &e.Code)
}

// Error implements the error interface.
func (e *RequestError) Error() string {
	return fmt.Sprintf(
		"error, status code: %d, message: %s",
		e.HTTPStatusCode,
		e.Err,
	)
}

// Unwrap returns the underlying error
func (e *RequestError) Unwrap() error {
	return e.Err
}

// ErrChatCompletionInvalidModel is an error that occurs when the model is not supported with the CreateChatCompletion method.
type ErrChatCompletionInvalidModel struct {
	Model    string
	Endpoint string
}

// Error implements the error interface.
func (e ErrChatCompletionInvalidModel) Error() string {
	return fmt.Errorf(
		"this model (%s) is not supported with this method of interaction over %s, please use CreateCompletion client method instead",
		e.Endpoint,
		e.Model,
	).
		Error()
}

// ErrChatCompletionStreamNotSupported is an error that occurs when streaming is not supported with the CreateChatCompletionStream method.
type ErrChatCompletionStreamNotSupported struct {
	model string
}

// Error implements the error interface.
func (e ErrChatCompletionStreamNotSupported) Error() string {
	return fmt.Errorf("streaming is not supported with this method, please use CreateChatCompletionStream client method instead").
		Error()
}

// ErrContentFieldsMisused is an error that occurs when both Content and MultiContent properties are set.
type ErrContentFieldsMisused struct {
	field string
}

// Error implements the error interface.
func (e ErrContentFieldsMisused) Error() string {
	return fmt.Errorf("can't use both Content and MultiContent properties simultaneously").
		Error()
}

// ErrCompletionUnsupportedModel is an error that occurs when the model is not supported with the CreateCompletion method.
type ErrCompletionUnsupportedModel struct{ Model string }

// Error implements the error interface.
func (e ErrCompletionUnsupportedModel) Error() string {
	return fmt.Errorf("this model (%s) is not supported with this method, please use CreateCompletion client method instead", e.Model).
		Error()
}

// ErrCompletionStreamNotSupported is an error that occurs when streaming is not supported with the CreateCompletionStream method.
type ErrCompletionStreamNotSupported struct{}

// Error implements the error interface.
func (e ErrCompletionStreamNotSupported) Error() string {
	return fmt.Errorf("streaming is not supported with this method, please use CreateCompletionStream client method instead").
		Error()
}

// ErrCompletionRequestPromptTypeNotSupported is an error that occurs when the type of CompletionRequest.Prompt only supports string and []string.
type ErrCompletionRequestPromptTypeNotSupported struct{}

// Error implements the error interface.
func (e ErrCompletionRequestPromptTypeNotSupported) Error() string {
	return fmt.Errorf("the type of CompletionRequest.Prompt only supports string and []string").
		Error()
}
