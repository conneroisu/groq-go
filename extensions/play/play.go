package play

import (
	"io"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	defaultBaseURL = "https://api.play.ai/"
)

// PlayAI is a PlayAI extension.
type (
	PlayAI struct {
		apiKey  string
		baseURL string
		header  builders.Header
		client  *http.Client
	}
	// Endpoint is an endpoint for the PlayAI api.
	Endpoint string
)

// New creates a new PlayAI extension.
func New(apiKey string, userID string, opts ...Option) (*PlayAI, error) {
	p := PlayAI{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		header: builders.Header{
			SetCommonHeaders: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer "+apiKey)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-USER-ID", userID)
			},
		},
	}
	for _, opt := range opts {
		opt(&p)
	}
	return &p, nil
}

func (p *PlayAI) sendRequest(req *http.Request, v any) error {
	req.Header.Set("Accept", "application/json")
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := p.client.Do(req)
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
		return nil
	}
}
