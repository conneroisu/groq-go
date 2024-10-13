package e2b

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/conneroisu/groq-go/pkg/builders"
)

type (
	// ProcessEvents is a process event type.
	// string
	ProcessEvents string
	// SandboxTemplate is a sandbox template.
	SandboxTemplate string
	// Sandbox is a code sandbox.
	//
	// The sandbox is like an isolated, but interactive system.
	Sandbox struct {
		ID         string            `json:"sandboxID"`
		Metadata   map[string]string `json:"metadata"`
		Template   SandboxTemplate   `json:"templateID"`
		Alias      string            `json:"alias"`
		ClientID   string            `json:"clientID"`
		apiKey     string
		baseAPIURL string
		logger     *slog.Logger
		header     builders.Header
		httpScheme string
		client     *http.Client
		wsH        *WSHandler
	}
	// Process is a process in the sandbox.
	Process struct {
		ext      *Sandbox
		ID       string
		ResultID string
		cmd      string
		cwd      string
		env      map[string]string
	}
	// SubscribeParams is the params for subscribing to a process event.
	SubscribeParams struct {
		Event ProcessEvents
		Ch    chan<- Event
	}
	// Option is an option for the sandbox.
	Option func(*Sandbox)
	// ProcessOption is an option for a process.
	ProcessOption func(*Process)
	// Event is a file system event.
	// {\"jsonrpc\":\"2.0\",\"method\":\"process_subscription\",\"params\":{\"subscription\":\"0xc900b4c1c65808e80174d22e2ce9ecf4\",\"result\":{\"type\":\"Stdout\",\"line\":\"Hello World!\",\"timestamp\":1728774047677344401}}}\n
	Event struct {
		Params struct {
			Subscription string `json:"subscription"`
			Result       struct {
				Type        string `json:"type"`
				Line        string `json:"line"`
				Timestamp   int64  `json:"timestamp"`
				IsDirectory bool   `json:"isDirectory"`
				Error       string `json:"error"`
			} `json:"result"`
		}
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
		// Error is the possible error of the event.
		Error string
	}
	// OperationType is an operation type.
	OperationType int
	// Method is a JSON-RPC method.
	Method string
	// Request is a JSON-RPC request.
	Request struct {
		// JSONRPC is the JSON-RPC version of the message.
		JSONRPC string `json:"jsonrpc"`
		// Method is the method of the message.
		Method Method `json:"method"`
		// ID is the ID of the message.
		ID int `json:"id"`
		// Params is the params of the message.
		Params []any `json:"params"`
		// ResponseCh is the channel to write the response to.
		ResponseCh chan []byte `json:"-"`
	}
	// Response is a JSON-RPC response.
	Response[T any, Q any] struct {
		// ID is the ID of the message.
		ID int `json:"id"`
		// Result is the result of the message.
		Result T `json:"result"`
		// Error is the error of the message.
		Error Q `json:"error"`
	}
	// LsResult is a result of the list request.
	LsResult struct {
		Name  string `json:"name"`
		IsDir bool   `json:"isDir"`
	}
	// APIError is the error of the API.
	APIError struct {
		Code    int    `json:"code,omitempty"` // Code is the code of the error.
		Message string `json:"message"`        // Message is the message of the error.
	}
)

const (
	// OnStdout is the event for the stdout.
	OnStdout ProcessEvents = "onStdout"
	// OnStderr is the event for the stderr.
	OnStderr ProcessEvents = "onStderr"
	// OnExit is the event for the exit.
	OnExit ProcessEvents = "onExit"

	rpc                         = "2.0"
	charset                     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	filesystemWrite      Method = "filesystem_write"
	filesystemRead       Method = "filesystem_read"
	filesystemList       Method = "filesystem_list"
	filesystemRemove     Method = "filesystem_remove"
	filesystemMakeDir    Method = "filesystem_makeDir"
	filesystemReadBytes  Method = "filesystem_readBase64"
	filesystemWriteBytes Method = "filesystem_writeBase64"
	processSubscribe     Method = "process_subscribe"
	processUnsubscribe   Method = "process_unsubscribe"
	processStart         Method = "process_start"
	// TODO: Check this one.
	filesystemSubscribe = "filesystem_subscribe"
	defaultBaseURL      = "https://api.e2b.dev"
	defaultWSScheme     = "wss"
	wsRoute             = "/ws"
	fileRoute           = "/file"
	sandboxesRoute      = "/sandboxes"  // (GET/POST /sandboxes)
	deleteSandboxRoute  = "/sandboxes/" // (DELETE /sandboxes/:id)
	defaultHTTPScheme   = "https"
	// EventTypeCreate is a type of event for the creation of a file/dir.
	EventTypeCreate OperationType = iota
	// EventTypeWrite is a type of event for the write to a file.
	EventTypeWrite
	// EventTypeRemove is a type of event for the removal of a file/dir.
	EventTypeRemove
)

// NewSandbox creates a new sandbox.
func NewSandbox(
	ctx context.Context,
	apiKey string,
	opts ...Option,
) (Sandbox, error) {
	sb := Sandbox{
		apiKey:     apiKey,
		Template:   "base",
		baseAPIURL: defaultBaseURL,
		Metadata: map[string]string{
			"sdk": "groq-go v1",
		},
		client:     http.DefaultClient,
		logger:     slog.Default(),
		httpScheme: defaultHTTPScheme,
	}
	for _, opt := range opts {
		opt(&sb)
	}
	sb.header.SetCommonHeaders = func(req *http.Request) {
		req.Header.Set("X-API-Key", sb.apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}
	req, err := builders.NewRequest(
		ctx,
		sb.header,
		http.MethodPost,
		fmt.Sprintf("%s%s", sb.baseAPIURL, sandboxesRoute),
		builders.WithBody(sb),
	)
	if err != nil {
		return sb, err
	}
	err = sb.sendRequest(req, &sb)
	if err != nil {
		return sb, err
	}
	u := sb.wsURL()
	sb.logger.Debug("Connecting to sandbox", "url", u.String())
	sb.wsH, err = newWSHandler(ctx, u.String())
	if err != nil {
		return sb, err
	}
	return sb, nil
}

// KeepAlive keeps the sandbox alive.
func (s *Sandbox) KeepAlive(timeout time.Duration) error {
	time.Sleep(timeout)
	// TODO: implement
	return nil
}

// Reconnect reconnects to the sandbox.
func (s *Sandbox) Reconnect(ctx context.Context) (err error) {
	u := s.wsURL()
	s.logger.Debug("Reconnecting to sandbox", "url", u.String())
	s.wsH, err = newWSHandler(ctx, u.String())
	return err
}

// Disconnect disconnects from the sandbox.
func (s *Sandbox) Disconnect() error {
	return s.wsH.ws.Close()
}

// Stop stops the sandbox.
func (s *Sandbox) Stop(ctx context.Context) error {
	req, err := builders.NewRequest(
		ctx,
		s.header,
		http.MethodDelete,
		fmt.Sprintf("%s%s%s", s.baseAPIURL, deleteSandboxRoute, s.ID),
		builders.WithBody(interface{}(nil)),
	)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("request to delete sandbox failed: %s", resp.Status)
	}
	return nil
}

// Mkdir makes a directory in the sandbox file system.
func (s *Sandbox) Mkdir(path string) error {
	respCh := make(chan []byte)
	err := s.wsH.Write(Request{
		JSONRPC:    rpc,
		Method:     filesystemMakeDir,
		Params:     []any{path},
		ResponseCh: respCh,
	})
	if err != nil {
		return err
	}
	var resp Response[string, string]
	err = json.Unmarshal(<-respCh, &resp)
	if err != nil {
		return fmt.Errorf("failed to mkdir: %w", err)
	}
	if resp.Error != "" {
		return fmt.Errorf("failed to write to file: %s", resp.Error)
	}
	return nil
}

// Ls lists the files and/or directories in the sandbox file system at
// the given path.
func (s *Sandbox) Ls(path string) ([]LsResult, error) {
	respCh := make(chan []byte)
	err := s.wsH.Write(Request{
		Params:     []any{path},
		JSONRPC:    rpc,
		Method:     filesystemList,
		ResponseCh: respCh,
	})
	if err != nil {
		return nil, err
	}
	var res Response[[]LsResult, string]
	err = json.Unmarshal(<-respCh, &res)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

// Read reads a file from the sandbox file system.
func (s *Sandbox) Read(
	path string,
) (string, error) {
	respCh := make(chan []byte)
	err := s.wsH.Write(Request{
		JSONRPC:    rpc,
		Method:     filesystemRead,
		Params:     []any{path},
		ResponseCh: respCh,
	})
	if err != nil {
		return "", err
	}
	var res Response[string, string]
	err = json.Unmarshal(<-respCh, &res)
	if err != nil {
		return "", err
	}
	if res.Error != "" {
		return "", fmt.Errorf("failed to read file: %s", res.Error)
	}
	return res.Result, nil
}

// Write writes to a file to the sandbox file system.
func (s *Sandbox) Write(path string, data []byte) error {
	respCh := make(chan []byte)
	err := s.wsH.Write(Request{
		JSONRPC:    rpc,
		Method:     filesystemWrite,
		Params:     []any{path, string(data)},
		ResponseCh: respCh,
	})
	if err != nil {
		return err
	}
	err = json.Unmarshal(<-respCh, &Request{})
	if err != nil {
		return err
	}
	return nil
}

// ReadBytes reads a file from the sandbox file system.
func (s *Sandbox) ReadBytes(path string) ([]byte, error) {
	resCh := make(chan []byte)
	err := s.wsH.Write(Request{
		JSONRPC:    rpc,
		Method:     filesystemReadBytes,
		Params:     []any{path},
		ResponseCh: resCh,
	})
	if err != nil {
		return nil, err
	}
	var rR Response[string, string]
	err = json.Unmarshal(<-resCh, &rR)
	if err != nil {
		return nil, err
	}
	sDec, err := base64.StdEncoding.DecodeString(rR.Result)
	if err != nil {
		return nil, err
	}
	return sDec, nil
}

// Watch watches a directory in the sandbox file system.
//
// This is intended to be run in a goroutine as it will block until the
// connection is closed, an error occurs, or the context is canceled.
func (s *Sandbox) Watch(
	ctx context.Context,
	path string,
) (<-chan Event, error) {
	// TODO: implement
	return nil, nil
}

// Upload uploads a file to the sandbox file system.
func (s *Sandbox) Upload(r io.Reader, path string) error {
	// TODO: implement
	return nil
}

// Download downloads a file from the sandbox file system.
func (s *Sandbox) Download(path string) (io.ReadCloser, error) {
	// TODO: implement
	return nil, nil
}

// NewProcess creates a new process startable in the sandbox.
func (s *Sandbox) NewProcess(
	cmd string,
	opts ...ProcessOption,
) (Process, error) {
	proc := Process{
		ID:  createProcessID(),
		ext: s,
		cmd: cmd,
	}
	for _, opt := range opts {
		opt(&proc)
	}
	return proc, nil
}

// Start starts a process in the sandbox.
func (p *Process) Start() error {
	if p.env == nil {
		p.env = map[string]string{"PYTHONUNBUFFERED": "1"}
	}
	respCh := make(chan []byte)
	err := p.ext.wsH.Write(Request{
		JSONRPC:    rpc,
		Method:     processStart,
		Params:     []any{p.ID, p.cmd, p.env, p.cwd},
		ResponseCh: respCh,
	})
	if err != nil {
		return err
	}
	var res Response[string, APIError]
	err = json.Unmarshal(<-respCh, &res)
	if err != nil {
		return err
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("process start failed(%d): %s", res.Error.Code, res.Error.Message)
	}
	if res.Result == "" || len(res.Result) == 0 {
		return fmt.Errorf("process start failed got empty result id")
	}
	if p.ID != res.Result {
		return fmt.Errorf("process start failed got wrong result id; want %s, got %s", p.ID, res.Result)
	}
	return nil
}

// Close closes the sandbox.
func (s *Sandbox) Close() error {
	return s.wsH.ws.Close()
}

// Done returns a channel that is closed when the process is done.
func (p *Process) Done() <-chan struct{} {
	rCh, ok := p.ext.wsH.idMap.Load(p.ID)
	if !ok {
		return nil
	}
	return rCh.(<-chan struct{})

}

// Subscribe subscribes to a process event.
func (p *Process) Subscribe(
	ctx context.Context,
	event ProcessEvents,
	ch chan<- Event,
) error {
	responseCh := make(chan []byte)
	err := p.ext.wsH.Write(Request{
		JSONRPC:    rpc,
		Method:     processSubscribe,
		Params:     []any{event, p.ID},
		ResponseCh: responseCh,
	})
	if err != nil {
		return err
	}
	var res Response[string, APIError]
	resp := <-responseCh
	err = json.Unmarshal(resp, &res)
	if err != nil {
		return err
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("process subscribe failed(%d): %s", res.Error.Code, res.Error.Message)
	}
	p.ext.logger.Debug("subscribed to process", "id", p.ID, "event", event, "subID", res.Result)
	go func() {
		select {
		case <-ctx.Done():
			break
		case <-p.Done():
			return
		}
		p.ext.logger.Debug("unsubscribing from process", "id", p.ID, "event", event, "subID", res.Result)
		respCh := make(chan []byte)
		err = p.ext.wsH.Write(Request{
			JSONRPC:    rpc,
			Method:     processUnsubscribe,
			Params:     []any{res.Result},
			ResponseCh: respCh,
		})
		if err != nil {
			println(err)
		}
		var unsubRes Response[bool, string]
		err = json.Unmarshal(<-respCh, &unsubRes)
		if err != nil {
			println(err)
		}
	}()
	eventByCh := make(chan []byte)
	p.ext.wsH.Sub(res.Result, eventByCh)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var event Event
				err = json.Unmarshal(<-eventByCh, &event)
				if err != nil {
					return
				}
				if event.Error != "" {
					return
				}
				if event.Params.Subscription != res.Result {
					continue
				}
				ch <- event
			}
		}
	}()
	return nil
}

func (s *Sandbox) wsURL() url.URL {
	return url.URL{
		Scheme: defaultWSScheme,
		Host: fmt.Sprintf("49982-%s-%s.e2b.dev",
			s.ID,
			s.ClientID,
		),
		Path: wsRoute,
	}
}

func createProcessID() string {
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
