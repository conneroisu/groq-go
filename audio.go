package groq

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	// TranscriptionTimestampGranularityWord is the word timestamp
	// granularity.
	TranscriptionTimestampGranularityWord TranscriptionTimestampGranularity = "word"
	// TranscriptionTimestampGranularitySegment is the segment timestamp
	// granularity.
	TranscriptionTimestampGranularitySegment TranscriptionTimestampGranularity = "segment"
)

type (
	// TranscriptionTimestampGranularity is the timestamp granularity for
	// the transcription.
	//
	// string
	TranscriptionTimestampGranularity string
	// AudioRequest represents a request structure for audio API.
	AudioRequest struct {
		// Model is the model to use for the transcription.
		Model AudioModel
		// FilePath is either an existing file in your filesystem or a
		// filename representing the contents of Reader.
		FilePath string
		// Reader is an optional io.Reader when you do not want to use
		// an existing file.
		Reader io.Reader
		// Prompt is the prompt for the transcription.
		Prompt string
		// Temperature is the temperature for the transcription.
		Temperature float32
		// Language is the language for the transcription. Only for
		// transcription.
		Language string
		// Format is the format for the response.
		Format Format
	}
	// AudioResponse represents a response structure for audio API.
	AudioResponse struct {
		// Task is the task of the response.
		Task string `json:"task"`
		// Language is the language of the response.
		Language string `json:"language"`
		// Duration is the duration of the response.
		Duration float64 `json:"duration"`
		// Segments is the segments of the response.
		Segments Segments `json:"segments"`
		// Words is the words of the response.
		Words Words `json:"words"`
		// Text is the text of the response.
		Text string `json:"text"`

		Header http.Header // Header is the header of the response.
	}
	// Words is the words of the audio response.
	Words []struct {
		// Word is the textual representation of a word in the audio
		// response.
		Word string `json:"word"`
		// Start is the start of the words in seconds.
		Start float64 `json:"start"`
		// End is the end of the words in seconds.
		End float64 `json:"end"`
	}
	// Segments is the segments of the response.
	Segments []struct {
		// ID is the ID of the segment.
		ID int `json:"id"`
		// Seek is the seek of the segment.
		Seek int `json:"seek"`
		// Start is the start of the segment.
		Start float64 `json:"start"`
		// End is the end of the segment.
		End float64 `json:"end"`
		// Text is the text of the segment.
		Text string `json:"text"`
		// Tokens is the tokens of the segment.
		Tokens []int `json:"tokens"`
		// Temperature is the temperature of the segment.
		Temperature float64 `json:"temperature"`
		// AvgLogprob is the avg log prob of the segment.
		AvgLogprob float64 `json:"avg_logprob"`
		// CompressionRatio is the compression ratio of the segment.
		CompressionRatio float64 `json:"compression_ratio"`
		// NoSpeechProb is the no speech prob of the segment.
		NoSpeechProb float64 `json:"no_speech_prob"`
		// Transient is the transient of the segment.
		Transient bool `json:"transient"`
	}
	// audioTextResponse is the response structure for the audio API when the
	// response format is text.
	audioTextResponse struct {
		// Text is the text of the response.
		Text string `json:"text"`
		// Header is the response header.
		header http.Header `json:"-"`
	}
)

// SetHeader sets the header of the response.
func (r *AudioResponse) SetHeader(header http.Header) { r.Header = header }

// SetHeader sets the header of the audio text response.
func (r *audioTextResponse) SetHeader(header http.Header) { r.header = header }

// toAudioResponse converts the audio text response to an audio response.
func (r *audioTextResponse) toAudioResponse() AudioResponse {
	return AudioResponse{Text: r.Text, Header: r.header}
}

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

func (r AudioRequest) hasJSONResponse() bool {
	return r.Format == "" || r.Format == FormatJSON ||
		r.Format == FormatVerboseJSON
}

// AudioMultipartForm creates a form with audio file contents and the name of
// the model to use for audio processing.
func AudioMultipartForm(request AudioRequest, b builders.FormBuilder) error {
	err := createFileField(request, b)
	if err != nil {
		return err
	}
	err = b.WriteField("model", string(request.Model))
	if err != nil {
		return fmt.Errorf("writing model name: %w", err)
	}
	// Create a form field for the prompt (if provided)
	if request.Prompt != "" {
		err = b.WriteField("prompt", request.Prompt)
		if err != nil {
			return fmt.Errorf("writing prompt: %w", err)
		}
	}
	// Create a form field for the format (if provided)
	if request.Format != "" {
		err = b.WriteField("response_format", string(request.Format))
		if err != nil {
			return fmt.Errorf("writing format: %w", err)
		}
	}
	// Create a form field for the temperature (if provided)
	if request.Temperature != 0 {
		err = b.WriteField(
			"temperature",
			fmt.Sprintf("%.2f", request.Temperature),
		)
		if err != nil {
			return fmt.Errorf("writing temperature: %w", err)
		}
	}
	// Create a form field for the language (if provided)
	if request.Language != "" {
		err = b.WriteField("language", request.Language)
		if err != nil {
			return fmt.Errorf("writing language: %w", err)
		}
	}
	return b.Close()
}

func createFileField(
	request AudioRequest,
	b builders.FormBuilder,
) (err error) {
	if request.Reader != nil {
		err := b.CreateFormFileReader("file", request.Reader, request.FilePath)
		if err != nil {
			return fmt.Errorf("creating form using reader: %w", err)
		}
		return nil
	}
	f, err := os.Open(request.FilePath)
	if err != nil {
		return fmt.Errorf("opening audio file: %w", err)
	}
	defer f.Close()
	err = b.CreateFormFile("file", f)
	if err != nil {
		return fmt.Errorf("creating form file: %w", err)
	}
	return nil
}
