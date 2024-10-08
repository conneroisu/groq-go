package e2b

import (
	"log/slog"
	"net/http"
)

// WithBaseURL sets the base URL for the e2b sandbox.
func (s *Sandbox) WithBaseURL(baseURL string) Option {
	return func(s *Sandbox) { s.baseAPIURL = baseURL }
}

// WithClient sets the client for the e2b sandbox.
func WithClient(client *http.Client) Option {
	return func(s *Sandbox) { s.client = client }
}

// WithTemplate sets the template for the e2b sandbox.
func (s *Sandbox) WithTemplate(template SandboxTemplate) Option {
	return func(s *Sandbox) { s.Template = template }
}

// WithLogger sets the logger for the e2b sandbox.
func WithLogger(logger *slog.Logger) Option {
	return func(s *Sandbox) { s.logger = logger }
}

// WithTemplate sets the template for the e2b sandbox.
func WithTemplate(template SandboxTemplate) Option {
	return func(s *Sandbox) { s.Template = template }
}

// WithMetaData sets the meta data for the e2b sandbox.
func WithMetaData(metaData map[string]string) Option {
	return func(s *Sandbox) { s.Metadata = metaData }
}
