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
		ID         string            `json:"sandboxID"`  // ID of the sandbox.
		Metadata   map[string]string `json:"metadata"`   // Metadata of the sandbox.
		Template   SandboxTemplate   `json:"templateID"` // Template of the sandbox.
		ClientID   string            `json:"clientID"`   // ClientID of the sandbox.
		Cwd        string            `json:"cwd"`        // Cwd is the sandbox's current working directory.
		logger     *slog.Logger      `json:"-"`          // logger is the sandbox's logger.
		apiKey     string            `json:"-"`          // apiKey is the sandbox's api key.
		baseAPIURL string            `json:"-"`          // baseAPIURL is the base api url of the sandbox.
		httpScheme string            `json:"-"`          // httpScheme is the sandbox's http scheme.
		client     *http.Client      `json:"-"`          // client is the sandbox's http client.
		header     builders.Header   `json:"-"`          // header is the sandbox's request header builder.
		ws         *websocket.Conn   `json:"-"`          // ws is the sandbox's websocket connection.
		Map        sync.Map          `json:"-"`          // Map is the map of the sandbox.
		msgCnt     int               `json:"-"`          // msgCnt is the message count.
	}
	// Process is a process in the sandbox.
	Process struct {
		sb  *Sandbox          // sb is the sandbox the process belongs to.
		id  string            // ID is process id.
		cmd string            // cmd is process's command.
		Cwd string            // cwd is process's current working directory.
		Env map[string]string // env is process's environment variables.
	}
	// SubscribeParams is the params for subscribing to a process event.
	SubscribeParams struct {
		Event ProcessEvents // Event is the event to subscribe to.
		Ch    chan<- Event  // Ch is the channel to write the event to.
	}
	// Option is an option for the sandbox.
	Option func(*Sandbox)
	// ProcessOption is an option for a process.
	ProcessOption func(*Process)
	// Event is a file system event.
	Event struct {
		Path      string        `json:"path"`      // Path is the path of the event.
		Name      string        `json:"name"`      // Name is the name of file or directory.
		Timestamp int64         `json:"timestamp"` // Timestamp is the timestamp of the event.
		Error     string        `json:"error"`     // Error is the possible error of the event.
		Params    EventParams   `json:"params"`    // Params is the parameters of the event.
		Operation OperationType `json:"operation"` // Operation is the operation type of the event.
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
	// OperationType is an operation type.
	OperationType int
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
	Method  string
	decResp struct {
		Method string `json:"method"`
		ID     int    `json:"id"`
		Params struct {
			Subscription string `json:"subscription"`
		}
	}
)

const (
	OnStdout ProcessEvents = "onStdout" // OnStdout is the event for the stdout.
	OnStderr ProcessEvents = "onStderr" // OnStderr is the event for the stderr.
	OnExit   ProcessEvents = "onExit"   // OnExit is the event for the exit.

	EventTypeCreate OperationType = iota // EventTypeCreate is an event for the creation of a file/dir.
	EventTypeWrite                       // EventTypeWrite is an event for the write to a file.
	EventTypeRemove                      // EventTypeRemove is an event for the removal of a file/dir.

	rpc                = "2.0"
	charset            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	defaultBaseURL     = "api.e2b.dev"
	defaultWSScheme    = "wss"
	wsRoute            = "/ws"
	fileRoute          = "/file"
	sandboxesRoute     = "/sandboxes"  // (GET/POST /sandboxes)
	deleteSandboxRoute = "/sandboxes/" // (DELETE /sandboxes/:id)
	defaultHTTPScheme  = "https"

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
		apiKey:     apiKey,
		Template:   "base",
		baseAPIURL: defaultBaseURL,
		Metadata: map[string]string{
			"sdk": "groq-go v1",
		},
		client:     http.DefaultClient,
		logger:     slog.New(slog.NewJSONHandler(io.Discard, nil)),
		httpScheme: defaultHTTPScheme,
		msgCnt:     1,
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
		ctx, sb.header, http.MethodPost,
		fmt.Sprintf("%s://%s%s", sb.httpScheme, sb.baseAPIURL, sandboxesRoute),
		builders.WithBody(&sb),
	)
	if err != nil {
		return &sb, err
	}
	err = sb.sendRequest(req, &sb)
	if err != nil {
		return &sb, err
	}
	sb.ws, _, err = websocket.DefaultDialer.Dial(sb.wsURL().String(), nil)
	if err != nil {
		return &sb, err
	}
	go func() {
		err := sb.read(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return &sb, nil
}

// KeepAlive keeps the sandbox alive.
func (s *Sandbox) KeepAlive(ctx context.Context, timeout time.Duration) error {
	req, err := builders.NewRequest(
		ctx, s.header, http.MethodPost,
		fmt.Sprintf("%s://%s/sandboxes/%s/refreshes", s.httpScheme, s.baseAPIURL, s.ID),
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
		return fmt.Errorf("request to create sandbox failed: %s\nbody: %s", resp.Status, getBody(resp))
	}
	return nil
}

// Reconnect reconnects to the sandbox.
func (s *Sandbox) Reconnect(ctx context.Context) (err error) {
	if err := s.ws.Close(); err != nil {
		return err
	}
	u := s.wsURL()
	s.logger.Debug("reconnecting to sandbox", "url", u.String())
	s.ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
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
		fmt.Sprintf("%s://%s%s%s", s.httpScheme, s.baseAPIURL, deleteSandboxRoute, s.ID),
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
	err := s.WriteRequest(filesystemMakeDir, []any{path}, respCh)
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
	err := s.WriteRequest(filesystemList, []any{path}, respCh)
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
	path string,
) (string, error) {
	respCh := make(chan []byte)
	err := s.WriteRequest(filesystemRead, []any{path}, respCh)
	if err != nil {
		return "", err
	}
	res, err := decodeResponse[string, string](<-respCh)
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
	err := s.WriteRequest(filesystemWrite, []any{path, string(data)}, respCh)
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
func (s *Sandbox) ReadBytes(ctx context.Context, path string) ([]byte, error) {
	resCh := make(chan []byte)
	defer close(resCh)
	err := s.WriteRequest(filesystemReadBytes, []any{path}, resCh)
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
func (s *Sandbox) Watch(
	ctx context.Context,
	path string,
	eCh chan<- Event,
) error {
	respCh := make(chan []byte)
	defer close(respCh)
	err := s.WriteRequest(filesystemSubscribe, []any{"watchDir", path}, respCh)
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
	proc Process,
) (*Process, error) {
	proc.cmd = cmd
	proc.id = createProcessID()
	proc.sb = s
	if proc.Cwd == "" {
		proc.Cwd = s.Cwd
	}
	return &proc, nil
}

// Start starts a process in the sandbox.
func (p *Process) Start(ctx context.Context) (err error) {
	if p.Env == nil {
		p.Env = map[string]string{"PYTHONUNBUFFERED": "1"}
	}
	respCh := make(chan []byte)
	if err = p.sb.WriteRequest(processStart, []any{p.id, p.cmd, p.Env, p.Cwd}, respCh); err != nil {
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

// Subscribe subscribes to a process event.
func (p *Process) Subscribe(
	ctx context.Context,
	event ProcessEvents,
	ch chan<- Event,
) error {
	respCh := make(chan []byte)
	err := p.sb.WriteRequest(processSubscribe, []any{event, p.id}, respCh)
	if err != nil {
		return err
	}
	res, err := decodeResponse[string, APIError](<-respCh)
	if err != nil {
		return err
	}
	if res.Error.Code != 0 {
		return fmt.Errorf("process subscribe failed(%d): %s", res.Error.Code, res.Error.Message)
	}
	eventByCh := make(chan []byte)
	p.sb.Map.Store(res.Result, eventByCh)
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
	select {
	case <-ctx.Done():
		p.sb.Map.Delete(res.Result)
		break
	case <-p.Done():
		p.sb.Map.Delete(res.Result)
		break
	}
	p.sb.Map.Delete(res.Result)
	err = p.sb.WriteRequest(processUnsubscribe, []any{res.Result}, respCh)
	if err != nil {
		println(err)
	}
	unsubRes, err := decodeResponse[bool, string](<-respCh)
	if err != nil {
		println(err)
	}
	if unsubRes.Error != "" {
		println(unsubRes.Error)
	}
	return nil
}
func (s *Sandbox) wsURL() *url.URL {
	return &url.URL{
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
func (s *Sandbox) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/json")
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < http.StatusOK ||
		res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("request to create sandbox failed: %s\nbody: %s", res.Status, getBody(res))
	}
	if v == nil {
		return nil
	}
	switch o := v.(type) {
	case *string:
		return decodeString(res.Body, o)
	default:
		return json.NewDecoder(res.Body).Decode(v)
	}
}
func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}
func getBody(resp *http.Response) string {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(b)
}

// WriteRequest writes a request to the websocket.
func (s *Sandbox) WriteRequest(method Method, params []any, respCh chan []byte) error {
	req := Request{
		Method:  method,
		JSONRPC: rpc,
		Params:  params,
		ID:      s.msgCnt,
	}
	defer func() { s.msgCnt++ }()
	s.logger.Debug("write", "method", req.Method, "id", req.ID, "params", req.Params)
	s.Map.Store(req.ID, respCh)
	jsVal, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = s.ws.WriteMessage(websocket.TextMessage, jsVal)
	if err != nil {
		return fmt.Errorf("failed to write %s request (%d): %w", req.Method, req.ID, err)
	}
	return nil
}

// Read reads a response from the websocket.
//
// If the context is cancelled, the websocket will be closed.
func (s *Sandbox) read(ctx context.Context) (err error) {
	defer func() {
		err = s.ws.Close()
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, body, err := s.ws.ReadMessage()
			if err != nil {
				return err
			}
			var decResp decResp
			err = json.Unmarshal(body, &decResp)
			if err != nil {
				return err
			}
			s.logger.Debug("read", "method", decResp.Method, "id", decResp.ID, "body", body)
			if decResp.Params.Subscription != "" {
				toR, ok := s.Map.Load(decResp.Params.Subscription)
				if !ok {
					s.logger.Debug("subscription not found", "id", decResp.Params.Subscription)
				}
				toRCh, ok := toR.(chan []byte)
				if !ok {
					s.logger.Debug("subscription not found", "id", decResp.Params.Subscription)
				}
				toRCh <- body
				continue
			}
			// response has an id
			toR, ok := s.Map.Load(decResp.ID)
			if !ok {
				s.logger.Debug("response not found", "id", decResp.ID)
			}
			toRCh, ok := toR.(chan []byte)
			if !ok {
				s.logger.Debug("responsech not found", "id", decResp.ID)
			}
			toRCh <- body
		}
	}
}

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

// WithCwd sets the current working directory.
func WithCwd(cwd string) Option {
	return func(s *Sandbox) { s.Cwd = cwd }
}

func decodeResponse[T any, Q any](body []byte) (*Response[T, Q], error) {
	decResp := new(Response[T, Q])
	err := json.Unmarshal(body, decResp)
	if err != nil {
		return nil, err
	}
	return decResp, nil
}
