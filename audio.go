package groq

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Whisper Defines the models provided by OpenAI to use when processing audio with OpenAI.
const (
	Whisper1 = "whisper-1"
	// AudioResponseFormatJSON is the JSON format.
	AudioResponseFormatJSON AudioResponseFormat = "json"
	// AudioResponseFormatText is the text format.
	AudioResponseFormatText AudioResponseFormat = "text"
	// AudioResponseFormatSRT is the SRT format.
	AudioResponseFormatSRT AudioResponseFormat = "srt"
	// AudioResponseFormatVerboseJSON is the verbose JSON format.
	AudioResponseFormatVerboseJSON AudioResponseFormat = "verbose_json"
	// AudioResponseFormatVTT is the VTT format.
	AudioResponseFormatVTT AudioResponseFormat = "vtt"
	// TranscriptionTimestampGranularityWord is the word timestamp granularity.
	TranscriptionTimestampGranularityWord TranscriptionTimestampGranularity = "word"
	// TranscriptionTimestampGranularitySegment is the segment timestamp granularity.
	TranscriptionTimestampGranularitySegment TranscriptionTimestampGranularity = "segment"
)

// AudioResponseFormat is the response format for the audio API.
//
// Response formats; Whisper uses AudioResponseFormatJSON by default.
// string
type AudioResponseFormat string

// TranscriptionTimestampGranularity is the timestamp granularity for the transcription.
// string
type TranscriptionTimestampGranularity string

// AudioRequest represents a request structure for audio API.
type AudioRequest struct {
	Model string // Model is the model to use for the transcription.

	FilePath string // FilePath is either an existing file in your filesystem or a filename representing the contents of Reader.

	Reader io.Reader // Reader is an optional io.Reader when you do not want to use an existing file.

	Prompt                 string                              // Prompt is the prompt for the transcription.
	Temperature            float32                             // Temperature is the temperature for the transcription.
	Language               string                              // Language is the language for the transcription. Only for transcription.
	Format                 AudioResponseFormat                 // Format is the format for the response.
	TimestampGranularities []TranscriptionTimestampGranularity // Only for transcription. TimestampGranularities is the timestamp granularities for the transcription.
}

// AudioResponse represents a response structure for audio API.
type AudioResponse struct {
	Task     string  `json:"task"`
	Language string  `json:"language"`
	Duration float64 `json:"duration"`
	Segments []struct {
		ID               int     `json:"id"`
		Seek             int     `json:"seek"`
		Start            float64 `json:"start"`
		End              float64 `json:"end"`
		Text             string  `json:"text"`
		Tokens           []int   `json:"tokens"`
		Temperature      float64 `json:"temperature"`
		AvgLogprob       float64 `json:"avg_logprob"`
		CompressionRatio float64 `json:"compression_ratio"`
		NoSpeechProb     float64 `json:"no_speech_prob"`
		Transient        bool    `json:"transient"`
	} `json:"segments"`
	Words []struct {
		Word  string  `json:"word"`
		Start float64 `json:"start"`
		End   float64 `json:"end"`
	} `json:"words"`
	Text string `json:"text"`

	http.Header
}

// SetHeader sets the header of the response.
func (r AudioResponse) SetHeader(header http.Header) {
	r.Header = header
}

// audioTextResponse is the response structure for the audio API when the response format is text.
type audioTextResponse struct {
	Text string `json:"text"`

	http.Header
}

func (r *audioTextResponse) SetHeader(header http.Header) {
	r.Header = header
}

// ToAudioResponse converts the audio text response to an audio response.
func (r *audioTextResponse) ToAudioResponse() AudioResponse {
	return AudioResponse{
		Text:   r.Text,
		Header: r.Header,
	}
}

// CreateTranscription — API call to create a transcription. Returns transcribed text.
func (c *Client) CreateTranscription(
	ctx context.Context,
	request AudioRequest,
) (response AudioResponse, err error) {
	return c.callAudioAPI(ctx, request, "transcriptions")
}

// CreateTranslation — API call to translate audio into English.
func (c *Client) CreateTranslation(
	ctx context.Context,
	request AudioRequest,
) (response AudioResponse, err error) {
	return c.callAudioAPI(ctx, request, "translations")
}

// callAudioAPI — API call to an audio endpoint.
func (c *Client) callAudioAPI(
	ctx context.Context,
	request AudioRequest,
	endpointSuffix string,
) (response AudioResponse, err error) {
	var formBody bytes.Buffer
	builder := c.createFormBuilder(&formBody)

	if err = audioMultipartForm(request, builder); err != nil {
		return AudioResponse{}, err
	}

	urlSuffix := fmt.Sprintf("/audio/%s", endpointSuffix)
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix, withModel(request.Model)),
		withBody(&formBody),
		withContentType(builder.FormDataContentType()),
	)
	if err != nil {
		return AudioResponse{}, err
	}

	if request.HasJSONResponse() {
		err = c.sendRequest(req, &response)
	} else {
		var textResponse audioTextResponse
		err = c.sendRequest(req, &textResponse)
		response = textResponse.ToAudioResponse()
	}
	if err != nil {
		return AudioResponse{}, err
	}
	return
}

// HasJSONResponse returns true if the response format is JSON.
func (r AudioRequest) HasJSONResponse() bool {
	return r.Format == "" || r.Format == AudioResponseFormatJSON ||
		r.Format == AudioResponseFormatVerboseJSON
}

// audioMultipartForm creates a form with audio file contents and the name of the model to use for
// audio processing.
func audioMultipartForm(request AudioRequest, b FormBuilder) error {
	err := createFileField(request, b)
	if err != nil {
		return err
	}

	err = b.WriteField("model", request.Model)
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

	if len(request.TimestampGranularities) > 0 {
		for _, tg := range request.TimestampGranularities {
			err = b.WriteField("timestamp_granularities[]", string(tg))
			if err != nil {
				return fmt.Errorf("writing timestamp_granularities[]: %w", err)
			}
		}
	}

	// Close the multipart writer
	return b.Close()
}

// createFileField creates the "file" form field from either an existing file or by using the reader.
func createFileField(request AudioRequest, b FormBuilder) error {
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
