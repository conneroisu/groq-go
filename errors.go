package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type (
	// DefaultErrorAccumulator is a default implementation of ErrorAccumulator
	DefaultErrorAccumulator struct {
		Buffer errorBuffer
	}
	// APIError provides error information returned by the Groq API.
	APIError struct {
		Code           any     `json:"code,omitempty"`  // Code is the code of the error.
		Message        string  `json:"message"`         // Message is the message of the error.
		Param          *string `json:"param,omitempty"` // Param is the param of the error.
		Type           string  `json:"type"`            // Type is the type of the error.
		HTTPStatusCode int     `json:"-"`               // HTTPStatusCode is the status code of the error.
	}
	// ErrContentFieldsMisused is an error that occurs when both Content and
	// MultiContent properties are set.
	ErrContentFieldsMisused struct {
		field string
	}
	// ErrToolNotFound is returned when a tool is not found.
	ErrToolNotFound struct {
		ToolName string
	}
	// ErrTooManyEmptyStreamMessages is returned when the stream has sent too many
	// empty messages.
	ErrTooManyEmptyStreamMessages struct{}
	errorAccumulator              interface {
		// Write method writes bytes to the error accumulator
		//
		// It implements the io.Writer interface.
		Write(p []byte) error
		// Bytes method returns the bytes of the error accumulator.
		Bytes() []byte
	}
	errorBuffer interface {
		io.Writer
		Len() int
		Bytes() []byte
	}
	requestError struct {
		HTTPStatusCode int
		Err            error
	}
	errorResponse struct {
		Error *APIError `json:"error,omitempty"`
	}
)

// Error implements the error interface.
func (e ErrContentFieldsMisused) Error() string {
	return fmt.Errorf("can't use both Content and MultiContent properties simultaneously").
		Error()
}

// Error returns the error message.
func (e ErrTooManyEmptyStreamMessages) Error() string {
	return "stream has sent too many empty messages"
}

// newErrorAccumulator creates a new error accumulator
func newErrorAccumulator() errorAccumulator {
	return &DefaultErrorAccumulator{
		Buffer: &bytes.Buffer{},
	}
}

// Write method writes bytes to the error accumulator.
func (e *DefaultErrorAccumulator) Write(p []byte) error {
	_, err := e.Buffer.Write(p)
	if err != nil {
		return fmt.Errorf("error accumulator write error, %w", err)
	}
	return nil
}

// Bytes method returns the bytes of the error accumulator.
func (e *DefaultErrorAccumulator) Bytes() (errBytes []byte) {
	if e.Buffer.Len() == 0 {
		return
	}
	errBytes = e.Buffer.Bytes()
	return
}

// Error method implements the error interface on APIError.
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
func (e *requestError) Error() string {
	return fmt.Sprintf(
		"error, status code: %d, message: %s",
		e.HTTPStatusCode,
		e.Err,
	)
}

// Unwrap unwraps the error.
func (e *requestError) Unwrap() error {
	return e.Err
}

// Error implements the error interface.
func (e ErrToolNotFound) Error() string {
	return fmt.Sprintf("tool %s not found", e.ToolName)
}
