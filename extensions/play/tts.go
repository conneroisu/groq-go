package play

import (
	"io"
)

const (
	ttsStreamEndpoint Endpoint = "/api/v1/tts/stream"
)

// TTSStreamRequest represents a request structure for TTS API.
type TTSStreamRequest struct {
	Model                     string `json:"model"`
	Text                      string `json:"text"`
	Voice                     string `json:"voice"`
	Voice2                    string `json:"voice2"`
	OutputFormat              string `json:"outputFormat"`
	Speed                     int    `json:"speed"`
	SampleRate                int    `json:"sampleRate"`
	Seed                      any    `json:"seed"`
	Temperature               any    `json:"temperature"`
	TurnPrefix                string `json:"turnPrefix"`
	TurnPrefix2               string `json:"turnPrefix2"`
	Prompt                    string `json:"prompt"`
	Prompt2                   string `json:"prompt2"`
	VoiceConditioningSeconds  int    `json:"voiceConditioningSeconds"`
	VoiceConditioningSeconds2 int    `json:"voiceConditioningSeconds2"`
}

func (p *PlayAI) TTSStream(req TTSStreamRequest) (io.Reader, error) {
	var res string
	return res, err
}
