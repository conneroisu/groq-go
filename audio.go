package groq

import (
	"bytes"
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

// CreateTranscription calls the transcriptions endpoint with the given request.
//
// Returns transcribed text in the response_format specified in the request.
func (c *Client) CreateTranscription(
	ctx context.Context,
	request AudioRequest,
) (AudioResponse, error) {
	return c.callAudioAPI(ctx, request, transcriptionsSuffix)
}

// CreateTranslation calls the translations endpoint with the given request.
//
// Returns the translated text in the response_format specified in the request.
func (c *Client) CreateTranslation(
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
	err = AudioMultipartForm(request, c.requestFormBuilder)
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
