// Package toolhouse provides a Toolhouse extension for groq-go.
package toolhouse

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/conneroisu/groq-go/pkg/builders"
)

const (
	defaultBaseURL   = "https://api.toolhouse.ai/v1"
	getToolsEndpoint = "/get_tools"
	runToolEndpoint  = "/run_tools"
	applicationJSON  = "application/json"
)

type (
	// Toolhouse is a Toolhouse extension.
	Toolhouse struct {
		apiKey   string
		baseURL  string
		provider string
		bundle   string
		client   *http.Client
		metadata map[string]any
		logger   *slog.Logger
		header   builders.Header
	}

	// Options is a function that sets options for a Toolhouse extension.
	Options func(*Toolhouse)
)

// NewExtension creates a new Toolhouse extension.
func NewExtension(apiKey string, opts ...Options) (e *Toolhouse, err error) {
	e = &Toolhouse{
		apiKey:   apiKey,
		baseURL:  defaultBaseURL,
		client:   http.DefaultClient,
		bundle:   "default",
		provider: "openai",
		logger:   slog.Default(),
	}
	e.header.SetCommonHeaders = func(req *http.Request) {
		req.Header.Set("User-Agent", "Toolhouse/1.2.1 Python/3.11.0")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.apiKey))
		req.Header.Set("Content-Type", applicationJSON)
	}
	for _, opt := range opts {
		opt(e)
	}
	if e.apiKey == "" {
		err = fmt.Errorf("api key is required")
		return
	}
	return e, nil
}

func (e *Toolhouse) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json")
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK ||
		res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("failed to send http request: %s", res.Status)
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
		e.logger.Debug("decoding json response")
		err = json.NewDecoder(res.Body).Decode(v)
		if err != nil {
			read, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			e.logger.Debug("failed to decode response", "response", string(read))
			return fmt.Errorf("failed to decode response: %s", string(read))
		}
		return nil
	}
}
