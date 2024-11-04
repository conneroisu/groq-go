package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	ttsEndpoint     Endpoint = "/v1/ai/tts"
	accentsEndpoint Endpoint = "/v1/audio/speaker_voice_accents"
)

type (
	// TTSOption is an option for the TTS request.
	TTSOption func(*ttsRequest)
	// ttsRequest represents a request structure for TTS API.
	ttsRequest struct {
		// Text is the text to convert to speech.
		// Required.
		Text string `json:"text"`
		// Accent is the accent of the speaker voice to use.
		//
		// Not required if the FileKey or SpeakerURL is not provided.
		Accent string `json:"accent,omitempty"`
		// SpeakerURL is the url of the speaker voice to use.
		//
		// Not required if the FileKey is not provided.
		SpeakerURL string `json:"speaker_clone_url,omitempty"`
		// FileKey is the key of the file to use as the speaker voice.
		//
		// Not required if the SpeakerURL is not provided.
		FileKey string `json:"speaker_clone_file_store_key,omitempty"`
	}
)

// WithAccent sets the accent of the speaker voice to use.
func WithAccent(accent string) TTSOption {
	return func(r *ttsRequest) { r.Accent = accent }
}

// WithSpeakerURL sets the url of the speaker voice to use.
func WithSpeakerURL(url string) TTSOption {
	return func(r *ttsRequest) { r.SpeakerURL = url }
}

// WithFileKey sets the file key of the speaker voice to use.
func WithFileKey(key string) TTSOption {
	return func(r *ttsRequest) { r.FileKey = key }
}

// AudioTTS creates a text to speech (TTS) audio file.
//
// It only support one option at a time, but does support no options.
//
// POST https://api.jigsawstack.com/v1/ai/tts
//
// https://docs.jigsawstack.com/api-reference/ai/text-to-speech
func (j *JigsawStack) AudioTTS(
	ctx context.Context,
	text string,
	options ...TTSOption,
) (mp3 string, err error) {
	body := ttsRequest{
		Text: text,
	}
	for _, option := range options {
		option(&body)
	}
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+string(ttsEndpoint),
		builders.WithBody(body),
	)
	if err != nil {
		return "", err
	}
	var resp string
	err = j.sendRequest(req, &resp)
	if err != nil {
		return "", err
	}
	return resp, nil
}
