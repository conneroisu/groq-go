// Package main is a script to generate the jigsaw accents for the project.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/conneroisu/groq-go/pkg/test"
)

const (
	defaultBaseURL  = "https://api.jigsawstack.com"
	accentsEndpoint = "/v1/ai/tts"
)

type (
	// Client is a JigsawStack extension.
	Client struct {
		baseURL string
		apiKey  string
		client  *http.Client
		logger  *slog.Logger
		header  builders.Header
	}
	// SpeakerVoiceAccent represents a speaker voice accent.
	SpeakerVoiceAccent struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Accents []struct {
			Accent     string `json:"accent"`
			LocaleName string `json:"locale_name"`
			Gender     string `json:"gender"`
		} `json:"accents"`
	}
)

// AudioGetSpeakerVoiceAccents gets the speaker voice accents.
//
// GET https://api.jigsawstack.com/v1/ai/tts
//
// https://docs.jigsawstack.com/api-reference/audio/speaker-voice-accents
func (c *Client) AudioGetSpeakerVoiceAccents(
	ctx context.Context,
) (response SpeakerVoiceAccent, err error) {
	uri := c.baseURL + accentsEndpoint
	req, err := builders.NewRequest(
		ctx,
		c.header,
		http.MethodGet,
		uri,
	)
	if err != nil {
		return
	}
	var resp SpeakerVoiceAccent
	err = c.sendRequest(req, &resp)
	if err != nil {
		return
	}
	return resp, nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func newClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		header: builders.Header{SetCommonHeaders: func(r *http.Request) {
			r.Header.Set("x-api-key", apiKey)
		}},
		client:  http.DefaultClient,
		baseURL: defaultBaseURL,
	}
}

func run(ctx context.Context) error {
	println("generating accents")
	key, err := test.GetAPIKey("JIGSAWSTACK_API_KEY")
	if err != nil {
		return err
	}
	client := newClient(key)
	accents, err := client.AudioGetSpeakerVoiceAccents(ctx)
	if err != nil {
		return err
	}
	println(len(accents.Accents))
	if !accents.Success {
		return fmt.Errorf("failed to get accents: %v", accents.Message)
	}
	for _, accent := range accents.Accents {
		println(accent.Accent)
	}
	return nil
}

func (c *Client) sendRequest(req *http.Request, v any) error {
	req.Header.Set("Accept", "application/json")
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK ||
		res.StatusCode >= http.StatusBadRequest {
		return nil
	}
	if v == nil {
		return nil
	}
	switch o := v.(type) {
	case *string:
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		*o = string(b)
		return nil
	default:
		err = json.NewDecoder(res.Body).Decode(v)
		if err != nil {
			read, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			c.logger.Debug("failed to decode response", "response", string(read))
			return fmt.Errorf("failed to decode response: %w\nbody: %s", err, string(read))
		}
		return nil
	}
}
