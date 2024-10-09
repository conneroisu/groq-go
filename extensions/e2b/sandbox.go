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
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type (
	processSubscribeParams struct {
		event ProcessEvents
		id    string
	}

	// ProcessEvents is a process event type.
	// string
	ProcessEvents string

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
		logger          *slog.Logger
		requestBuilder  requestBuilder
		httpScheme      string
		defaultKernelID string
	}
	// Process is a process in the sandbox.
	Process struct {
		ext *Sandbox

		ID       string
		ResultID string
		cmd      string
		cwd      string
		env      map[string]string
	}
	// CreateSandboxResponse represents the response of the create sandbox
	// http method.
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
	// ProcessRequestParams represents the params of the process request.
	ProcessRequestParams struct {
		// ID is the ID of the process.
		ID string
		// Command is the command to run.
		Command string `json:"command"`
		// Env is the environment variables.
		Env map[string]string `json:"env"`
		// Cwd is the current working directory.
		//
		// Blank means the current directory.
		Cwd string `json:"cwd"`
	}
)

const (
	rpc = "2.0"

	onStdout ProcessEvents = "onStdout"
	onStderr ProcessEvents = "onStderr"
	onExit   ProcessEvents = "onExit"

	charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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
	// Routes
	wsRoute   = "/ws"
	fileRoute = "/file"
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
			"sdk": "groq-go v1",
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
		return sb, fmt.Errorf(
			"request to create sandbox failed: %s\nbody: %s",
			resp.Status,
			string(body))
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
	kernelID, err := sb.Read("/root/.jupyter/kernel_id")
	if err != nil {
		return sb, err
	}
	sb.defaultKernelID = string(kernelID)
	return sb, nil
}

// KeepAlive keeps the sandbox alive.
func (s *Sandbox) KeepAlive(timeout time.Duration) error {
	time.Sleep(timeout)
	// TODO: implement
	return nil
}

// Reconnect reconnects to the sandbox.
func (s *Sandbox) Reconnect() error {
	u := s.wsURL()
	s.logger.Debug("Reconnecting to sandbox", "url", u.String())
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	s.ws = ws
	return nil
}

// Disconnect disconnects from the sandbox.
func (s *Sandbox) Disconnect() error {
	return s.ws.Close()
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

func (s *Sandbox) hostname(id string) string {
	return fmt.Sprintf("https://%s-%s-%s.e2b.dev",
		id,
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

// ListKernels lists the kernels in the sandbox.
//
// Make sure that the sandbox supports kernels before calling this method.
// The template must be set to "code-interpreter-stateful" or similar.
func (s *Sandbox) ListKernels(ctx context.Context) ([]ListKernelResponse, error) {
	req, err := s.newRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s%s", s.hostname("8888"), kernelsRoute),
	)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res []ListKernelResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Mkdir makes a directory in the sandbox file system.
func (s *Sandbox) Mkdir(path string) error {
	s.logger.Debug("Making directory", "path", path)
	s.msgCnt++
	jsVal, err := json.Marshal(Request{
		Params:  []any{path},
		JSONRPC: rpc,
		ID:      s.msgCnt,
		Method:  filesystemMakeDir,
	})
	if err != nil {
		return err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsVal)
	if err != nil {
		return err
	}
	_, _, err = s.ws.ReadMessage()
	if err != nil {
		return err
	}
	return nil
}

// Ls lists the files and/or directories in the sandbox file system at
// the given path.
func (s *Sandbox) Ls(path string) ([]LsResult, error) {
	s.logger.Debug("Listing files and dirs", "path", path)
	s.msgCnt++
	jsVal, err := json.Marshal(Request{
		Params:  []any{path},
		JSONRPC: rpc,
		ID:      s.msgCnt,
		Method:  filesystemList,
	})
	if err != nil {
		return nil, err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsVal)
	if err != nil {
		return nil, err
	}
	_, msr, err := s.ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	var res LsResponse
	err = json.Unmarshal(msr, &res)
	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

// Read reads a file from the sandbox file system.
func (s *Sandbox) Read(
	path string,
) ([]byte, error) {
	s.logger.Debug("Reading from file", "path", path)
	s.msgCnt++
	jsnV, err := json.Marshal(Request{
		JSONRPC: rpc,
		Method:  filesystemRead,
		Params:  []any{path},
		ID:      s.msgCnt,
	})
	if err != nil {
		return nil, err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsnV)
	if err != nil {
		return nil, err
	}
	_, message, err := s.ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	var resp ReadResponse
	err = json.Unmarshal(message, &resp)
	if err != nil {
		return nil, err
	}
	return []byte(resp.Result), nil
}

// Write writes to a file to the sandbox file system.
func (s *Sandbox) Write(path string, data []byte) error {
	s.logger.Debug("Writing to file", "path", path)
	s.msgCnt++
	jsnV, err := json.Marshal(Request{
		JSONRPC: rpc,
		Method:  filesystemWrite,
		Params: []any{
			path,
			string(data),
		},
		ID: s.msgCnt,
	})
	if err != nil {
		return err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsnV)
	if err != nil {
		return err
	}
	_, _, err = s.ws.ReadMessage()
	if err != nil {
		return err
	}
	return nil
}

// ReadBytes reads a file from the sandbox file system.
func (s *Sandbox) ReadBytes(path string) ([]byte, error) {
	s.logger.Debug("Reading Bytes", "path", path)
	s.msgCnt++
	jsnV, err := json.Marshal(Request{
		JSONRPC: rpc,
		Method:  filesystemReadBytes,
		Params: []any{
			path,
		},
		ID: s.msgCnt,
	})
	if err != nil {
		return nil, err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsnV)
	if err != nil {
		return nil, err
	}
	_, message, err := s.ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	var rR ReadResponse
	err = json.Unmarshal(message, &rR)
	if err != nil {
		return nil, err
	}
	sDec, err := base64.StdEncoding.DecodeString(string(rR.Result))
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
	s.logger.Debug("Uploading file", "path", path)
	// TODO: implement
	return nil
}

// Download downloads a file from the sandbox file system.
func (s *Sandbox) Download(path string) (io.ReadCloser, error) {
	s.logger.Debug("Downloading file", "path", path)
	// TODO: implement
	return nil, nil
}

// {"jsonrpc": "2.0", "method": "process_start", "params": ["KkLECSZQiN5B", "cat file0.txt", {"PYTHONUNBUFFERED": "1"}, ""], "id": 12}
// {"jsonrpc": "2.0", "method": "process_start", "params": ["Z9SalhcNx641", "cat file9.txt", {"PYTHONUNBUFFERED": "1"}, ""], "id": 341}
// {"jsonrpc": "2.0", "method": "process_subscribe", "params": ["onExit", "N5hJqKkNXj1i"], "id": 15}
// {"jsonrpc": "2.0", "method": "process_subscribe", "params": ["onStdout", "N5hJqKkNXj1i"], "id": 16}
// {"jsonrpc": "2.0", "method": "process_unsubscribe", "params": ["0xa7966b61d145231b3b3ab8cd440edf58"], "id": 14}
// {"jsonrpc": "2.0", "method": "process_unsubscribe", "params": ["0xb6b65c652bc5576751debfc82e864156"], "id": 17}

type processSubscribeRequest struct {
	// JSONRPC is the JSON-RPC version of the message.
	JSONRPC string `json:"jsonrpc"`
	// Method is the method of the message.
	Method Method `json:"method"`
	// ID is the ID of the message.
	ID int `json:"id"`
	// Params is the params of the message.
	Params []any `json:"params"`
}

func createProcessID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// StartProcess starts a process in the sandbox.
func (s *Sandbox) StartProcess(
	cmd string,
) (proc Process, err error) {
	s.msgCnt++
	proc.ID = createProcessID(12)
	req := Request{
		JSONRPC: rpc,
		Method:  processStart,
		ID:      s.msgCnt,
		Params: []any{
			createProcessID(12),
			cmd,
			map[string]string{},
			"",
		},
	}
	jsVal, err := json.Marshal(req)
	if err != nil {
		return proc, err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsVal)
	if err != nil {
		return proc, err
	}
	_, msr, err := s.ws.ReadMessage()
	if err != nil {
		return proc, err
	}

	// {"jsonrpc":"2.0","id":2,"result":"ewMUmGQ0vVmW"}
	var res processStartResponse
	err = json.Unmarshal(msr, &res)
	if err != nil {
		return proc, err
	}
	if res.Result == "" {
		return proc, fmt.Errorf("process start failed got empty result id")
	}
	return Process{
		ID:       proc.ID,
		ResultID: res.Result,
	}, nil
}

type processStartResponse struct {
	// JSONRPC is the JSON-RPC version of the message.
	JSONRPC string `json:"jsonrpc"`
	// Method is the method of the message.
	Method Method `json:"method"`
	// ID is the ID of the message.
	ID int `json:"id"`
	// Result is the result of the message.
	Result string `json:"result"`
}

func (s *Sandbox) subscribeProcess(ctx context.Context, id string, event ProcessEvents) error {
	s.logger.Debug("Subscribing to process", "id", id)
	req := Request{
		JSONRPC: rpc,
		Method:  processSubscribe,
		ID:      s.msgCnt,
		Params: []any{
			event,
			id,
		},
	}
	jsVal, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsVal)
	if err != nil {
		return err
	}
	_, msr, err := s.ws.ReadMessage()
	if err != nil {
		return err
	}
	println(string(msr))
	return nil
}

func (s *Sandbox) unsubscribeProcess(ctx context.Context, id string, event ProcessEvents) error {
	s.logger.Debug("Unsubscribing from process", "id", id)
	req := Request{
		JSONRPC: rpc,
		Method:  processUnsubscribe,
		ID:      s.msgCnt,
		Params: []any{
			event,
			id,
		},
	}
	jsVal, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsVal)
	if err != nil {
		return err
	}
	_, msr, err := s.ws.ReadMessage()
	if err != nil {
		return err
	}
	println(string(msr))
	return nil
}
