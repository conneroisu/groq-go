package jigsawstack

import (
	"context"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	ttsEndpoint     = "/v1/ai/tts"
	accentsEndpoint = "/v1/audio/speaker_voice_accents"
)

type (
	// TTSOption is an option for the TTS request.
	TTSOption func(*TTSRequest)
	// TTSRequest represents a request structure for TTS API.
	TTSRequest struct {
		Text string `json:"text"`
		// Accent is the accent of the speaker voice to use.
		//
		// It is only used if the FileKey or SpeakerURL is not provided.
		Accent string `json:"accent,omitempty"`
		// SpeakerURL is the url of the speaker voice to use.
		//
		// It is only used if the FileKey is not provided.
		SpeakerURL string `json:"speaker_clone_url,omitempty"`
		// FileKey is the key of the file to use as the speaker voice.
		//
		// It is only used if the SpeakerURL is not provided.
		FileKey string `json:"speaker_clone_file_store_key,omitempty"`
	}
	// SpeakerVoiceAccent represents a speaker voice accent.
	SpeakerVoiceAccent struct {
		Success bool `json:"success"`
		Accents []struct {
			Accent     string `json:"accent"`
			LocaleName string `json:"locale_name"`
			Gender     string `json:"gender"`
		} `json:"accents"`
	}
)

// AudioTTS creates a text to speech (TTS) audio file.
//
// POST https://api.jigsawstack.com/v1/ai/tts
//
// https://docs.jigsawstack.com/api-reference/ai/text-to-speech
func (j *JigsawStack) AudioTTS(
	ctx context.Context,
	request TTSRequest,
) (mp3 string, err error) {
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodPost,
		j.baseURL+ttsEndpoint,
		builders.WithBody(request),
	)
	if err != nil {
		return
	}
	var resp string
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}

// AudioGetSpeakerVoiceAccents gets the speaker voice accents.
//
// GET https://api.jigsawstack.com/v1/ai/tts
//
// https://docs.jigsawstack.com/api-reference/audio/speaker-voice-accents
func (j *JigsawStack) AudioGetSpeakerVoiceAccents(
	ctx context.Context,
) (response []SpeakerVoiceAccent, err error) {
	uri := j.baseURL + accentsEndpoint
	req, err := builders.NewRequest(
		ctx,
		j.header,
		http.MethodGet,
		uri,
	)
	if err != nil {
		return
	}
	var resp []SpeakerVoiceAccent
	err = j.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}
