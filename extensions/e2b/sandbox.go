package e2b

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	gen "github.com/conneroisu/groq-go/extensions/e2b/gen"
)

func init() {
	gen.Init()
}

type (
	// SandboxTemplate is a sandbox template.
	SandboxTemplate string
	// Sandbox is a code sandbox.
	//
	// The sandbox is like an isolated runtime or playground for the LLM.
	Sandbox struct {
		apiKey  string
		baseURL string
		client  *http.Client

		metaData map[string]string
		template SandboxTemplate
		// cwd      string
		// envVars  map[string]string
		logger *slog.Logger
	}
	// Process is a process in the sandbox.
	Process struct {
		ext *Sandbox
	}
	// Kernel is a code kernel.
	//
	// It is effectively a separate runtime environment inside of a sandbox.
	//
	// You can imagine kernel as a separate environment where code is
	// executed.
	//
	// You can have multiple kernels running at the same time.
	//
	// Each kernel has its own state, so you can have multiple kernels
	// running different code at the same time.
	//
	// A kernel will be kept alive with the sandbox even if you disconnect.
	// So, it may be useful to defer the shutdown of the kernel.
	Kernel struct {
		ext *Sandbox
	}
	// Event is a file system event.
	Event struct {
		// Path is the path of the event.
		Path string
		// Name is the name of file or directory.
		Name string
		// Operation is the operation type of the event.
		Operation OperationType
		// Timestamp is the timestamp of the event.
		Timestamp int64
		// IsDir is true if the event is a directory.
		IsDir bool
	}
	// OperationType is an operation type.
	OperationType int
	request       struct {
		apiKey  string
		cwd     string
		envVars map[string]string
		event   Event
		timeout time.Duration
	}
	// Option is an option for the sandbox.
	Option func(*Sandbox)
)

const (
	// EventTypeCreate is the type of event for the creation of a file or
	// directory.
	EventTypeCreate OperationType = iota
	// EventTypeWrite is the type of event for the write to a file.
	EventTypeWrite
	// EventTypeRemove is the type of event for the removal of a file or
	// directory.
	EventTypeRemove

	defaultBaseURL = "https://api.e2b.dev/v2"

	wsRoute   = "/ws"
	fileRoute = "/file"
	// (GET/POST /sandboxes)
	sandboxesRoute = "/sandboxes"
	// (DELETE /sandboxes/:id)
	deleteSandboxRoute = "/sandboxes/%s"
)

// NewSandbox creates a new sandbox.
func NewSandbox(
	apiKey string,
	template SandboxTemplate,
	opts ...Option,
) (Sandbox, error) {
	sb := Sandbox{
		apiKey:   apiKey,
		template: "base",
		baseURL:  defaultBaseURL,
		client:   http.DefaultClient,
		logger:   slog.Default(),
	}
	for _, opt := range opts {
		opt(&sb)
	}
	return sb, nil
}

// WithBaseURL sets the base URL for the e2b sandbox.
func (s *Sandbox) WithBaseURL(baseURL string) Option {
	return func(s *Sandbox) { s.baseURL = baseURL }
}

// WithClient sets the client for the e2b sandbox.
func WithClient(client *http.Client) Option {
	return func(s *Sandbox) { s.client = client }
}

// WithTemplate sets the template for the e2b sandbox.
func (s *Sandbox) WithTemplate(template SandboxTemplate) Option {
	return func(s *Sandbox) { s.template = template }
}

// WithLogger sets the logger for the e2b sandbox.
func WithLogger(logger *slog.Logger) Option {
	return func(s *Sandbox) { s.logger = logger }
}

// WithMetaData sets the meta data for the e2b sandbox.
func WithMetaData(metaData map[string]string) Option {
	return func(s *Sandbox) { s.metaData = metaData }
}

// Read reads a file from the sandbox file system.
func (s *Sandbox) Read(
	path string,
) ([]byte, error) {
	return nil, nil
}

// Write writes to a file to the sandbox file system.
func (s *Sandbox) Write(path string, data []byte) error {
	return nil
}

// ReadBytes reads a file from the sandbox file system.
func (s *Sandbox) ReadBytes(path string) ([]byte, error) {
	return nil, nil
}

// Watch watches a directory in the sandbox file system.
func (s *Sandbox) Watch(path string) (<-chan Event, error) {
	return nil, nil
}

// Upload uploads a file to the sandbox file system.
func (s *Sandbox) Upload(r io.Reader, path string) error {
	return nil
}

// Download downloads a file from the sandbox file system.
func (s *Sandbox) Download(path string) (io.ReadCloser, error) {
	return nil, nil
}

// KeepAlive keeps the sandbox alive.
func (s *Sandbox) KeepAlive(timeout time.Duration) error {
	return nil
}

// Reconnect reconnects to the sandbox.
func (s *Sandbox) Reconnect(id string) error {
	return nil
}

// StartProcess starts a process in the sandbox.
//
// If the context is cancelled, the process will be killed.
func (s *Sandbox) StartProcess(
	ctx context.Context,
	cmd string,
) (*Process, error) {
	if ctx.Done() == nil {
		return nil, ctx.Err()
	}
	return nil, nil
}

// Mkdir makes a directory in the sandbox file system.
func (s *Sandbox) Mkdir(path string) error {
	return nil
}

// Ls lists the files and/or directories in the sandbox file system at
// the given path.
func (s *Sandbox) Ls(path string) ([]string, error) {
	return nil, nil
}

// ListKernels lists the kernels avaliable to the extension.
func (s *Sandbox) ListKernels() ([]Kernel, error) {
	return nil, nil
}

// CreateKernel creates a new kernel.
func (s *Sandbox) CreateKernel() (Kernel, error) {
	return Kernel{}, nil
}

// Shutdown shutdowns a kernel.
func (k *Kernel) Shutdown(kernel Kernel) error {
	return nil
}

// RestartKernel restarts a kernel.
func (k *Kernel) RestartKernel(kernel Kernel) error {
	return nil
}
