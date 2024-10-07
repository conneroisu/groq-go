package e2b

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type (
	// SandboxTemplate is a sandbox template.
	SandboxTemplate string
	// Sandbox is a code sandbox.
	//
	// The sandbox is like an isolated runtime or playground for the LLM.
	Sandbox struct {
		ID       string            `json:"sandboxID"`
		Metadata map[string]string `json:"metadata"`
		Template SandboxTemplate   `json:"templateID"`
		Alias    string            `json:"alias"`
		ClientID string            `json:"clientID"`
		apiKey   string
		baseURL  string
		wsURL    string
		client   *http.Client
		ws       *websocket.Conn
		msgCnt   int
		mu       *sync.Mutex
		// cwd      string
		// envVars  map[string]string
		logger *slog.Logger
	}
	// CreateSandboxResponse represents the response of the create sandbox http method.
	CreateSandboxResponse struct {
		Alias       string `json:"alias"`
		ClientID    string `json:"clientID"`
		EnvdVersion string `json:"envdVersion"`
		SandboxID   string `json:"sandboxID"`
		TemplateID  string `json:"templateID"`
	}
	// Process is a process in the sandbox.
	Process struct {
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
	defaultBaseURL  = "https://api.e2b.dev"
	defaultWSScheme = "wss"
	wsRoute         = "/ws"
	fileRoute       = "/file"
	// (GET/POST /sandboxes)
	sandboxesRoute = "/sandboxes"
	// (DELETE /sandboxes/:id)
	deleteSandboxRoute = "/sandboxes/%s"
	// Kernels Endpoint
	kernelsRoute = "/api/kernels"
)

// NewSandbox creates a new sandbox.
func NewSandbox(
	apiKey string,
	template SandboxTemplate,
	opts ...Option,
) (Sandbox, error) {
	sb := Sandbox{
		mu:       &sync.Mutex{},
		apiKey:   apiKey,
		Template: "base",
		baseURL:  defaultBaseURL,
		Metadata: map[string]string{
			"name": "groq-go",
		},
		client: http.DefaultClient,
		logger: slog.Default(),
	}
	jsVal, err := json.Marshal(sb)
	if err != nil {
		return Sandbox{}, err
	}
	for _, opt := range opts {
		opt(&sb)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s%s", sb.baseURL, sandboxesRoute),
		bytes.NewBuffer([]byte(jsVal)),
	)
	if err != nil {
		return Sandbox{}, err
	}
	req.Header.Set("X-API-Key", sb.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := sb.client.Do(req)
	if err != nil {
		return Sandbox{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return Sandbox{}, fmt.Errorf("request to create sandbox failed: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Sandbox{}, err
	}
	var res CreateSandboxResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return Sandbox{}, err
	}
	sb.ID = res.SandboxID
	sb.Alias = res.Alias
	sb.ClientID = res.ClientID
	sb.wsURL = fmt.Sprintf("49982-%s-%s.e2b.dev",
		res.SandboxID,
		res.ClientID,
	)
	u := url.URL{
		Scheme: defaultWSScheme,
		Host:   sb.wsURL,
		Path:   wsRoute,
	}
	sb.logger.Debug("Connecting to sandbox", "url", u.String())
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return Sandbox{}, err
	}
	sb.ws = ws
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
	return func(s *Sandbox) { s.Template = template }
}

// WithLogger sets the logger for the e2b sandbox.
func WithLogger(logger *slog.Logger) Option {
	return func(s *Sandbox) { s.logger = logger }
}

// WithMetaData sets the meta data for the e2b sandbox.
func WithMetaData(metaData map[string]string) Option {
	return func(s *Sandbox) { s.Metadata = metaData }
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
) (proc *Process, err error) {
	if ctx.Done() == nil {
		return nil, ctx.Err()
	}
	return nil, nil
}

// ListKernels lists the kernels in the sandbox.
func (s *Sandbox) ListKernels() ([]Kernel, error) {
	url := fmt.Sprintf("https://%s%s%s", "8888", s.wsURL[len("49982"):], kernelsRoute)
	println(url)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))
	// var res ListKernelsResponse
	return nil, nil
}

// CreateKernel creates a new kernel.
func (s *Sandbox) CreateKernel() (Kernel, error) {
	return Kernel{}, nil
}

// Close closes the sandbox.
func (s *Sandbox) Close() error {
	return s.ws.Close()
}
