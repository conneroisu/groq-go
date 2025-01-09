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
	"sync"
	"time"

	"github.com/conneroisu/groq-go/pkg/builders"
	"github.com/gorilla/websocket"
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
		ID       string                  `json:"sandboxID"`  // ID of the sandbox.
		ClientID string                  `json:"clientID"`   // ClientID of the sandbox.
		Cwd      string                  `json:"cwd"`        // Cwd is the sandbox's current working directory.
		apiKey   string                  `json:"-"`          // apiKey is the sandbox's api key.
		Template SandboxTemplate         `json:"templateID"` // Template of the sandbox.
		baseURL  string                  `json:"-"`          // baseAPIURL is the base api url of the sandbox.
		Metadata map[string]string       `json:"metadata"`   // Metadata of the sandbox.
		logger   *slog.Logger            `json:"-"`          // logger is the sandbox's logger.
		client   *http.Client            `json:"-"`          // client is the sandbox's http client.
		header   builders.Header         `json:"-"`          // header is the sandbox's request header builder.
		ws       *websocket.Conn         `json:"-"`          // ws is the sandbox's websocket connection.
		wsURL    func(s *Sandbox) string `json:"-"`          // wsURL is the sandbox's websocket url.
		Map      *sync.Map               `json:"-"`          // Map is the map of the sandbox.
		idCh     chan int                `json:"-"`          // idCh is the channel to generate ids for requests.
		toolW    ToolingWrapper          `json:"-"`          // toolW is the tooling wrapper for the sandbox.
	}
	// Option is an option for the sandbox.
	Option func(*Sandbox)
	// Process is a process in the sandbox.
	Process struct {
		id  string            // ID is process id.
		cmd string            // cmd is process's command.
		Cwd string            // cwd is process's current working directory.
		sb  *Sandbox          // sb is the sandbox the process belongs to.
		Env map[string]string // env is process's environment variables.
	}
	// ProcessOption is an option for the process.
	ProcessOption func(*Process)
	// Event is a file system event.
	Event struct {
		Path      string      `json:"path"`      // Path is the path of the event.
		Name      string      `json:"name"`      // Name is the name of file or directory.
		Timestamp int64       `json:"timestamp"` // Timestamp is the timestamp of the event.
		Error     string      `json:"error"`     // Error is the possible error of the event.
		Params    EventParams `json:"params"`    // Params is the parameters of the event.
	}
	// EventParams is the params for subscribing to a process event.
	EventParams struct {
		Subscription string      `json:"subscription"` // Subscription is the subscription id of the event.
		Result       EventResult `json:"result"`       // Result is the result of the event.
	}
	// EventResult is a file system event response.
	EventResult struct {
		Type        string `json:"type"`
		Line        string `json:"line"`
		Timestamp   int64  `json:"timestamp"`
		IsDirectory bool   `json:"isDirectory"`
		Error       string `json:"error"`
	}
	// Request is a JSON-RPC request.
	Request struct {
		JSONRPC string `json:"jsonrpc"` // JSONRPC is the JSON-RPC version of the request.
		Method  Method `json:"method"`  // Method is the request method.
		ID      int    `json:"id"`      // ID of the request.
		Params  []any  `json:"params"`  // Params of the request.
	}
	// Response is a JSON-RPC response.
	Response[T any, Q any] struct {
		ID     int `json:"id"`     // ID of the response.
		Result T   `json:"result"` // Result of the response.
		Error  Q   `json:"error"`  // Error of the message.
	}
	// LsResult is a result of the list request.
	LsResult struct {
		Name  string `json:"name"`  // Name is the name of the file or directory.
		IsDir bool   `json:"isDir"` // isDir is true if the entry is a directory.
	}
	// APIError is the error of the API.
	APIError struct {
		Code    int    `json:"code,omitempty"` // Code is the code of the error.
		Message string `json:"message"`        // Message is the message of the error.
	}
	// Method is a JSON-RPC method.
	Method string
)

const (
	OnStdout ProcessEvents = "onStdout" // OnStdout is the event for the stdout.
	OnStderr ProcessEvents = "onStderr" // OnStderr is the event for the stderr.
	OnExit   ProcessEvents = "onExit"   // OnExit is the event for the exit.

	rpc                = "2.0"
	charset            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	defaultBaseURL     = "https://api.e2b.dev"
	defaultWSScheme    = "wss"
	wsRoute            = "/ws"
	fileRoute          = "/file"
	sandboxesRoute     = "/sandboxes"  // (GET/POST /sandboxes)
	deleteSandboxRoute = "/sandboxes/" // (DELETE /sandboxes/:id)

	filesystemWrite      Method = "filesystem_write"
	filesystemRead       Method = "filesystem_read"
	filesystemList       Method = "filesystem_list"
	filesystemRemove     Method = "filesystem_remove"
	filesystemMakeDir    Method = "filesystem_makeDir"
	filesystemReadBytes  Method = "filesystem_readBase64"
	filesystemWriteBytes Method = "filesystem_writeBase64"
	filesystemSubscribe  Method = "filesystem_subscribe"
	processSubscribe     Method = "process_subscribe"
	processUnsubscribe   Method = "process_unsubscribe"
	processStart         Method = "process_start"
)

// NewSandbox creates a new sandbox.
func NewSandbox(
	ctx context.Context,
	apiKey string,
	opts ...Option,
) (*Sandbox, error) {
	sb := Sandbox{
		apiKey:   apiKey,
		Template: "base",
		baseURL:  defaultBaseURL,
		Metadata: map[string]string{
			"sdk": "groq-go v1",
		},
		client: http.DefaultClient,
		logger: slog.New(slog.NewJSONHandler(io.Discard, nil)),
		toolW:  defaultToolWrapper,
		idCh:   make(chan int),
		Map:    new(sync.Map),
		wsURL: func(s *Sandbox) string {
			return fmt.Sprintf("wss://49982-%s-%s.e2b.dev/ws", s.ID, s.ClientID)
		},
		header: builders.Header{
			SetCommonHeaders: func(req *http.Request) {
				req.Header.Set("X-API-Key", apiKey)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Accept", "application/json")
			},
		},
	}
	for _, opt := range opts {
		opt(&sb)
	}
	req, err := builders.NewRequest(
		ctx, sb.header, http.MethodPost,
		fmt.Sprintf("%s%s", sb.baseURL, sandboxesRoute),
		builders.WithBody(&sb),
	)
	if err != nil {
		return &sb, err
	}
	err = sb.sendRequest(req, &sb)
	if err != nil {
		return &sb, err
	}
	sb.ws, _, err = websocket.DefaultDialer.Dial(sb.wsURL(&sb), nil)
	if err != nil {
		return &sb, err
	}
	go sb.identify(ctx)
	go func() {
		err := sb.read(ctx)
		if err != nil {
			sb.logger.Error("failed to read sandbox", "error", err)
		}
	}()
	return &sb, nil
}

// KeepAlive keeps the sandbox alive.
func (s *Sandbox) KeepAlive(ctx context.Context, timeout time.Duration) error {
	req, err := builders.NewRequest(
		ctx, s.header, http.MethodPost,
		fmt.Sprintf("%s/sandboxes/%s/refreshes", s.baseURL, s.ID),
		builders.WithBody(struct {
			Duration int `json:"duration"`
		}{Duration: int(timeout.Seconds())}),
	)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("request to create sandbox failed: %s", resp.Status)
	}
	return nil
}

// Reconnect reconnects to the sandbox.
func (s *Sandbox) Reconnect(ctx context.Context) (err error) {
	if err := s.ws.Close(); err != nil {
		return err
	}
	urlu := fmt.Sprintf("wss://49982-%s-%s.e2b.dev/ws", s.ID, s.ClientID)
	s.ws, _, err = websocket.DefaultDialer.Dial(urlu, nil)
	if err != nil {
		return err
	}
	go func() {
		err := s.read(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return err
}

// Stop stops the sandbox.
func (s *Sandbox) Stop(ctx context.Context) error {
	req, err := builders.NewRequest(
		ctx, s.header, http.MethodDelete,
		fmt.Sprintf("%s%s%s", s.baseURL, deleteSandboxRoute, s.ID),
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
	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("request to delete sandbox failed: %s", resp.Status)
	}
	return nil
}

// Mkdir makes a directory in the sandbox file system.
func (s *Sandbox) Mkdir(ctx context.Context, path string) error {
	respCh := make(chan []byte)
	err := s.writeRequest(ctx, filesystemMakeDir, []any{path}, respCh)
	if err != nil {
		return err
	}
	select {
	case body := <-respCh:
		resp, err := decodeResponse[string, APIError](body)
		if err != nil {
			return fmt.Errorf("failed to mkdir: %w", err)
		}
		if resp.Error.Code != 0 {
			return fmt.Errorf("failed to write to file: %s", resp.Error.Message)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Ls lists the files and/or directories in the sandbox file system at
// the given path.
func (s *Sandbox) Ls(ctx context.Context, path string) ([]LsResult, error) {
	respCh := make(chan []byte)
	defer close(respCh)
	err := s.writeRequest(ctx, filesystemList, []any{path}, respCh)
	if err != nil {
		return nil, err
	}
	select {
	case body := <-respCh:
		res, err := decodeResponse[[]LsResult, string](body)
		if err != nil {
			return nil, err
		}
		return res.Result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Read reads a file from the sandbox file system.
func (s *Sandbox) Read(
	ctx context.Context,
	path string,
) (string, error) {
	respCh := make(chan []byte)
	err := s.writeRequest(ctx, filesystemRead, []any{path}, respCh)
	if err != nil {
		return "", err
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case body := <-respCh:
		res, err := decodeResponse[string, string](body)
		if err != nil {
			return "", err
		}
		if res.Error != "" {
			return "", fmt.Errorf("failed to read file: %s", res.Error)
		}
		return res.Result, nil
	}
}

// Write writes to a file to the sandbox file system.
func (s *Sandbox) Write(ctx context.Context, path string, data []byte) error {
	respCh := make(chan []byte)
	err := s.writeRequest(ctx, filesystemWrite, []any{path, string(data)}, respCh)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case resp := <-respCh:
		err = json.Unmarshal(resp, &Request{})
		if err != nil {
			return err
		}
		return nil
	}
}

// ReadBytes reads a file from the sandbox file system.
func (s *Sandbox) ReadBytes(ctx context.Context, path string) ([]byte, error) {
	resCh := make(chan []byte)
	defer close(resCh)
	err := s.writeRequest(ctx, filesystemReadBytes, []any{path}, resCh)
	if err != nil {
		return nil, err
	}
	select {
	case body := <-resCh:
		res, err := decodeResponse[string, string](body)
		if err != nil {
			return nil, err
		}
		sDec, err := base64.StdEncoding.DecodeString(res.Result)
		if err != nil {
			return nil, err
		}
		return sDec, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Watch watches a directory in the sandbox file system.
//
// This is intended to be run in a goroutine as it will block until the
// connection is closed, an error occurs, or the context is canceled.
//
// While blocking, filesystem events will be written to the provided channel.
func (s *Sandbox) Watch(
	ctx context.Context,
	path string,
	eCh chan<- Event,
) error {
	respCh := make(chan []byte)
	defer close(respCh)
	err := s.writeRequest(ctx, filesystemSubscribe, []any{"watchDir", path}, respCh)
	if err != nil {
		return err
	}
	res, err := decodeResponse[string, string](<-respCh)
	if err != nil {
		return err
	}
	s.Map.Store(res.Result, eCh)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var event Event
				err := json.Unmarshal(<-respCh, &event)
				if err != nil {
					return
				}
				if event.Error != "" {
					return
				}
				if event.Params.Subscription != path {
					continue
				}
				eCh <- event
			}
		}
	}()
	return nil
}

// NewProcess creates a new process startable in the sandbox.
func (s *Sandbox) NewProcess(
	cmd string,
	opts ...ProcessOption,
) (*Process, error) {
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	proc := &Process{
		id:  string(b),
		sb:  s,
		cmd: cmd,
	}
	for _, opt := range opts {
		opt(proc)
	}
	if proc.Cwd == "" {
		proc.Cwd = s.Cwd
	}
	return proc, nil
}

// Start starts a process in the sandbox.
func (p *Process) Start(ctx context.Context) (err error) {
	if p.Env == nil {
		p.Env = map[string]string{"PYTHONUNBUFFERED": "1"}
	}
	respCh := make(chan []byte)
	err = p.sb.writeRequest(
		ctx,
		processStart,
		[]any{p.id, p.cmd, p.Env, p.Cwd},
		respCh,
	)
	if err != nil {
		return err
	}
	select {
	case body := <-respCh:
		res, err := decodeResponse[string, APIError](body)
		if err != nil {
			return err
		}
		if res.Error.Code != 0 {
			return fmt.Errorf("process start failed(%d): %s", res.Error.Code, res.Error.Message)
		}
		if res.Result == "" || len(res.Result) == 0 {
			return fmt.Errorf("process start failed got empty result id")
		}
		if p.id != res.Result {
			return fmt.Errorf("process start failed got wrong result id; want %s, got %s", p.id, res.Result)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Done returns a channel that is closed when the process is done.
func (p *Process) Done() <-chan struct{} {
	rCh, ok := p.sb.Map.Load(p.id)
	if !ok {
		return nil
	}
	return rCh.(<-chan struct{})
}

// SubscribeStdout subscribes to the process's stdout.
func (p *Process) SubscribeStdout(ctx context.Context) (chan Event, chan error) {
	return p.subscribe(ctx, OnStdout)
}

// SubscribeStderr subscribes to the process's stderr.
func (p *Process) SubscribeStderr(ctx context.Context) (chan Event, chan error) {
	return p.subscribe(ctx, OnStderr)
}

// SubscribeExit subscribes to the process's exit.
func (p *Process) SubscribeExit(ctx context.Context) (chan Event, chan error) {
	return p.subscribe(ctx, OnExit)
}

// Subscribe subscribes to a process event.
//
// It creates a go routine to read the process events into the provided channel.
func (p *Process) subscribe(
	ctx context.Context,
	event ProcessEvents,
) (chan Event, chan error) {
	events := make(chan Event)
	errs := make(chan error)
	go func(errCh chan error) {
		respCh := make(chan []byte)
		defer close(respCh)
		err := p.sb.writeRequest(ctx, processSubscribe, []any{event, p.id}, respCh)
		if err != nil {
			errCh <- err
		}
		res, err := decodeResponse[string, any](<-respCh)
		if err != nil {
			errCh <- err
		}
		p.sb.Map.Store(res.Result, respCh)
	loop:
		for {
			select {
			case eventBd := <-respCh:
				var event Event
				_ = json.Unmarshal(eventBd, &event)
				if event.Error != "" {
					p.sb.logger.Error("failed to read event", "error", event.Error)
					continue
				}
				events <- event
			case <-ctx.Done():
				break loop
			case <-p.Done():
				break loop
			}
		}

		p.sb.Map.Delete(res.Result)
		finishCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		p.sb.logger.Debug("unsubscribing from process", "event", event, "id", res.Result)
		_ = p.sb.writeRequest(finishCtx, processUnsubscribe, []any{res.Result}, respCh)
		unsubRes, _ := decodeResponse[bool, string](<-respCh)
		if unsubRes.Error != "" || !unsubRes.Result {
			p.sb.logger.Debug("failed to unsubscribe from process", "error", unsubRes.Error)
		}
	}(errs)
	return events, errs
}
func (s *Sandbox) sendRequest(req *http.Request, v interface{}) error {
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK ||
		res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("failed to create sandbox: %s", res.Status)
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
		return json.NewDecoder(res.Body).Decode(v)
	}
}
func decodeResponse[T any, Q any](body []byte) (*Response[T, Q], error) {
	decResp := new(Response[T, Q])
	err := json.Unmarshal(body, decResp)
	if err != nil {
		return nil, err
	}
	return decResp, nil
}
func (s *Sandbox) read(ctx context.Context) error {
	var body []byte
	var err error
	type decResp struct {
		Method string `json:"method"`
		ID     int    `json:"id"`
		Params struct {
			Subscription string `json:"subscription"`
		}
	}
	defer func() {
		err := s.ws.Close()
		if err != nil {
			s.logger.Error("failed to close sandbox", "error", err)
		}
	}()
	msgCh := make(chan []byte, 10)
	for {
		select {
		case body = <-msgCh:
			var decResp decResp
			err = json.Unmarshal(body, &decResp)
			if err != nil {
				return err
			}
			var key any
			key = decResp.Params.Subscription
			if decResp.ID != 0 {
				key = decResp.ID
			}
			toR, ok := s.Map.Load(key)
			if !ok {
				msgCh <- body
				continue
			}
			toRCh, ok := toR.(chan []byte)
			if !ok {
				msgCh <- body
				continue
			}
			s.logger.Debug("read",
				"subscription", decResp.Params.Subscription,
				"body", body,
				"sandbox", s.ID,
			)
			toRCh <- body
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, msg, err := s.ws.ReadMessage()
			if err != nil {
				return err
			}
			msgCh <- msg
		}
	}
}
func (s *Sandbox) writeRequest(
	ctx context.Context,
	method Method,
	params []any,
	respCh chan []byte,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case id := <-s.idCh:
		req := Request{
			Method:  method,
			JSONRPC: rpc,
			Params:  params,
			ID:      id,
		}
		s.logger.Debug("request",
			"sandbox", id,
			"method", method,
			"id", id,
			"params", params,
		)
		s.Map.Store(req.ID, respCh)
		jsVal, err := json.Marshal(req)
		if err != nil {
			return err
		}
		err = s.ws.WriteMessage(websocket.TextMessage, jsVal)
		if err != nil {
			return fmt.Errorf(
				"writing %s request failed (%d): %w",
				method,
				req.ID,
				err,
			)
		}
		return nil
	}
}
func (s *Sandbox) identify(ctx context.Context) {
	id := 1
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.idCh <- id
			id++
		}
	}
}
