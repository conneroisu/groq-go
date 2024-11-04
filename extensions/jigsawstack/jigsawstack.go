package jigsawstack

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	defaultBaseURL = "https://api.jigsawstack.com"
)

type (
	// JigsawStack is a JigsawStack extension.
	JigsawStack struct {
		baseURL string
		client  *http.Client
		logger  *slog.Logger
		header  builders.Header
	}
	// Option is an option for the JigsawStack client.
	Option func(*JigsawStack)
	// Endpoint is the endpoint for the JigsawStack api.
	Endpoint string
)

// NewJigsawStack creates a new JigsawStack extension.
func NewJigsawStack(apiKey string, opts ...Option) (*JigsawStack, error) {
	j := &JigsawStack{
		baseURL: defaultBaseURL,
		client:  http.DefaultClient,
		logger:  slog.Default(),
	}
	for _, opt := range opts {
		opt(j)
	}
	j.header.SetCommonHeaders = func(req *http.Request) {
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}
	return j, nil
}

// WithBaseURL sets the base URL for the JigsawStack extension.
func WithBaseURL(baseURL string) Option {
	return func(j *JigsawStack) { j.baseURL = baseURL }
}

// WithClient sets the client for the JigsawStack extension.
func WithClient(client *http.Client) Option {
	return func(j *JigsawStack) { j.client = client }
}

// WithLogger sets the logger for the JigsawStack extension.
func WithLogger(logger *slog.Logger) Option {
	return func(j *JigsawStack) { j.logger = logger }
}

func (j *JigsawStack) sendRequest(req *http.Request, v any) error {
	j.header.SetCommonHeaders(req)
	resp, err := j.client.Do(req)
	if err != nil {
		return err
	}
	j.logger.Debug("received http response from jigsawstack", "status", resp.Status, "headers", resp.Header, "url", req.URL)
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusBadRequest {
		read, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("bad status code: %d\nbdy: %s\n headers: %v", resp.StatusCode, read, resp.Header)
	}
	if v == nil {
		return nil
	}
	switch o := v.(type) {
	case *string:
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		*o = string(b)
		return nil
	default:
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			read, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			j.logger.Debug("failed to decode response", "response", string(read))
			return fmt.Errorf("failed to decode response: %w\nbody: %s", err, string(read))
		}
		return nil
	}
}
