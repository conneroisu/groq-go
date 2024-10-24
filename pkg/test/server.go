package test

import (
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
	handlers map[string]handler
}

// handler is a function that handles a request.
type handler func(w http.ResponseWriter, r *http.Request)

// NewTestServer creates a new test server.
func NewTestServer() *ServerTest {
	return &ServerTest{
		handlers: make(map[string]handler),
	}
}

// RegisterHandler registers a handler for a path.
func (ts *ServerTest) RegisterHandler(path string, handler handler) {
	// to make the registered paths friendlier to a regex match in the route handler
	// in GroqTestServer
	path = strings.ReplaceAll(path, "*", ".*")
	ts.handlers[path] = handler
}
