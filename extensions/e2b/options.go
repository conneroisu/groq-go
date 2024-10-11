package e2b

import (
	"log/slog"
	"net/http"
)

// WithBaseURL sets the base URL for the e2b sandbox.
func (s *Sandbox) WithBaseURL(baseURL string) Option {
	return func(s *Sandbox) { s.baseAPIURL = baseURL }
}

// WithTemplate sets the template for the e2b sandbox.
func (s *Sandbox) WithTemplate(template SandboxTemplate) Option {
	return func(s *Sandbox) { s.Template = template }
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

// WithHTTPScheme sets the http scheme for the e2b sandbox.
func WithHTTPScheme(scheme string) Option {
	return func(s *Sandbox) { s.httpScheme = scheme }
}

// WithEnv sets the environment variables.
func WithEnv(env map[string]string) ProcessOption {
	return func(p *Process) {
		p.env = env
	}
}

// WithCwd sets the current working directory.
func WithCwd(cwd string) ProcessOption {
	return func(p *Process) {
		p.cwd = cwd
	}
}
