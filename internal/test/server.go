package test

import (
	"log/slog"
	"net/http"
	"strings"
)

const (
	testAPI = "this-is-my-secure-token-do-not-steal!!"
)

// GetTestToken returns the test token.
func GetTestToken() string {
	return testAPI
}

// ServerTest is a test server for the groq api.
type ServerTest struct {
	handlers map[string]Handler
	logger   *slog.Logger
}

// Handler is a function that handles a request.
type Handler func(w http.ResponseWriter, r *http.Request)

// NewTestServer creates a new test server.
func NewTestServer() *ServerTest {
	return &ServerTest{
		handlers: make(map[string]Handler),
		logger:   DefaultLogger,
	}
}

// RegisterHandler registers a handler for a path.
func (ts *ServerTest) RegisterHandler(path string, handler Handler) {
	// to make the registered paths friendlier to a regex match in the route handler
	// in GroqTestServer
	path = strings.ReplaceAll(path, "*", ".*")
	ts.handlers[path] = handler
}
