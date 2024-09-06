package groq

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CompletionStream is a stream of completions.
type CompletionStream struct {
	*streamReader[CompletionResponse]
}

// CreateCompletionStream — API call to create a completion w/ streaming
// support.
//
// Recv receives a response from the stream.
// It sets whether to stream back partial progress.
//
// If set, tokens will be sent as data-only server-sent events as they become
// available, with the stream terminated by a data: [DONE] message.
func (c *Client) CreateCompletionStream(
	ctx context.Context,
	request CompletionRequest,
) (*CompletionStream, error) {
	var err error
	if !endpointSupportsModel(completionsSuffix, request.Model) {
		return nil, ErrCompletionStreamNotSupported{}
	}
	if !checkPromptType(request.Prompt) {
		return nil, ErrCompletionRequestPromptTypeNotSupported{}
	}
	request.Stream = true
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(completionsSuffix, withModel(request.Model)),
		withBody(request))
	if err != nil {
		return nil, err
	}
	resp, err := sendRequestStream[CompletionResponse](c, req)
	if err != nil {
		return nil, err
	}
	return &CompletionStream{
		streamReader: resp,
	}, nil
}

type streamer interface {
	ChatCompletionStreamResponse | CompletionResponse
}

type streamReader[T streamer] struct {
	emptyMessagesLimit uint
	isFinished         bool

	reader         *bufio.Reader
	response       *http.Response
	errAccumulator errorAccumulator

	Header http.Header // Header is the header of the response.
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
		headerData  = []byte("data: ")
		errorPrefix = []byte(`data: {"error":`)

		emptyMessagesCount uint
		hasErrorPrefix     bool
	)
	for {
		rawLine, readErr := stream.reader.ReadBytes('\n')
		if readErr != nil || hasErrorPrefix {
			respErr := stream.unmarshalError()
			if respErr != nil {
				return *new(T), fmt.Errorf("error, %w", respErr.Error)
			}
			return *new(T), readErr
		}
		noSpaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasPrefix(noSpaceLine, errorPrefix) {
			hasErrorPrefix = true
		}
		if !bytes.HasPrefix(noSpaceLine, headerData) || hasErrorPrefix {
			if hasErrorPrefix {
				noSpaceLine = bytes.TrimPrefix(noSpaceLine, headerData)
			}
			writeErr := stream.errAccumulator.Write(noSpaceLine)
			if writeErr != nil {
				return *new(T), writeErr
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
