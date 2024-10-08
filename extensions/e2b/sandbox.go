package e2b

import (
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
		ID         string            `json:"sandboxID"`
		Metadata   map[string]string `json:"metadata"`
		Template   SandboxTemplate   `json:"templateID"`
		Alias      string            `json:"alias"`
		ClientID   string            `json:"clientID"`
		apiKey     string
		baseAPIURL string
		client     *http.Client
		ws         *websocket.Conn
		msgCnt     int
		mu         *sync.Mutex
		// cwd      string
		// envVars  map[string]string
		logger         *slog.Logger
		requestBuilder requestBuilder
		httpScheme     string
	}
	// CreateSandboxResponse represents the response of the create sandbox http method.
	CreateSandboxResponse struct {
		Alias       string `json:"alias"`
		ClientID    string `json:"clientID"`
		EnvdVersion string `json:"envdVersion"`
		SandboxID   string `json:"sandboxID"`
		TemplateID  string `json:"templateID"`
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
	// Option is an option for the sandbox.
	Option func(*Sandbox)

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
	}
	// LsResponse is a JSON-RPC response when listing files and directories.
	LsResponse struct {
		// JSONRPC is the JSON-RPC version of the message.
		JSONRPC string `json:"jsonrpc"`
		// Method is the method of the message.
		Method Method `json:"method"`
		// ID is the ID of the message.
		ID     int        `json:"id"`
		Result []LsResult `json:"result"`
	}
	// ReadResponse is a JSON-RPC response when reading a file.
	ReadResponse struct {
		// JSONRPC is the JSON-RPC version of the message.
		JSONRPC string `json:"jsonrpc"`
		// Method is the method of the message.
		Method Method `json:"method"`
		// ID is the ID of the message.
		ID     int    `json:"id"`
		Result string `json:"result"`
	}
	// LsResult is a result of the list request.
	LsResult struct {
		Name  string `json:"name"`
		IsDir bool   `json:"isDir"`
	}
)

const (
	rpc = "2.0"

	filesystemWrite      Method = "filesystem_write"
	filesystemRead       Method = "filesystem_read"
	filesystemList       Method = "filesystem_list"
	filesystemRemove     Method = "filesystem_remove"
	filesystemMakeDir    Method = "filesystem_makeDir"
	filesystemReadBytes  Method = "filesystem_readBase64"
	filesystemWriteBytes Method = "filesystem_writeBase64"
	// TODO: Check this one.
	filesystemSubscribe = "filesystem_subscribe"
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

	defaultHTTPScheme = "https"
)

// NewSandbox creates a new sandbox.
func NewSandbox(
	apiKey string,
	opts ...Option,
) (Sandbox, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sb := Sandbox{
		mu:         &sync.Mutex{},
		apiKey:     apiKey,
		Template:   "base",
		baseAPIURL: defaultBaseURL,
		Metadata: map[string]string{
			"name": "groq-go",
		},
		client:         http.DefaultClient,
		logger:         slog.Default(),
		requestBuilder: newRequestBuilder(),
	}
	for _, opt := range opts {
		opt(&sb)
	}
	req, err := sb.newRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s%s", sb.baseAPIURL, sandboxesRoute),
		withBody(sb),
	)
	if err != nil {
		return sb, err
	}
	resp, err := sb.client.Do(req)
	if err != nil {
		return sb, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return sb, err
	}
	if resp.StatusCode != http.StatusCreated {
		return sb, fmt.Errorf("request to create sandbox failed: %s\nbody: %s", resp.Status, string(body))
	}
	var res CreateSandboxResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return sb, err
	}
	sb.ID = res.SandboxID
	sb.Alias = res.Alias
	sb.ClientID = res.ClientID
	u := sb.wsURL()
	sb.logger.Debug("Connecting to sandbox", "url", u.String())
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return sb, err
	}
	sb.ws = ws
	return sb, nil
}

// KeepAlive keeps the sandbox alive.
func (s *Sandbox) KeepAlive(timeout time.Duration) error {
	time.Sleep(timeout)
	// TODO: implement
	return nil
}

// Reconnect reconnects to the sandbox.
func (s *Sandbox) Reconnect( /* id string */ ) error {
	u := s.wsURL()
	s.logger.Debug("Reconnecting to sandbox", "url", u.String())
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	s.ws = ws
	return nil
}

// NewProcess starts a process in the sandbox.
//
// If the context is cancelled, the process will be killed.
func (s *Sandbox) NewProcess(
	cmd string,
) (proc Process, err error) {
	return Process{
		cmd: cmd,
	}, nil
}

// Stop stops the sandbox.
func (s *Sandbox) Stop(ctx context.Context) error {
	req, err := s.newRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s%s", s.baseAPIURL, fmt.Sprintf(deleteSandboxRoute, s.ID)),
		withBody(interface{}(nil)),
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

func (s *Sandbox) hostname() string {
	return fmt.Sprintf("49982-%s-%s.e2b.dev",
		s.ID,
		s.ClientID,
	)
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

func (s *Sandbox) httpURL(path string) url.URL {
	return url.URL{
		Scheme: defaultHTTPScheme,
		Host: fmt.Sprintf("49982-%s-%s.e2b.dev",
			s.ID,
			s.ClientID,
		),
		Path: path,
	}
}
