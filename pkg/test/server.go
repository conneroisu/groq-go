package test

import (
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
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
	Logger   *slog.Logger
}

// handler is a function that handles a request.
type handler func(w http.ResponseWriter, r *http.Request)

// NewTestServer creates a new test server.
func NewTestServer() *ServerTest {
	return &ServerTest{
		handlers: make(map[string]handler),
		Logger:   defaultLogger,
	}
}

// RegisterHandler registers a handler for a path.
func (ts *ServerTest) RegisterHandler(path string, handler handler) {
	// to make the registered paths friendlier to a regex match in the route handler
	// in GroqTestServer
	path = strings.ReplaceAll(path, "*", ".*")
	ts.handlers[path] = handler
}

// GroqTestServer Creates a mocked Groq server which can pretend to handle requests during testing.
func (ts *ServerTest) GroqTestServer() *httptest.Server {
	return httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf(
				"received a %s request at path %q\n",
				r.Method,
				r.URL.Path,
			)

			// check auth
			if r.Header.Get("Authorization") != "Bearer "+GetTestToken() &&
				r.Header.Get("api-key") != GetTestToken() {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Handle /path/* routes.
			// Note: the * is converted to a .* in register handler for proper regex handling
			for route, handler := range ts.handlers {
				// Adding ^ and $ to make path matching deterministic since go map iteration isn't ordered
				pattern, _ := regexp.Compile("^" + route + "$")
				if pattern.MatchString(r.URL.Path) {
					handler(w, r)
					return
				}
			}
			http.Error(
				w,
				"the resource path doesn't exist",
				http.StatusNotFound,
			)
		}),
	)
}

// E2bTestServer creates a test server for emulating the e2b api.
func (ts *ServerTest) E2bTestServer() *httptest.Server {
	return httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf(
				"received a %s request at path %q\n",
				r.Method,
				r.URL.Path,
			)

			// check auth
			if r.Header.Get("X-API-Key") != GetTestToken() &&
				r.Header.Get("api-key") != GetTestToken() {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Handle /path/* routes.
			// Note: the * is converted to a .* in register handler for proper regex handling
			for route, handler := range ts.handlers {
				// Adding ^ and $ to make path matching deterministic since go map iteration isn't ordered
				pattern, _ := regexp.Compile("^" + route + "$")
				if pattern.MatchString(r.URL.Path) {
					handler(w, r)
					return
				}
			}
			http.Error(
				w,
				"the resource path doesn't exist",
				http.StatusNotFound,
			)
		}),
	)
}

var (
	defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			if a.Key == "level" {
				return slog.Attr{}
			}
			if a.Key == slog.SourceKey {
				str := a.Value.String()
				split := strings.Split(str, "/")
				if len(split) > 2 {
					a.Value = slog.StringValue(strings.Join(split[len(split)-2:], "/"))
					a.Value = slog.StringValue(strings.Replace(a.Value.String(), "}", "", -1))
				}
				a.Key = a.Value.String()
				a.Value = slog.IntValue(0)
			}
			if a.Key == "body" {
				a.Value = slog.StringValue(strings.Replace(a.Value.String(), "/", "", -1))
				a.Value = slog.StringValue(strings.Replace(a.Value.String(), "\n", "", -1))
				a.Value = slog.StringValue(strings.Replace(a.Value.String(), "\"", "", -1))
			}
			return a
		}}))
)
