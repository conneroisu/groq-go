package e2b

import (
	"log/slog"
	"net/http"
)

// E2B Sandbox Options

// WithBaseURL sets the base URL for the e2b sandbox.
func WithBaseURL(baseURL string) Option {
	return func(s *Sandbox) { s.baseURL = baseURL }
}

// WithClient sets the client for the e2b sandbox.
func WithClient(client *http.Client) Option {
	return func(s *Sandbox) { s.client = client }
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

// WithCwd sets the current working directory.
func WithCwd(cwd string) Option {
	return func(s *Sandbox) { s.Cwd = cwd }
}

// WithWsURL sets the websocket url resolving function for the e2b sandbox.
func WithWsURL(wsURL func(s *Sandbox) string) Option {
	return func(s *Sandbox) { s.wsURL = wsURL }
}

// Process Options
