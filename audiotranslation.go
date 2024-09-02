// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package groq

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/conneroisu/groq-go/internal/apiform"
	"github.com/conneroisu/groq-go/internal/apijson"
	"github.com/conneroisu/groq-go/internal/param"
	"github.com/conneroisu/groq-go/internal/requestconfig"
	"github.com/conneroisu/groq-go/option"
)

// AudioTranslationService contains methods and other services that help with
// interacting with the groq API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAudioTranslationService] method instead.
type AudioTranslationService struct {
	Options []option.RequestOption
}

// NewAudioTranslationService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAudioTranslationService(opts ...option.RequestOption) (r *AudioTranslationService) {
	r = &AudioTranslationService{}
	r.Options = opts
	return
}

// Translates audio into English.
func (r *AudioTranslationService) New(ctx context.Context, body AudioTranslationNewParams, opts ...option.RequestOption) (res *AudioTranslationNewResponse, err error) {
	opts = append(r.Options[:], opts...)
	path := "openai/v1/audio/translations"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

type AudioTranslationNewResponse struct {
	Text string                          `json:"text,required"`
	JSON audioTranslationNewResponseJSON `json:"-"`
}

// audioTranslationNewResponseJSON contains the JSON metadata for the struct
// [AudioTranslationNewResponse]
type audioTranslationNewResponseJSON struct {
	Text        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AudioTranslationNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r audioTranslationNewResponseJSON) RawJSON() string {
	return r.raw
}

type AudioTranslationNewParams struct {
	// The audio file object (not file name) translate, in one of these formats: flac,
	// mp3, mp4, mpeg, mpga, m4a, ogg, wav, or webm.
	File param.Field[io.Reader] `json:"file,required" format:"binary"`
	// ID of the model to use. Only `whisper-large-v3` is currently available.
	Model param.Field[AudioTranslationNewParamsModel] `json:"model,required"`
	// An optional text to guide the model's style or continue a previous audio
	// segment. The [prompt](/docs/guides/speech-to-text/prompting) should be in
	// English.
	Prompt param.Field[string] `json:"prompt"`
	// The format of the transcript output, in one of these options: `json`, `text`, or
	// `verbose_json`.
	ResponseFormat param.Field[AudioTranslationNewParamsResponseFormat] `json:"response_format"`
	// The sampling temperature, between 0 and 1. Higher values like 0.8 will make the
	// output more random, while lower values like 0.2 will make it more focused and
	// deterministic. If set to 0, the model will use
	// [log probability](https://en.wikipedia.org/wiki/Log_probability) to
	// automatically increase the temperature until certain thresholds are hit.
	Temperature param.Field[float64] `json:"temperature"`
}

func (r AudioTranslationNewParams) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err != nil {
		writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

type AudioTranslationNewParamsModel string

const (
	AudioTranslationNewParamsModelWhisperLargeV3 AudioTranslationNewParamsModel = "whisper-large-v3"
)

func (r AudioTranslationNewParamsModel) IsKnown() bool {
	switch r {
	case AudioTranslationNewParamsModelWhisperLargeV3:
		return true
	}
	return false
}

// The format of the transcript output, in one of these options: `json`, `text`, or
// `verbose_json`.
type AudioTranslationNewParamsResponseFormat string

const (
	AudioTranslationNewParamsResponseFormatJson        AudioTranslationNewParamsResponseFormat = "json"
	AudioTranslationNewParamsResponseFormatText        AudioTranslationNewParamsResponseFormat = "text"
	AudioTranslationNewParamsResponseFormatVerboseJson AudioTranslationNewParamsResponseFormat = "verbose_json"
)

func (r AudioTranslationNewParamsResponseFormat) IsKnown() bool {
	switch r {
	case AudioTranslationNewParamsResponseFormatJson, AudioTranslationNewParamsResponseFormatText, AudioTranslationNewParamsResponseFormatVerboseJson:
		return true
	}
	return false
}
