package toolhouse

import (
	"log/slog"
	"net/http"
)

// WithBaseURL sets the base URL for the Toolhouse extension.
func WithBaseURL(baseURL string) Options {
	return func(e *Toolhouse) {
		e.baseURL = baseURL
	}
}

// WithClient sets the client for the Toolhouse extension.
func WithClient(client *http.Client) Options {
	return func(e *Toolhouse) {
		e.client = client
	}
}

// WithMetadata sets the metadata for the get tools request.
func WithMetadata(metadata map[string]any) Options {
	return func(r *Toolhouse) {
		r.metadata = metadata
	}
}

// WithLogger sets the logger for the Toolhouse extension.
func WithLogger(logger *slog.Logger) Options {
	return func(r *Toolhouse) {
		r.logger = logger
	}
}
