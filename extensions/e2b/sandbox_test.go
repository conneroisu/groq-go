package e2b

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/conneroisu/groq-go/pkg/test"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{}

func echo(a *assert.Assertions) func(w http.ResponseWriter, r *http.Request) {
	mu := sync.Mutex{}
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			a.NoError(err)
			defaultLogger.Debug("server read message", "msg", message)
			req := decode(message)
			switch req.Method {
			case filesystemList:
				err = c.WriteMessage(mt, encode(Response[[]LsResult, string]{
					ID:    req.ID,
					Error: "",
					Result: []LsResult{
						{
							Name:  "hello.txt",
							IsDir: false,
						},
					},
				}))
				a.NoError(err)
			case filesystemRead:
				err = c.WriteMessage(mt, encode(Response[string, string]{
					ID:     req.ID,
					Error:  "",
					Result: "hello",
				}))
				a.NoError(err)
			case filesystemWrite:
				err = c.WriteMessage(mt, encode(Response[string, string]{
					ID:     req.ID,
					Error:  "",
					Result: "hello",
				}))
				a.NoError(err)
			case processStart:
				err = c.WriteMessage(mt, encode(Response[string, APIError]{
					ID:     req.ID,
					Error:  APIError{},
					Result: req.Params[0].(string),
				}))
				a.NoError(err)
			case processSubscribe:
				err = c.WriteMessage(mt, encode(Response[string, APIError]{
					ID:     req.ID,
					Error:  APIError{},
					Result: "test-proc-id",
				}))
				a.NoError(err)
				err = c.WriteMessage(mt, encode(Response[
					EventParams, APIError,
				]{
					ID:    req.ID,
					Error: APIError{},
					Result: EventParams{
						Subscription: "test-proc-id",
						Result: EventResult{
							Type:        "Stdout",
							Line:        "hello",
							Timestamp:   0,
							IsDirectory: false,
							Error:       "",
						},
					},
				}))
				a.NoError(err)
			case filesystemMakeDir:
				err = c.WriteMessage(mt, encode(Response[string, APIError]{
					ID:     req.ID,
					Error:  APIError{},
					Result: "",
				}))
			default:
				err = c.WriteMessage(mt, message)
				a.NoError(err)
			}
		}
	}
}

func encode(v any) []byte {
	res, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return res
}
func decode(bod []byte) Request {
	var req Request
	err := json.Unmarshal(bod, &req)
	if err != nil {
		panic(err)
	}
	return req
}

func TestNewSandbox(t *testing.T) {
	a := assert.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv := test.NewTestServer()
	ts := srv.E2bTestServer()
	wsts := httptest.NewServer(http.HandlerFunc(echo(a)))
	id := "test-sandbox-id"
	srv.RegisterHandler("/sandboxes", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(encode(&Sandbox{ID: id}))
		a.NoError(err)
	})
	ts.Start()
	u := "ws" + strings.TrimPrefix(wsts.URL, "http")
	// Create a new sandbox.
	sb, err := NewSandbox(
		ctx,
		test.GetTestToken(),
		WithLogger(defaultLogger),
		WithBaseURL(ts.URL),
		WithWsURL(func(_ *Sandbox) string {
			return u + "/ws"
		}),
	)
	a.NoError(err, "NewSandbox error")
	a.NotNil(sb, "NewSandbox returned nil")
	a.Equal(sb.ID, id)
	// Call ls on the sandbox.
	lsRes, err := sb.Ls(ctx, ".")
	a.NoError(err)
	a.NotEmpty(lsRes)
	// Call mkdir on the sandbox.
	err = sb.Mkdir(ctx, "hello")
	a.NoError(err)
	// Call write on the sandbox.
	err = sb.Write(ctx, "hello.txt", []byte("hello"))
	a.NoError(err)
	// Call read on the sandbox.
	readRes, err := sb.Read(ctx, "hello.txt")
	a.NoError(err)
	a.Equal("hello", readRes)
	// create a process
	proc, err := sb.NewProcess("sleep 5 && echo 'hello world!'", Process{})
	a.NoError(err)
	events := make(chan Event, 10)
	err = proc.Start(ctx)
	a.NoError(err)
	err = proc.Subscribe(ctx, OnStdout, events)
	a.NoError(err)
	event := <-events
	jsnBytes, err := json.MarshalIndent(&event, "", "  ")
	a.NoError(err)
	t.Logf("test got event: %s", string(jsnBytes))
}
